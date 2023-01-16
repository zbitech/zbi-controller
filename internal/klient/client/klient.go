package client

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/zbitech/controller/internal/helper"
	"github.com/zbitech/controller/internal/utils"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
)

var (
	decUnstructured       = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	KUBERNETES_IN_CLUSTER = utils.GetEnv("KUBERNETES_IN_CLUSTER", "true")
	KUBECONFIG            = utils.GetEnv("KUBECONFIG", "/etc/zbi/kubeconfig")
)

type Klient struct {
	KubernetesClient kubernetes.Interface
	DynamicClient    dynamic.Interface
	restConfg        *rest.Config
}

func NewRestConfig(ctx context.Context) (*rest.Config, error) {

	var log = logger.GetLogger(ctx)
	var restConfig *rest.Config
	var err error
	if KUBERNETES_IN_CLUSTER == "true" {
		log.Infof("connecting to in-cluster kubernetes")
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		log.Infof("connecting to kubernetes with config file - %s", KUBECONFIG)
		restConfig, err = clientcmd.BuildConfigFromFlags("", KUBECONFIG)
		if err != nil {
			return nil, err
		}
	}
	return restConfig, nil
}

func NewKlient(ctx context.Context) (*Klient, error) {

	var log = logger.GetLogger(ctx)
	cfg, err := NewRestConfig(ctx)
	if err != nil {
		log.Errorf("Unable to create kubernetes configuration - %s", err)
		//		return nil, errs.NewApplicationError(errs.KubernetesError)
		return nil, err
	}

	kubernetesClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &Klient{
		KubernetesClient: kubernetesClient,
		DynamicClient:    dynamicClient,
		restConfg:        cfg,
	}, nil
}

func (k *Klient) GetKubernetesClient() kubernetes.Interface {
	return k.KubernetesClient
}

func (k *Klient) GetDynamicClient() dynamic.Interface {
	return k.DynamicClient
}

func (k *Klient) ApplyResource(ctx context.Context, object *unstructured.Unstructured) (*model.KubernetesResource, error) {

	var log = logger.GetServiceLogger(ctx, "klient.ApplyResource")
	defer func() { logger.LogServiceTime(log) }()

	data, err := json.Marshal(object)
	if err != nil {
		log.Errorf("failed to marshal resource - %s", err)
		// return nil, errs.NewApplicationError(errs.MarshalError, err)
		return nil, err
	}

	dr := helper.GetDynamicResourceInterface(k.DynamicClient, object)
	result, err := dr.Patch(ctx, object.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{FieldManager: "zbi-controller"})
	if err != nil {
		log.Errorf("failed to create resource - %s", err)
		return nil, fmt.Errorf("failed to create %s %s - %s", object.GetKind(), object.GetName(), err)
	} else {
		log.Infof("successfully created %s of kind %s", result.GetName(), result.GetKind())
	}

	// set all resources to active after initial creation
	//var status = utils.GetResourceStatusField(result)
	//var properties = utils.GetResourceProperties(result)
	//created := time.Now()
	//objType := model.ResourceObjectType(object.GetKind())
	//
	//resource := model.KubernetesResource{
	//	Name: object.GetName(),
	//	//		Project:    object.Project,
	//	//		Instance:   object.Instance,
	//	//		Namespace:  object.GetNamespace(),
	//	Type:    objType,
	//	Status:  status,
	//	Created: &created,
	//	//		Updated:    &created,
	//	Properties: properties,
	//}

	return helper.CreateUnstructuredResource(result), err
}

func (k *Klient) ApplyResources(ctx context.Context, objects []unstructured.Unstructured) ([]model.KubernetesResource, error) {

	var log = logger.GetServiceLogger(ctx, "klient.ApplyResources")
	defer func() { logger.LogServiceTime(log) }()

	results := make([]model.KubernetesResource, 0)

	var e error

	for _, object := range objects {
		res, err := k.ApplyResource(ctx, &object)
		if err != nil {
			log.WithFields(logrus.Fields{"error": err, "resource": object}).Errorf("failed to create resource")
			if e == nil {
				e = err
			} else {
				e = errors.Wrap(e, err.Error())
			}
		} else {
			log.WithFields(logrus.Fields{"resource": res}).Infof("created resource")
			results = append(results, *res)
		}
	}

	return results, e
}

func (k *Klient) DeleteResource(ctx context.Context, resource *model.KubernetesResource) error {

	var log = logger.GetServiceLogger(ctx, "klient.DeleteResource")
	defer func() { logger.LogServiceTime(log) }()

	err := k.DeleteDynamicResource(ctx, resource.Namespace, resource.Name, helper.GvrMap[resource.Type])
	if err != nil {
		return err
	}

	resource.Status = "Deleted"
	return nil
}

func (k *Klient) DeleteResources(ctx context.Context, resources []model.KubernetesResource) ([]model.KubernetesResource, error) {

	var log = logger.GetServiceLogger(ctx, "klient.DeleteResources")
	defer func() { logger.LogServiceTime(log) }()

	var e error

	var newResources = make([]model.KubernetesResource, 0)
	for _, resource := range resources {
		err := k.DeleteResource(ctx, &resource)
		if err != nil {
			if e == nil {
				e = err
			} else {
				e = errors.Wrap(e, err.Error())
			}
		}
		newResources = append(newResources, resource)
	}

	return newResources, e
}

func (k *Klient) DeleteDynamicResource(ctx context.Context, namespace, name string, resource schema.GroupVersionResource) error {
	var log = logger.GetServiceLogger(ctx, "klient.DeleteDynamicResource")
	defer func() { logger.LogServiceTime(log) }()

	return k.DynamicClient.Resource(resource).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (k *Klient) DeleteNamespace(ctx context.Context, namespace string) error {
	var log = logger.GetServiceLogger(ctx, "klient.DeleteNamespace")
	defer func() { logger.LogServiceTime(log) }()

	return k.KubernetesClient.CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
}

func (k *Klient) GetDynamicResource(ctx context.Context, namespace, name string, resource schema.GroupVersionResource) (*unstructured.Unstructured, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetDynamicResource")
	defer func() { logger.LogServiceTime(log) }()

	return k.DynamicClient.Resource(resource).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (k *Klient) GetDynamicResourceList(ctx context.Context, namespace string, resource schema.GroupVersionResource) ([]unstructured.Unstructured, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetDynamicResourceList")
	defer func() { logger.LogServiceTime(log) }()

	resultList, err := k.DynamicClient.Resource(resource).Namespace(namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return resultList.Items, nil
}

func (k *Klient) GetNamespaces(ctx context.Context, labels map[string]string) ([]corev1.Namespace, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetNamespaces")
	defer func() { logger.LogServiceTime(log) }()

	results, err := k.KubernetesClient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var list []corev1.Namespace
	for _, item := range results.Items {
		if helper.FilterLabels(item.Labels, labels) {
			list = append(list, item)
		}
	}

	return list, nil
}

func (k *Klient) GetStorageClass(ctx context.Context, name string) (*storagev1.StorageClass, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetStorageClass")
	defer func() { logger.LogServiceTime(log) }()

	return k.KubernetesClient.StorageV1().StorageClasses().Get(ctx, name, metav1.GetOptions{})
}

func (k *Klient) GetStorageClasses(ctx context.Context) ([]storagev1.StorageClass, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetStorageClasses")
	defer func() { logger.LogServiceTime(log) }()

	storageClasses, err := k.KubernetesClient.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return storageClasses.Items, nil
}

func (k *Klient) GetSnapshotClass(ctx context.Context, name string) (*unstructured.Unstructured, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetSnapshotClass")
	defer func() { logger.LogServiceTime(log) }()

	resource := helper.GvrMap[model.ResourceVolumeSnapshotClass]
	return k.DynamicClient.Resource(resource).Get(ctx, name, metav1.GetOptions{})
}

func (k *Klient) GetSnapshotClasses(ctx context.Context) ([]unstructured.Unstructured, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetSnapshotClasses")
	defer func() { logger.LogServiceTime(log) }()

	resource := helper.GvrMap[model.ResourceVolumeSnapshotClass]
	resultList, err := k.DynamicClient.Resource(resource).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return resultList.Items, nil

}

func (k *Klient) GetNamespace(ctx context.Context, name string) (*corev1.Namespace, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetNamespace")
	defer func() { logger.LogServiceTime(log) }()

	return k.KubernetesClient.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
}

func (k *Klient) GetDeploymentByName(ctx context.Context, namespace, name string) (*appsv1.Deployment, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetDeploymentByName")
	defer func() { logger.LogServiceTime(log) }()

	return k.KubernetesClient.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (k *Klient) GetDeployments(ctx context.Context, namespace string, labels map[string]string) []appsv1.Deployment {
	var log = logger.GetServiceLogger(ctx, "klient.GetDeployments")
	defer func() { logger.LogServiceTime(log) }()

	results, err := k.KubernetesClient.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil
	}

	var list []appsv1.Deployment
	for _, item := range results.Items {
		if helper.FilterLabels(item.Labels, labels) {
			list = append(list, item)
		}
	}

	return list
}

func (k *Klient) GetPodByName(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetPodByName")
	defer func() { logger.LogServiceTime(log) }()

	result, err := k.KubernetesClient.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (k *Klient) GetPods(ctx context.Context, namespace string, labels map[string]string) []corev1.Pod {
	var log = logger.GetServiceLogger(ctx, "klient.GetPods")
	defer func() { logger.LogServiceTime(log) }()

	results, err := k.KubernetesClient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return []corev1.Pod{}
	}

	var list []corev1.Pod
	for _, item := range results.Items {
		if helper.FilterLabels(item.Labels, labels) {
			list = append(list, item)
		}
	}

	return list
}

func (k *Klient) GetServiceByName(ctx context.Context, namespace, name string) (*corev1.Service, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetServiceByName")
	defer func() { logger.LogServiceTime(log) }()

	return k.KubernetesClient.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (k *Klient) GetServices(ctx context.Context, namespace string, labels map[string]string) ([]corev1.Service, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetServices")
	defer func() { logger.LogServiceTime(log) }()

	results, err := k.KubernetesClient.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var list []corev1.Service
	for _, item := range results.Items {
		if helper.FilterLabels(item.Labels, labels) {
			list = append(list, item)
		}
	}

	return list, nil
}

func (k *Klient) GetSecretByName(ctx context.Context, namespace, name string) (*corev1.Secret, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetSecretByName")
	defer func() { logger.LogServiceTime(log) }()

	return k.KubernetesClient.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (k *Klient) GetSecrets(ctx context.Context, namespace string, labels map[string]string) ([]corev1.Secret, error) {
	var log = logger.GetServiceLogger(ctx, "klient.Secrets")
	defer func() { logger.LogServiceTime(log) }()

	results, err := k.KubernetesClient.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Errorf("Found no secrets in %s- %s", namespace, err)
		return nil, err
	}

	log.Debugf("Found %d secrets in %s", len(results.Items), namespace)
	var list []corev1.Secret
	for _, item := range results.Items {
		if helper.FilterLabels(item.Labels, labels) {
			list = append(list, item)
		}
	}

	return list, nil
}

func (k *Klient) GetConfigMapByName(ctx context.Context, namespace, name string) (*corev1.ConfigMap, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetConfigMapByName")
	defer func() { logger.LogServiceTime(log) }()

	return k.KubernetesClient.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (k *Klient) GetConfigMaps(ctx context.Context, namespace string, labels map[string]string) ([]corev1.ConfigMap, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetConfigMaps")
	defer func() { logger.LogServiceTime(log) }()

	results, err := k.KubernetesClient.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var list []corev1.ConfigMap
	for _, item := range results.Items {
		if helper.FilterLabels(item.Labels, labels) {
			list = append(list, item)
		}
	}

	return list, nil
}

func (k *Klient) GetPersistentVolumeByName(ctx context.Context, name string) (*corev1.PersistentVolume, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetPersistentVolumeByName")
	defer func() { logger.LogServiceTime(log) }()

	return k.KubernetesClient.CoreV1().PersistentVolumes().Get(ctx, name, metav1.GetOptions{})
}

func (k *Klient) GetPersistentVolumes(ctx context.Context) ([]corev1.PersistentVolume, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetPersistentVolumes")
	defer func() { logger.LogServiceTime(log) }()

	pvs, err := k.KubernetesClient.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return pvs.Items, nil
}

func (k *Klient) GetPersistentVolumeClaimByName(ctx context.Context, namespace, name string) (*corev1.PersistentVolumeClaim, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetPersistentVolumeClaimByName")
	defer func() { logger.LogServiceTime(log) }()

	return k.KubernetesClient.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (k *Klient) GetPersistentVolumeClaims(ctx context.Context, namespace string, labels map[string]string) ([]corev1.PersistentVolumeClaim, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetPersistentVolumeClaims")
	defer func() { logger.LogServiceTime(log) }()

	results, err := k.KubernetesClient.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var list []corev1.PersistentVolumeClaim
	for _, item := range results.Items {
		if helper.FilterLabels(item.Labels, labels) {
			list = append(list, item)
		}
	}

	return list, nil
}

func (k *Klient) GenerateKubernetesObjects(ctx context.Context, specArr []string) ([]*unstructured.Unstructured, []*schema.GroupVersionKind, error) {

	var log = logger.GetServiceLogger(ctx, "klient.GenerateKubernetesObjects")
	defer func() { logger.LogServiceTime(log) }()

	var objects []*unstructured.Unstructured
	var gvks []*schema.GroupVersionKind
	for index, spec := range specArr {
		object, gvk, err := helper.DecodeFromYaml(spec)
		if err != nil {
			return nil, nil, err
		}
		log.Infof("%d. Generated object %s", index+1, object.GetName())

		objects = append(objects, object)
		gvks = append(gvks, gvk)
	}

	return objects, gvks, nil
}

func (k *Klient) GetVolumeSnapshot(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetVolumeSnapshot")
	defer func() { logger.LogServiceTime(log) }()

	return k.GetDynamicResource(ctx, namespace, name, helper.GvrMap[model.ResourceVolumeSnapshot])
}

func (k *Klient) GetVolumeSnapshots(ctx context.Context, namespace string, labels map[string]string) []unstructured.Unstructured {
	var log = logger.GetServiceLogger(ctx, "klient.GetVolumeSnapshots")
	defer func() { logger.LogServiceTime(log) }()

	var snapshots []unstructured.Unstructured

	items, err := k.GetDynamicResourceList(ctx, namespace, helper.GvrMap[model.ResourceVolumeSnapshot])
	if err != nil {
		log.Errorf("No volumesnapshots found in namespace %s - %s", namespace, err)
		return snapshots
	}

	for _, item := range items {
		if helper.FilterLabels(item.GetLabels(), labels) {
			snapshots = append(snapshots, item)
		}
	}

	return snapshots
}

func (k *Klient) GetSnapshotSchedule(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetSnapshotSchedule")
	defer func() { logger.LogServiceTime(log) }()

	return k.GetDynamicResource(ctx, namespace, name, helper.GvrMap[model.ResourceSnapshotSchedule])
}

func (k *Klient) GetSnapshotSchedules(ctx context.Context, namespace string, labels map[string]string) []unstructured.Unstructured {
	var log = logger.GetServiceLogger(ctx, "klient.GetSnapshotSchedules")
	defer func() { logger.LogServiceTime(log) }()

	var schedules []unstructured.Unstructured
	items, err := k.GetDynamicResourceList(ctx, namespace, helper.GvrMap[model.ResourceSnapshotSchedule])
	if err != nil {
		log.Errorf("No schedules found in %s - %s", namespace, err)
		return schedules
	}

	for _, item := range items {
		if helper.FilterLabels(item.GetLabels(), labels) {
			schedules = append(schedules, item)
		}
	}

	return schedules
}

func (k *Klient) GetIngress(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error) {
	var log = logger.GetServiceLogger(ctx, "klient.GetIngress")
	defer func() { logger.LogServiceTime(log) }()

	return k.GetDynamicResource(ctx, namespace, name, helper.GvrMap[model.ResourceHTTPProxy])
}

func (k *Klient) GetIngresses(ctx context.Context, namespace string, labels map[string]string) []unstructured.Unstructured {

	var log = logger.GetServiceLogger(ctx, "klient.GetIngresses")
	defer func() { logger.LogServiceTime(log) }()

	var ingresses []unstructured.Unstructured
	items, err := k.GetDynamicResourceList(ctx, namespace, helper.GvrMap[model.ResourceHTTPProxy])
	if err != nil {
		log.Errorf("No schedules found in %s - %s", namespace, err)
		return ingresses
	}

	for _, item := range items {
		if helper.FilterLabels(item.GetLabels(), labels) {
			ingresses = append(ingresses, item)
		}
	}

	return ingresses
}
