package interfaces

import (
	"context"
	"github.com/zbitech/controller/pkg/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type KlientFactoryIF interface {
	Init(ctx context.Context, repoSvc RepositoryServiceIF) error
	GetZBIClient() ZBIClientIF
	StartMonitor()
	StopMonitor()
}

type KlientMonitorIF interface {
	AddInformer(rType model.ResourceObjectType)
	Start()
	Stop()
}

type ZBIClientIF interface {
	GetProjects(ctx context.Context) ([]model.Project, error)
	GetProject(ctx context.Context, project string) (*model.Project, error)
	CreateProject(ctx context.Context, project *model.Project) error
	RepairProject(ctx context.Context, project *model.Project) error
	DeleteProject(ctx context.Context, project *model.Project, instances []model.Instance) error
	GetProjectResources(ctx context.Context, project string) ([]model.KubernetesResource, error)
	GetProjectResource(ctx context.Context, project, resourceName string, resourceType model.ResourceObjectType) (*model.KubernetesResource, error)

	CreateInstance(ctx context.Context, project *model.Project, instance *model.Instance) error
	GetAllInstances(ctx context.Context, project *model.Project) ([]model.Instance, error)
	GetInstances(ctx context.Context, project *model.Project, instances []string) ([]model.Instance, error)
	GetInstance(ctx context.Context, project *model.Project, instance string) (*model.Instance, error)
	GetInstanceResources(ctx context.Context, project *model.Project, instance string) (*model.KubernetesResources, error)
	GetInstanceResource(ctx context.Context, project *model.Project, instance, resourceName string, resourceType model.ResourceObjectType) (*model.KubernetesResource, error)
	DeleteInstanceResource(ctx context.Context, project *model.Project, instance *model.Instance, resourceName string, resourceType model.ResourceObjectType) error
	UpdateInstance(ctx context.Context, project *model.Project, instance *model.Instance) error
	DeleteInstance(ctx context.Context, project *model.Project, instance *model.Instance) error
	RepairInstance(ctx context.Context, project *model.Project, instance *model.Instance) error
	StopInstance(ctx context.Context, project *model.Project, instance *model.Instance) error
	StartInstance(ctx context.Context, project *model.Project, instance *model.Instance) error
	RotateInstanceCredentials(ctx context.Context, project *model.Project, instance *model.Instance) error
	CreateSnapshot(ctx context.Context, project *model.Project, instance *model.Instance) error
	CreateSnapshotSchedule(ctx context.Context, project *model.Project, instance *model.Instance, schedule model.SnapshotScheduleType) error
}

type KlientIF interface {
	GetKubernetesClient() kubernetes.Interface
	GetDynamicClient() dynamic.Interface

	DeleteDynamicResource(ctx context.Context, namespace, name string, resource schema.GroupVersionResource) error
	DeleteNamespace(ctx context.Context, namespace string) error

	ApplyResource(ctx context.Context, object *unstructured.Unstructured) (*model.KubernetesResource, error)
	ApplyResources(ctx context.Context, objects []unstructured.Unstructured) ([]model.KubernetesResource, error)

	DeleteResource(ctx context.Context, resource *model.KubernetesResource) error
	DeleteResources(ctx context.Context, resource []model.KubernetesResource) ([]model.KubernetesResource, error)

	GetDynamicResource(ctx context.Context, namespace, name string, resource schema.GroupVersionResource) (*unstructured.Unstructured, error)
	GetDynamicResourceList(ctx context.Context, namespace string, resource schema.GroupVersionResource) ([]unstructured.Unstructured, error)

	GetNamespace(ctx context.Context, name string) (*corev1.Namespace, error)
	GetNamespaces(ctx context.Context, labels map[string]string) ([]corev1.Namespace, error)

	GetStorageClass(ctx context.Context, name string) (*storagev1.StorageClass, error)
	GetStorageClasses(ctx context.Context) ([]storagev1.StorageClass, error)

	GetSnapshotClass(ctx context.Context, name string) (*unstructured.Unstructured, error)
	GetSnapshotClasses(ctx context.Context) ([]unstructured.Unstructured, error)

	GetDeploymentByName(ctx context.Context, namespace, name string) (*appsv1.Deployment, error)
	GetDeployments(ctx context.Context, namespace string, labels map[string]string) []appsv1.Deployment

	GetPodByName(ctx context.Context, namespace, name string) (*corev1.Pod, error)
	GetPods(ctx context.Context, namespace string, labels map[string]string) []corev1.Pod

	GetServiceByName(ctx context.Context, namespace, name string) (*corev1.Service, error)
	GetServices(ctx context.Context, namespace string, labels map[string]string) ([]corev1.Service, error)

	GetSecretByName(ctx context.Context, namespace, name string) (*corev1.Secret, error)
	GetSecrets(ctx context.Context, namespace string, labels map[string]string) ([]corev1.Secret, error)

	GetConfigMapByName(ctx context.Context, namespace, name string) (*corev1.ConfigMap, error)
	GetConfigMaps(ctx context.Context, namespace string, labels map[string]string) ([]corev1.ConfigMap, error)

	GetPersistentVolumeByName(ctx context.Context, name string) (*corev1.PersistentVolume, error)
	GetPersistentVolumes(ctx context.Context) ([]corev1.PersistentVolume, error)

	GetPersistentVolumeClaimByName(ctx context.Context, namespace, name string) (*corev1.PersistentVolumeClaim, error)
	GetPersistentVolumeClaims(ctx context.Context, namespace string, labels map[string]string) ([]corev1.PersistentVolumeClaim, error)

	GetVolumeSnapshot(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error)
	GetVolumeSnapshots(ctx context.Context, namespace string, labels map[string]string) []unstructured.Unstructured

	GetSnapshotSchedule(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error)
	GetSnapshotSchedules(ctx context.Context, namespace string, labels map[string]string) []unstructured.Unstructured

	GetIngress(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error)
	GetIngresses(ctx context.Context, namespace string, labels map[string]string) []unstructured.Unstructured
}
