package client

import (
	"context"
	"github.com/zbitech/controller/pkg/interfaces"
	"github.com/zbitech/controller/pkg/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	dynfake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

type FakeKlient struct {
	KubernetesClient kubernetes.Interface
	DynamicClient    dynamic.Interface

	FakeDeleteDynamicResource          func(ctx context.Context, namespace, name string, resource schema.GroupVersionResource) error
	FakeDeleteNamespace                func(ctx context.Context, namespace string) error
	FakeApplyResource                  func(ctx context.Context, object *unstructured.Unstructured) (*model.KubernetesResource, error)
	FakeApplyResources                 func(ctx context.Context, objects []unstructured.Unstructured) ([]model.KubernetesResource, error)
	FakeDeleteResource                 func(ctx context.Context, object *model.KubernetesResource) error
	FakeDeleteResources                func(ctx context.Context, objects []model.KubernetesResource) ([]model.KubernetesResource, error)
	FakeGetDynamicResource             func(ctx context.Context, namespace, name string, resource schema.GroupVersionResource) (*unstructured.Unstructured, error)
	FakeGetDynamicResourceList         func(ctx context.Context, namespace string, resource schema.GroupVersionResource) ([]unstructured.Unstructured, error)
	FakeGetNamespace                   func(ctx context.Context, name string) (*corev1.Namespace, error)
	FakeGetNamespaces                  func(ctx context.Context, labels map[string]string) ([]corev1.Namespace, error)
	FakeGetStorageClass                func(ctx context.Context, name string) (*storagev1.StorageClass, error)
	FakeGetStorageClasses              func(ctx context.Context) ([]storagev1.StorageClass, error)
	FakeGetSnapshotClass               func(ctx context.Context, name string) (*unstructured.Unstructured, error)
	FakeGetSnapshotClasses             func(ctx context.Context) ([]unstructured.Unstructured, error)
	FakeGetDeploymentByName            func(ctx context.Context, namespace, name string) (*appsv1.Deployment, error)
	FakeGetDeployments                 func(ctx context.Context, namespace string, labels map[string]string) []appsv1.Deployment
	FakeGetPodByName                   func(ctx context.Context, namespace, name string) (*corev1.Pod, error)
	FakeGetPods                        func(ctx context.Context, namespace string, labels map[string]string) []corev1.Pod
	FakeGetServiceByName               func(ctx context.Context, namespace, name string) (*corev1.Service, error)
	FakeGetServices                    func(ctx context.Context, namespace string, labels map[string]string) ([]corev1.Service, error)
	FakeGetSecretByName                func(ctx context.Context, namespace, name string) (*corev1.Secret, error)
	FakeGetSecrets                     func(ctx context.Context, namespace string, labels map[string]string) ([]corev1.Secret, error)
	FakeGetConfigMapByName             func(ctx context.Context, namespace, name string) (*corev1.ConfigMap, error)
	FakeGetConfigMaps                  func(ctx context.Context, namespace string, labels map[string]string) ([]corev1.ConfigMap, error)
	FakeGetPersistentVolumeByName      func(ctx context.Context, name string) (*corev1.PersistentVolume, error)
	FakeGetPersistentVolumes           func(ctx context.Context) ([]corev1.PersistentVolume, error)
	FakeGetPersistentVolumeClaimByName func(ctx context.Context, namespace, name string) (*corev1.PersistentVolumeClaim, error)
	FakeGetPersistentVolumeClaims      func(ctx context.Context, namespace string, labels map[string]string) ([]corev1.PersistentVolumeClaim, error)
	FakeGetVolumeSnapshot              func(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error)
	FakeGetVolumeSnapshots             func(ctx context.Context, namespace string, labels map[string]string) []unstructured.Unstructured
	FakeGetSnapshotSchedule            func(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error)
	FakeGetSnapshotSchedules           func(ctx context.Context, namespace string, labels map[string]string) []unstructured.Unstructured
	FakeGetIngress                     func(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error)
	FakeGetIngresses                   func(ctx context.Context, namespace string, labels map[string]string) []unstructured.Unstructured
	FakeGetKubernetesClient            func() kubernetes.Interface
	FakeGetDynamicClient               func() dynamic.Interface
}

func NewFakeKlient(ctx context.Context) (interfaces.KlientIF, error) {

	return &FakeKlient{
		KubernetesClient: k8sfake.NewSimpleClientset(),
		DynamicClient:    dynfake.NewSimpleDynamicClient(runtime.NewScheme()),
	}, nil
}

func (f FakeKlient) GetKubernetesClient() kubernetes.Interface {
	return f.FakeGetKubernetesClient()
}

func (f FakeKlient) GetDynamicClient() dynamic.Interface {
	return f.FakeGetDynamicClient()
}

func (f FakeKlient) GetResource(object *unstructured.Unstructured) *model.KubernetesResource {
	//TODO implement me
	panic("implement me")
}

func (f FakeKlient) GenerateKubernetesObjects(ctx context.Context, spec_arr []string) ([]*unstructured.Unstructured, []*schema.GroupVersionKind, error) {
	//TODO implement me
	panic("implement me")
}

func (f FakeKlient) DeleteDynamicResource(ctx context.Context, namespace, name string, resource schema.GroupVersionResource) error {
	return f.FakeDeleteDynamicResource(ctx, namespace, name, resource)
}

func (f FakeKlient) DeleteNamespace(ctx context.Context, namespace string) error {
	return f.FakeDeleteNamespace(ctx, namespace)
}

func (f FakeKlient) ApplyResource(ctx context.Context, object *unstructured.Unstructured) (*model.KubernetesResource, error) {
	return f.FakeApplyResource(ctx, object)
}

func (f FakeKlient) ApplyResources(ctx context.Context, objects []unstructured.Unstructured) ([]model.KubernetesResource, error) {
	return f.FakeApplyResources(ctx, objects)
}

func (f FakeKlient) DeleteResource(ctx context.Context, object *model.KubernetesResource) error {
	return f.FakeDeleteResource(ctx, object)
}

func (f FakeKlient) DeleteResources(ctx context.Context, objects []model.KubernetesResource) ([]model.KubernetesResource, error) {
	return f.FakeDeleteResources(ctx, objects)
}

func (f FakeKlient) GetDynamicResource(ctx context.Context, namespace, name string, resource schema.GroupVersionResource) (*unstructured.Unstructured, error) {
	return f.FakeGetDynamicResource(ctx, namespace, name, resource)
}

func (f FakeKlient) GetDynamicResourceList(ctx context.Context, namespace string, resource schema.GroupVersionResource) ([]unstructured.Unstructured, error) {
	return f.FakeGetDynamicResourceList(ctx, namespace, resource)
}

func (f FakeKlient) GetNamespace(ctx context.Context, name string) (*corev1.Namespace, error) {
	return f.FakeGetNamespace(ctx, name)
}

func (f FakeKlient) GetNamespaces(ctx context.Context, labels map[string]string) ([]corev1.Namespace, error) {
	return f.FakeGetNamespaces(ctx, labels)
}

func (f FakeKlient) GetStorageClass(ctx context.Context, name string) (*storagev1.StorageClass, error) {
	return f.FakeGetStorageClass(ctx, name)
}

func (f FakeKlient) GetStorageClasses(ctx context.Context) ([]storagev1.StorageClass, error) {
	return f.FakeGetStorageClasses(ctx)
}

func (f FakeKlient) GetSnapshotClass(ctx context.Context, name string) (*unstructured.Unstructured, error) {
	return f.FakeGetSnapshotClass(ctx, name)
}

func (f FakeKlient) GetSnapshotClasses(ctx context.Context) ([]unstructured.Unstructured, error) {
	return f.FakeGetSnapshotClasses(ctx)
}

func (f FakeKlient) GetDeploymentByName(ctx context.Context, namespace, name string) (*appsv1.Deployment, error) {
	return f.FakeGetDeploymentByName(ctx, namespace, name)
}

func (f FakeKlient) GetDeployments(ctx context.Context, namespace string, labels map[string]string) []appsv1.Deployment {
	return f.FakeGetDeployments(ctx, namespace, labels)
}

func (f FakeKlient) GetPodByName(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	return f.FakeGetPodByName(ctx, namespace, name)
}

func (f FakeKlient) GetPods(ctx context.Context, namespace string, labels map[string]string) []corev1.Pod {
	return f.FakeGetPods(ctx, namespace, labels)
}

func (f FakeKlient) GetServiceByName(ctx context.Context, namespace, name string) (*corev1.Service, error) {
	return f.FakeGetServiceByName(ctx, namespace, name)
}

func (f FakeKlient) GetServices(ctx context.Context, namespace string, labels map[string]string) ([]corev1.Service, error) {
	return f.FakeGetServices(ctx, namespace, labels)
}

func (f FakeKlient) GetSecretByName(ctx context.Context, namespace, name string) (*corev1.Secret, error) {
	return f.FakeGetSecretByName(ctx, namespace, name)
}

func (f FakeKlient) GetSecrets(ctx context.Context, namespace string, labels map[string]string) ([]corev1.Secret, error) {
	return f.FakeGetSecrets(ctx, namespace, labels)
}

func (f FakeKlient) GetConfigMapByName(ctx context.Context, namespace, name string) (*corev1.ConfigMap, error) {
	return f.FakeGetConfigMapByName(ctx, namespace, name)
}

func (f FakeKlient) GetConfigMaps(ctx context.Context, namespace string, labels map[string]string) ([]corev1.ConfigMap, error) {
	return f.FakeGetConfigMaps(ctx, namespace, labels)
}

func (f FakeKlient) GetPersistentVolumeByName(ctx context.Context, name string) (*corev1.PersistentVolume, error) {
	return f.FakeGetPersistentVolumeByName(ctx, name)
}

func (f FakeKlient) GetPersistentVolumes(ctx context.Context) ([]corev1.PersistentVolume, error) {
	return f.FakeGetPersistentVolumes(ctx)
}

func (f FakeKlient) GetPersistentVolumeClaimByName(ctx context.Context, namespace, name string) (*corev1.PersistentVolumeClaim, error) {
	return f.FakeGetPersistentVolumeClaimByName(ctx, namespace, name)
}

func (f FakeKlient) GetPersistentVolumeClaims(ctx context.Context, namespace string, labels map[string]string) ([]corev1.PersistentVolumeClaim, error) {
	return f.FakeGetPersistentVolumeClaims(ctx, namespace, labels)
}

func (f FakeKlient) GetVolumeSnapshot(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error) {
	return f.FakeGetVolumeSnapshot(ctx, namespace, name)
}

func (f FakeKlient) GetVolumeSnapshots(ctx context.Context, namespace string, labels map[string]string) []unstructured.Unstructured {
	return f.FakeGetVolumeSnapshots(ctx, namespace, labels)
}

func (f FakeKlient) GetSnapshotSchedule(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error) {
	return f.FakeGetSnapshotSchedule(ctx, namespace, name)
}

func (f FakeKlient) GetSnapshotSchedules(ctx context.Context, namespace string, labels map[string]string) []unstructured.Unstructured {
	return f.FakeGetSnapshotSchedules(ctx, namespace, labels)
}

func (f FakeKlient) GetIngress(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error) {
	return f.FakeGetIngress(ctx, namespace, name)
}

func (f FakeKlient) GetIngresses(ctx context.Context, namespace string, labels map[string]string) []unstructured.Unstructured {
	return f.FakeGetIngresses(ctx, namespace, labels)
}
