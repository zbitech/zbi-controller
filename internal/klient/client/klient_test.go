package client

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	fake_zbi "github.com/zbitech/controller/fake-zbi"
	"github.com/zbitech/controller/fake-zbi/data"
	"github.com/zbitech/controller/internal/helper"
	"github.com/zbitech/controller/internal/utils"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dynfake "k8s.io/client-go/dynamic/fake"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"testing"
)

var (
	//	initialized = test.InitTestConfig()

	PodYaml        = "apiVersion: v1\nkind: Pod\nmetadata:\n  name: static-web\n  namespace: default\n  labels:\n    role: myrole\nspec:\n  containers:\n    - name: web\n      image: nginx\n      ports:\n        - name: web\n          containerPort: 80\n          protocol: TCP"
	DeploymentYAML = "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: nginx\n  namespace: default\n  labels:\n    app: nginx\nspec:\n  selector:\n    matchLabels:\n      app: nginx\n  template:\n    metadata:\n      labels:\n        app: nginx\n    spec:\n      containers:\n      - name: webserver\n        image: nginx\n        imagePullPolicy: Always\n"

	NamespaceJSON  = "{\n  \"apiVersion\": \"v1\",\n  \"kind\": \"Namespace\",\n  \"metadata\": {\n    \"name\": \"project\"\n  }\n}\n"
	PodJSON        = "{\n  \"apiVersion\": \"v1\",\n  \"kind\": \"Pod\",\n  \"metadata\": {\n    \"labels\": {\n      \"app\": \"nginx\"\n    },\n    \"name\": \"static-web\",\n    \"namespace\": \"default\"\n  },\n  \"spec\": {\n    \"containers\": [\n      {\n        \"name\": \"web\",\n        \"image\": \"nginx\",\n        \"ports\":[\n          {\n            \"name\": \"web\",\n            \"containerPort\": 80,\n            \"protocol\": \"TCP\"\n          }\n        ]\n      }\n    ]\n  }\n}\n"
	DeploymentJSON = "{\n  \"apiVersion\": \"apps/v1\",\n  \"kind\": \"Deployment\",\n  \"metadata\": {\n    \"labels\": {\n      \"app\": \"nginx\"\n    },\n    \"name\": \"nginx-app-dyn\",\n    \"namespace\": \"default\"\n  },\n  \"spec\": {\n    \"selector\": {\n      \"matchLabels\": {\n        \"app\": \"nginx\"\n      }\n    },\n    \"template\": {\n      \"metadata\": {\n        \"labels\": {\n          \"app\": \"nginx\"\n        }\n      },\n      \"spec\": {\n        \"containers\": [\n          {\n            \"image\": \"nginx\",\n            \"imagePullPolicy\": \"Always\",\n            \"name\": \"webserver\"\n          }\n        ]\n      }\n    }\n  }\n}\n"
	ServiceJSON    = "{\"apiVersion\": \"v1\", \"kind\": \"Service\", \"metadata\": {\"labels\": {\"app\": \"nginx\"}, \"name\": \"nginx-svc\", \"namespace\": \"default\"}, \"spec\": {\"ports\": [{\"name\": \"http\",\"port\": 80,\"targetPort\": 80}],\"selector\": {\"app\": \"nginx\"}}}\n"

	NamespaceGVR  = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}
	PodGVR        = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	DeploymentGVR = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	ServiceGVR    = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
)

func FakeNewKlient(ctx context.Context) (*Klient, error) {

	var log = logger.GetLogger(ctx)
	cfg, err := NewRestConfig(ctx)
	if err != nil {
		log.Errorf("Unable to create kubernetes configuration - %s", err)
		// return nil, errs.KubernetesError
		return nil, errors.New("kubernetes error")
	}

	kubernetesClient := k8sfake.NewSimpleClientset()
	//	dynamicClient := dynfake.NewSimpleDynamicClientWithCustomListKinds()
	// pass in object of custom type
	dynamicClient := dynfake.NewSimpleDynamicClient(runtime.NewScheme())

	//	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(kubernetesClient.Discovery()))

	return &Klient{
		KubernetesClient: kubernetesClient,
		DynamicClient:    dynamicClient,
		//		Mapper:           mapper,
		restConfg: cfg,
	}, nil
}

func TestKlient_NewRestConfig(t *testing.T) {
	ctx := fake_zbi.InitContext()
	r, err := NewRestConfig(ctx)
	assert.NoError(t, err, "Failed to create rest config - %s", err)
	assert.NotNilf(t, r, "Failed to create rest config")
}

func TestKlient_NewKlient(t *testing.T) {
	ctx := fake_zbi.InitContext()
	k, err := NewKlient(ctx)
	assert.NoError(t, err, "Failed to create kubernetes client - %s", err)
	assert.NotNilf(t, k, "Failed to create kubernetes client")
}

func TestKlient_ApplyResources(t *testing.T) {
	ctx := fake_zbi.InitContext()
	k, err := NewKlient(ctx)

	helper.Config.LoadConfig(ctx)
	helper.Config.LoadTemplates(ctx)

	zcash, err := data.GetZcashResources()
	assert.NoError(t, err)
	zcashResources, err := k.ApplyResources(ctx, zcash)
	assert.NoError(t, err)
	assert.NotNil(t, zcashResources)

	//for _, resource := range zcashResources {
	//	t.Log()
	//}

	//lwd, err := data.GetZcashResources()
	//assert.NoError(t, err)
	//lwdResources, err := k.ApplyResources(ctx, lwd)
	//assert.NoError(t, err)
	//assert.NotNil(t, lwdResources)
}

func TestKlient_DeleteResource(t *testing.T) {
	ctx := fake_zbi.InitContext()
	k, err := NewKlient(ctx)
	assert.NoError(t, err)

	helper.Config.LoadConfig(ctx)
	helper.Config.LoadTemplates(ctx)

	objects, err := data.GetResource(model.ResourceConfigMap)
	assert.NoError(t, err)
	assert.NotNil(t, objects)

	resources, err := k.ApplyResources(ctx, objects)
	assert.NoError(t, err)
	for i := len(resources) - 1; i >= 0; i-- {
		t.Logf("Deleting resource: %s", utils.MarshalObject(resources[i]))
		err = k.DeleteResource(ctx, &resources[i])
		assert.NoError(t, err)
	}
}

//func Test_DeleteDynamicResource(t *testing.T) {
//
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err, "Failed to create kubernetes client - %s", err)
//	assert.NotNilf(t, k, "Failed to create kubernetes client")
//
//	var obj unstructured.Unstructured
//
//	err = json.Unmarshal([]byte(DeploymentJSON), &obj)
//	assert.NoError(t, err, "Failed to generate deployment object - %s", err)
//	var robj = model.ResourceObject{Unstructured: obj, Properties: nil}
//
//	_, err = k.ApplyResource(ctx, &robj)
//	assert.NoError(t, err, "Failed to create deployment resource - %s", err)
//
//	err = k.DeleteDynamicResource(ctx, obj.GetNamespace(), obj.GetName(), DeploymentGVR)
//	assert.NoError(t, err, "Failed to delete %s resource after creation", obj.GetName())
//}
//
//func Test_DeleteNamespace(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err, "Failed to create kubernetes client - %s", err)
//	assert.NotNilf(t, k, "Failed to create kubernetes client")
//
//	var obj unstructured.Unstructured
//	err = json.Unmarshal([]byte(NamespaceJSON), &obj)
//	assert.NoError(t, err, "Failed to generate namespace object - %s", err)
//	var robj = model.ResourceObject{Unstructured: obj, Properties: nil}
//
//	_, err = k.ApplyResource(ctx, &robj)
//	assert.NoError(t, err, "Failed to create deployment object - %s", err)
//
//	namespace := obj.GetName()
//	err = k.DeleteNamespace(ctx, namespace)
//	assert.NoError(t, err, "Failed to delete namespace %s", obj.GetName())
//}
//
//func Test_GetDynamicResource(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err, "Failed to create kubernetes client - %s", err)
//	assert.NotNilf(t, k, "Failed to create kubernetes client")
//
//	var object unstructured.Unstructured
//	err = json.Unmarshal([]byte(DeploymentJSON), &object)
//	assert.NoError(t, err, "Failed to generate object - %s", err)
//	var robj = model.ResourceObject{Unstructured: object, Properties: nil}
//
//	_, err = k.ApplyResource(ctx, &robj)
//	assert.NoError(t, err, "Failed to create deployment object - %s", err)
//
//	result, err := k.GetDynamicResource(ctx, "default", object.GetName(), DeploymentGVR)
//	assert.NoError(t, err, "Failed to get deployment resource - %s", err)
//	assert.NotNilf(t, result, "Failed to get deployment resource")
//
//	err = k.DeleteDynamicResource(ctx, object.GetNamespace(), object.GetName(), DeploymentGVR)
//}
//
//func Test_GetDynamicResourceList(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err, "Failed to create kubernetes client - %s", err)
//	assert.NotNilf(t, k, "Failed to create kubernetes client")
//
//	namespace := "kube-system"
//	resourceList, err := k.GetDynamicResourceList(ctx, namespace, PodGVR)
//	assert.NoError(t, err, "Error getting pods from %s - %s", namespace, err)
//	assert.NotNilf(t, resourceList, "Failed to get list of pods from namespace %s", namespace)
//}
//
//func Test_GetNamespaces(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err, "Failed to create kubernetes client - %s", err)
//	assert.NotNilf(t, k, "Failed to create kubernetes client")
//
//	namespaces, err := k.GetNamespaces(ctx, nil)
//	assert.NoError(t, err, "Got an error while retrieving namespaces")
//	assert.NotNilf(t, namespaces, "Unable to get namespaces")
//}
//
//func Test_GetStorageClass(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err, "Failed to create kubernetes client - %s", err)
//	assert.NotNilf(t, k, "Failed to create kubernetes client")
//
//	className := "standard"
//	class, err := k.GetStorageClass(ctx, className)
//	assert.NoError(t, err, "Got an error while retrieving %s storage classes - %s", className, err)
//	assert.NotNilf(t, class, "Failed to get %s storage classes", className)
//}
//
//func Test_GetStorageClasses(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err, "Failed to create kubernetes client - %s", err)
//	assert.NotNilf(t, k, "Failed to create kubernetes client")
//
//	classes, err := k.GetStorageClasses(ctx)
//	assert.NoError(t, err, "Got an error while retrieving storage classes - %s", err)
//	assert.NotNilf(t, classes, "Failed to get storage classes")
//}
//
//func Test_GetNamespace(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err, "Failed to create kubernetes client - %s", err)
//	assert.NotNilf(t, k, "Failed to create kubernetes client")
//
//	name := "default"
//	namespace, err := k.GetNamespace(ctx, name)
//	assert.NoError(t, err, "Got an error while retrieving %s namespace - %s", name, err)
//	assert.NotNilf(t, namespace, "Failed to get %s namespace", name)
//}
//
//func Test_GetDeployments(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err, "Failed to create kubernetes client - %s", err)
//	assert.NotNilf(t, k, "Failed to create kubernetes client")
//
//	namespace := "kube-system"
//	deployments := k.GetDeployments(ctx, namespace, nil)
//	assert.NoError(t, err, "Got an error while retrieving deployments from %s - %s", namespace, err)
//	assert.NotNilf(t, deployments, "Failed to get deployments from %s", namespace)
//}
//
//func Test_GetDeploymentByName(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err, "Failed to create kubernetes client - %s", err)
//	assert.NotNilf(t, k, "Failed to create kubernetes client")
//
//	namespace := "kube-system"
//	name := "coredns"
//	deployment, err := k.GetDeploymentByName(ctx, namespace, name)
//	assert.NoError(t, err, "Got an error while retrieving deployment %s from %s - %s", name, namespace, err)
//	assert.NotNilf(t, deployment, "Failed to get deployment %s from %s", name, namespace)
//
//	var obj unstructured.Unstructured
//	err = json.Unmarshal([]byte(DeploymentJSON), &obj)
//	assert.NoError(t, err, "Failed to generate deployment object - %s", err)
//	var robj = model.ResourceObject{Unstructured: obj, Properties: nil}
//
//	_, err = k.ApplyResource(ctx, &robj)
//	assert.NoError(t, err, "Failed to create deployment resource - %s", err)
//
//	namespace = obj.GetNamespace()
//	name = obj.GetName()
//	deployment, err = k.GetDeploymentByName(ctx, namespace, name)
//	assert.NoError(t, err, "Got an error while retrieving deployment %s from %s - %s", name, namespace, err)
//	assert.NotNilf(t, deployment, "Failed to get deployment %s from %s", name, namespace)
//
//	k.DeleteDynamicResource(ctx, namespace, name, DeploymentGVR)
//}
//
//func Test_GetPods(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err, "Failed to create kubernetes client - %s", err)
//	assert.NotNilf(t, k, "Failed to create kubernetes client")
//
//	namespace := "kube-system"
//	pods, err := k.GetPods(ctx, namespace, nil)
//	assert.NoError(t, err, "Got an error while retrieving pods from %s - %s", namespace, err)
//	assert.NotNilf(t, pods, "Failed to get pods from %s", namespace)
//}
//
//func Test_GetPodByName(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err, "Failed to create kubernetes client - %s", err)
//	assert.NotNilf(t, k, "Failed to create kubernetes client")
//
//	//create pod
//	var obj unstructured.Unstructured
//	err = json.Unmarshal([]byte(PodJSON), &obj)
//	assert.NoError(t, err, "Failed to generate pod object - %s", err)
//	var robj = model.ResourceObject{Unstructured: obj, Properties: nil}
//
//	_, err = k.ApplyResource(ctx, &robj)
//	assert.NoError(t, err, "Failed to create pod resource - %s", err)
//
//	//get pod
//	namespace := obj.GetNamespace()
//	name := obj.GetName()
//	pod, err := k.GetPodByName(ctx, namespace, name)
//	assert.NoError(t, err, "Got an error while retrieving pod %s from %s - %s", name, namespace, err)
//	assert.NotNilf(t, pod, "Failed to get pod %s from %s", name, namespace)
//
//	k.DeleteDynamicResource(ctx, namespace, name, PodGVR)
//
//	//get non-existent pod
//	pod, err = k.GetPodByName(ctx, namespace, "fake-pod")
//	assert.Errorf(t, err, "Did not expected error for non-existent pod")
//	assert.Nilf(t, pod, "Got a result instead of nil")
//}
//
//func Test_GetServices(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err, "Failed to create kubernetes client - %s", err)
//	assert.NotNilf(t, k, "Failed to create kubernetes client")
//
//	namespace := "kube-system"
//	services, err := k.GetServices(ctx, namespace, nil)
//	assert.NoError(t, err, "Got an error while retrieving services from %s - %s", namespace, err)
//	assert.NotNilf(t, services, "Failed to get services from %s", namespace)
//}
//
//func Test_GetServiceByName(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err, "Failed to create kubernetes client - %s", err)
//	assert.NotNilf(t, k, "Failed to create kubernetes client")
//
//	namespace := "kube-system"
//	name := "kube-dns"
//	service, err := k.GetServiceByName(ctx, namespace, name)
//	assert.NoError(t, err, "Got an error while retrieving pod %s from %s - %s", name, namespace, err)
//	assert.NotNilf(t, service, "Failed to get pod %s from %s", name, namespace)
//
//}
//
//func Test_GetSecretByName(t *testing.T) {
//
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err)
//	assert.NotNil(t, k)
//
//	obj, err := data.GetGenericResource(model.ResourceSecret)
//	assert.NoError(t, err)
//	assert.NotNil(t, obj)
//
//	labels := map[string]string{"platform": "zbi", "project": "project", "instance": "instance"}
//	data := map[string]string{"password": utils.Base64EncodeString("password")}
//	manager.SetResourceField(obj, "data", data)
//	manager.SetResourceField(obj, "metadata.labels", labels)
//	var robj = model.ResourceObject{Unstructured: *obj, Properties: nil}
//
//	resource, err := k.ApplyResource(ctx, &robj)
//	assert.NoError(t, err)
//	assert.NotNil(t, resource)
//
//	secret, err := k.GetSecretByName(ctx, obj.GetNamespace(), obj.GetName())
//	assert.NoError(t, err)
//	assert.NotNil(t, secret)
//	assert.Equal(t, obj.GetName(), secret.Name)
//	//	assert.Equal(t, obj.UnstructuredContent()["data"], secret.StringData)
//
//	k.DeleteDynamicResource(ctx, resource.Namespace, resource.Name, helper.GvrMap[resource.Type])
//}
//
//func Test_GetSecrets(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err)
//	assert.NotNil(t, k)
//
//	obj, err := data.GetGenericResource(model.ResourceSecret)
//	assert.NoError(t, err)
//	assert.NotNil(t, obj)
//
//	labels := map[string]string{"platform": "zbi", "project": "project", "instance": "instance"}
//	data := map[string]string{"password": utils.Base64EncodeString("password")}
//	manager.SetResourceField(obj, "data", data)
//	manager.SetResourceField(obj, "metadata.labels", labels)
//
//	var robj = model.ResourceObject{Unstructured: *obj, Properties: nil}
//	resource, err := k.ApplyResource(ctx, &robj)
//	assert.NoError(t, err)
//	assert.NotNil(t, resource)
//
//	target, err := k.GetSecrets(ctx, obj.GetNamespace(), labels)
//	assert.NoError(t, err)
//	assert.NotNil(t, target)
//	assert.Len(t, target, 1)
//
//	k.DeleteDynamicResource(ctx, resource.Namespace, resource.Name, helper.GvrMap[resource.Type])
//}
//
//func Test_GetConfigMapByName(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err)
//	assert.NotNil(t, k)
//
//	obj, err := data.GetGenericResource(model.ResourceConfigMap)
//	assert.NoError(t, err)
//	assert.NotNil(t, obj)
//
//	labels := map[string]string{"platform": "zbi", "project": "project", "instance": "instance"}
//	data := map[string]string{"username": "zbiuser"}
//	manager.SetResourceField(obj, "data", data)
//	manager.SetResourceField(obj, "metadata.labels", labels)
//
//	var robj = model.ResourceObject{Unstructured: *obj, Properties: nil}
//	resource, err := k.ApplyResource(ctx, &robj)
//	assert.NoError(t, err)
//	assert.NotNil(t, resource)
//
//	target, err := k.GetConfigMapByName(ctx, obj.GetNamespace(), obj.GetName())
//	assert.NoError(t, err)
//	assert.NotNil(t, target)
//	assert.Equal(t, obj.GetName(), target.Name)
//	assert.Equal(t, obj.UnstructuredContent()["data"], target.Data)
//
//	k.DeleteDynamicResource(ctx, resource.Namespace, resource.Name, helper.GvrMap[resource.Type])
//}
//
//func Test_GetConfigMaps(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err)
//	assert.NotNil(t, k)
//
//	obj, err := data.GetGenericResource(model.ResourceConfigMap)
//	assert.NoError(t, err)
//	assert.NotNil(t, obj)
//
//	labels := map[string]string{"platform": "zbi", "project": "project", "instance": "instance"}
//	data := map[string]string{"username": "zbiuser"}
//	manager.SetResourceField(obj, "data", data)
//	manager.SetResourceField(obj, "metadata.labels", labels)
//
//	var robj = model.ResourceObject{Unstructured: *obj, Properties: nil}
//
//	resource, err := k.ApplyResource(ctx, &robj)
//	assert.NoError(t, err)
//	assert.NotNil(t, resource)
//
//	target, err := k.GetConfigMaps(ctx, obj.GetNamespace(), labels)
//	assert.NoError(t, err)
//	assert.NotNil(t, target)
//	assert.Len(t, target, 1)
//
//	k.DeleteDynamicResource(ctx, resource.Namespace, resource.Name, helper.GvrMap[resource.Type])
//}
//
//func Test_GetPersistentVolumeByName(t *testing.T) {
//	//ctx := context.Background()
//	//k, err := NewKlient(ctx)
//	//assert.NoError(t, err)
//	//assert.NotNil(t, k)
//	//
//	//obj, err := data.GetGenericResource(ztypes.ResourcePersistentVolumeClaim)
//	//assert.NoError(t, err)
//	//assert.NotNil(t, obj)
//	//
//	//labels := map[string]string{"platform": "zbi", "project": "project", "instance": "instance"}
//	//utils.RemoveResourceField(obj, "spec.dataSource")
//	//utils.SetResourceField(obj, "metadata.labels", labels)
//	//
//	//resource, err := k.ApplyResource(ctx, obj)
//	//assert.NoError(t, err)
//	//assert.NotNil(t, resource)
//	//
//	//target, err := k.GetPersistentVolumeByName(ctx, obj.GetName())
//	//assert.NoError(t, err)
//	//assert.NotNil(t, target)
//	//assert.Equal(t, obj.GetName(), target.Name)
//
//	//	k.DeleteDynamicResource(ctx, resource.Namespace, resource.Name, helper.GvrMap[resource.Type])
//}
//
//func Test_GetPersistentVolumes(t *testing.T) {
//	//ctx := context.Background()
//	//k, err := NewKlient(ctx)
//	//assert.NoError(t, err)
//	//assert.NotNil(t, k)
//	//
//	//obj, err := data.GetGenericResource(ztypes.ResourcePersistentVolumeClaim)
//	//assert.NoError(t, err)
//	//assert.NotNil(t, obj)
//	//
//	//labels := map[string]string{"platform": "zbi", "project": "project", "instance": "instance"}
//	//utils.RemoveResourceField(obj, "spec.dataSource")
//	//utils.SetResourceField(obj, "metadata.labels", labels)
//	//
//	//resource, err := k.ApplyResource(ctx, obj)
//	//assert.NoError(t, err)
//	//assert.NotNil(t, resource)
//	//
//	//target, err := k.GetPersistentVolumes(ctx, labels)
//	//assert.NoError(t, err)
//	//assert.NotNil(t, target)
//	//assert.Len(t, target, 1)
//	//
//	//k.DeleteDynamicResource(ctx, resource.Namespace, resource.Name, helper.GvrMap[resource.Type])
//}
//
//func Test_GetPersistentVolumeClaimByName(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err)
//	assert.NotNil(t, k)
//
//	obj, err := data.GetGenericResource(model.ResourcePersistentVolumeClaim)
//	assert.NoError(t, err)
//	assert.NotNil(t, obj)
//
//	labels := map[string]string{"platform": "zbi", "project": "project", "instance": "instance"}
//	manager.RemoveResourceField(obj, "spec.dataSource")
//	manager.SetResourceField(obj, "metadata.labels", labels)
//
//	var robj = model.ResourceObject{Unstructured: *obj, Properties: nil}
//	resource, err := k.ApplyResource(ctx, &robj)
//	assert.NoError(t, err)
//	assert.NotNil(t, resource)
//
//	target, err := k.GetPersistentVolumeClaimByName(ctx, obj.GetNamespace(), obj.GetName())
//	assert.NoError(t, err)
//	assert.NotNil(t, target)
//	assert.Equal(t, obj.GetName(), target.Name)
//
//	k.DeleteDynamicResource(ctx, resource.Namespace, resource.Name, helper.GvrMap[resource.Type])
//}
//
//func Test_GetPersistentVolumeClaims(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err)
//	assert.NotNil(t, k)
//
//	obj, err := data.GetGenericResource(model.ResourcePersistentVolumeClaim)
//	assert.NoError(t, err)
//	assert.NotNil(t, obj)
//
//	labels := map[string]string{"platform": "zbi", "project": "project", "instance": "instance"}
//	manager.RemoveResourceField(obj, "spec.dataSource")
//	manager.SetResourceField(obj, "metadata.labels", labels)
//	var robj = model.ResourceObject{Unstructured: *obj, Properties: nil}
//
//	resource, err := k.ApplyResource(ctx, &robj)
//	assert.NoError(t, err)
//	assert.NotNil(t, resource)
//
//	target, err := k.GetPersistentVolumeClaims(ctx, obj.GetNamespace(), labels)
//	assert.NoError(t, err)
//	assert.NotNil(t, target)
//	assert.Len(t, target, 1)
//
//	k.DeleteDynamicResource(ctx, resource.Namespace, resource.Name, helper.GvrMap[resource.Type])
//}
//
//func Test_GetResource(t *testing.T) {
//
//}
//
//func Test_GenerateKubernetesObjects(t *testing.T) {
//
//}
//
//func Test_GetVolumeSnapshot(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err)
//	assert.NotNil(t, k)
//
//	claim, err := data.GetGenericResource(model.ResourcePersistentVolumeClaim)
//	assert.NoError(t, err)
//	assert.NotNil(t, claim)
//	labels := map[string]string{"platform": "zbi", "project": "project", "instance": "instance", "volume": "test"}
//	manager.RemoveResourceField(claim, "spec.dataSource")
//	manager.SetResourceField(claim, "metadata.labels", labels)
//	var robj = model.ResourceObject{Unstructured: *claim, Properties: nil}
//
//	claimResource, err := k.ApplyResource(ctx, &robj)
//	assert.NoError(t, err)
//	assert.NotNil(t, claimResource)
//
//	obj, err := data.GetGenericResource(model.ResourceVolumeSnapshot)
//	assert.NoError(t, err)
//	assert.NotNil(t, obj)
//	manager.SetResourceField(obj, "metadata.namespace", claim.GetNamespace())
//	manager.SetResourceField(obj, "spec.source.persistentVolumeClaimName", claim.GetName())
//	manager.SetResourceField(obj, "metadata.labels", labels)
//
//	robj = model.ResourceObject{Unstructured: *obj, Properties: nil}
//	resource, err := k.ApplyResource(ctx, &robj)
//	assert.NoError(t, err)
//	assert.NotNil(t, resource)
//
//	target, err := k.GetVolumeSnapshot(ctx, obj.GetNamespace(), obj.GetName())
//	assert.NoError(t, err)
//	assert.NotNil(t, target)
//
//	k.DeleteDynamicResource(ctx, resource.Namespace, resource.Name, helper.GvrMap[resource.Type])
//	k.DeleteDynamicResource(ctx, claimResource.Namespace, claimResource.Name, helper.GvrMap[claimResource.Type])
//}
//
//func Test_GetVolumeSnapshots(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err)
//	assert.NotNil(t, k)
//
//	claim, err := data.GetGenericResource(model.ResourcePersistentVolumeClaim)
//	assert.NoError(t, err)
//	assert.NotNil(t, claim)
//	labels := map[string]string{"platform": "zbi", "project": "project", "instance": "instance", "volume": "test"}
//	manager.RemoveResourceField(claim, "spec.dataSource")
//	manager.SetResourceField(claim, "metadata.labels", labels)
//
//	var robj = model.ResourceObject{Unstructured: *claim, Properties: nil}
//	claimResource, err := k.ApplyResource(ctx, &robj)
//	assert.NoError(t, err)
//	assert.NotNil(t, claimResource)
//	//	t.Logf("Claim: %s", utils.MarshalIndentObject(claim))
//
//	obj, err := data.GetGenericResource(model.ResourceVolumeSnapshot)
//	assert.NoError(t, err)
//	assert.NotNil(t, obj)
//	manager.SetResourceField(obj, "metadata.namespace", claim.GetNamespace())
//	manager.SetResourceField(obj, "spec.source.persistentVolumeClaimName", claim.GetName())
//	manager.SetResourceField(obj, "metadata.labels", labels)
//
//	robj = model.ResourceObject{Unstructured: *obj, Properties: nil}
//	resource, err := k.ApplyResource(ctx, &robj)
//	assert.NoError(t, err)
//	assert.NotNil(t, resource)
//	//	t.Logf("Snapshot: %s", utils.MarshalIndentObject(obj))
//
//	target := k.GetVolumeSnapshots(ctx, obj.GetNamespace(), labels)
//	assert.NoError(t, err)
//	assert.NotNil(t, target)
//	assert.Len(t, target, 1)
//
//	k.DeleteDynamicResource(ctx, resource.Namespace, resource.Name, helper.GvrMap[resource.Type])
//	k.DeleteDynamicResource(ctx, claimResource.Namespace, claimResource.Name, helper.GvrMap[claimResource.Type])
//}
//
//func Test_GetSnapshotSchedule(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err)
//	assert.NotNil(t, k)
//
//	claim, err := data.GetGenericResource(model.ResourcePersistentVolumeClaim)
//	assert.NoError(t, err)
//	assert.NotNil(t, claim)
//	pvcLabels := map[string]string{"platform": "zbi", "project": "project", "instance": "instance", "volume": "test"}
//	manager.RemoveResourceField(claim, "spec.dataSource")
//	manager.SetResourceField(claim, "metadata.labels", pvcLabels)
//
//	obj, err := data.GetGenericResource(model.ResourceSnapshotSchedule)
//	assert.NoError(t, err)
//	assert.NotNil(t, obj)
//	manager.SetResourceField(obj, "metadata.labels", pvcLabels)
//	manager.SetResourceField(obj, "spec.claimSelector.matchLabels", pvcLabels)
//	manager.SetResourceField(obj, "spec.snapshotTemplate.labels", pvcLabels)
//	//	t.Logf("Schedule Resource: %s", utils.MarshalIndentObject(obj))
//
//	var robj = model.ResourceObject{Unstructured: *claim, Properties: nil}
//	resource, err := k.ApplyResource(ctx, &robj)
//	assert.NoError(t, err)
//	assert.NotNil(t, resource)
//
//	target, err := k.GetSnapshotSchedule(ctx, obj.GetNamespace(), obj.GetName())
//	assert.NoError(t, err)
//	assert.NotNil(t, target)
//
//	k.DeleteDynamicResource(ctx, resource.Namespace, resource.Name, helper.GvrMap[resource.Type])
//	//	t.Logf("Schedule Result: %s", utils.MarshalIndentObject(target))
//}
//
//func Test_GetSnapshotSchedules(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err)
//	assert.NotNil(t, k)
//
//	claim, err := data.GetGenericResource(model.ResourcePersistentVolumeClaim)
//	assert.NoError(t, err)
//	assert.NotNil(t, claim)
//	pvcLabels := map[string]string{"platform": "zbi", "project": "project", "instance": "instance", "volume": "test"}
//	manager.RemoveResourceField(claim, "spec.dataSource")
//	manager.SetResourceField(claim, "metadata.labels", pvcLabels)
//
//	obj, err := data.GetGenericResource(model.ResourceSnapshotSchedule)
//	assert.NoError(t, err)
//	assert.NotNil(t, obj)
//	manager.SetResourceField(obj, "metadata.labels", pvcLabels)
//	manager.SetResourceField(obj, "spec.claimSelector.matchLabels", pvcLabels)
//	manager.SetResourceField(obj, "spec.snapshotTemplate.labels", pvcLabels)
//
//	var robj = model.ResourceObject{Unstructured: *obj, Properties: nil}
//	resource, err := k.ApplyResource(ctx, &robj)
//	assert.NoError(t, err)
//	assert.NotNil(t, resource)
//
//	target := k.GetSnapshotSchedules(ctx, obj.GetNamespace(), pvcLabels)
//	assert.NotNil(t, target)
//	assert.Len(t, target, 1)
//
//	k.DeleteDynamicResource(ctx, resource.Namespace, resource.Name, helper.GvrMap[resource.Type])
//}
//
//func Test_GetIngress(t *testing.T) {
//	ctx := context.Background()
//	k, err := NewKlient(ctx)
//	assert.NoError(t, err)
//	assert.NotNil(t, k)
//
//	zbiIngress, err := data.GetGenericResource(model.ResourceHTTPProxy)
//	assert.NoError(t, err)
//	assert.NotNil(t, zbiIngress)
//	zbiLabels := map[string]string{"platform": "zbi"}
//	manager.RemoveResourceField(zbiIngress, "spec.routes")
//	manager.SetResourceField(zbiIngress, "metadata.name", "zbi-proxy")
//	manager.SetResourceField(zbiIngress, "metadata.labels", zbiLabels)
//	//	t.Logf("ZBI App: %s", utils.MarshalIndentObject(zbiIngress))
//
//	var robj = model.ResourceObject{Unstructured: *zbiIngress, Properties: nil}
//	zbiResource, err := k.ApplyResource(ctx, &robj)
//	assert.NoError(t, err)
//	assert.NotNil(t, zbiResource)
//
//	zbiTarget, err := k.GetIngress(ctx, zbiIngress.GetNamespace(), zbiIngress.GetName())
//	assert.NoError(t, err)
//	assert.NotNil(t, zbiTarget)
//
//	//	t.Logf("ZBI App: %s", utils.MarshalIndentObject(zbiTarget))
//
//	projIngress, err := data.GetGenericResource(model.ResourceHTTPProxy)
//	assert.NoError(t, err)
//	assert.NotNil(t, zbiIngress)
//	projLabels := map[string]string{"platform": "zbi", "project": "project"}
//	manager.RemoveResourceField(projIngress, "spec.includes")
//	manager.RemoveResourceField(projIngress, "spec.virtualhost")
//	manager.SetResourceField(projIngress, "metadata.name", "project-proxy")
//	manager.SetResourceField(projIngress, "metadata.labels", projLabels)
//	//	t.Logf("Project App: %s", utils.MarshalIndentObject(projIngress))
//
//	robj = model.ResourceObject{Unstructured: *projIngress, Properties: nil}
//	projResource, err := k.ApplyResource(ctx, &robj)
//	assert.NoError(t, err)
//	assert.NotNil(t, projResource)
//
//	projTarget, err := k.GetIngress(ctx, projIngress.GetNamespace(), projIngress.GetName())
//	assert.NoError(t, err)
//	assert.NotNil(t, projTarget)
//
//	k.DeleteDynamicResource(ctx, zbiResource.Namespace, zbiResource.Name, helper.GvrMap[zbiResource.Type])
//	k.DeleteDynamicResource(ctx, projResource.Namespace, projResource.Name, helper.GvrMap[projResource.Type])
//
//	//	t.Logf("ZBI App: %s", utils.MarshalIndentObject(projTarget))
//
//}
