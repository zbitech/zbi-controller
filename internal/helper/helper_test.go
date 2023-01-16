package helper

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	ASSETPATH  = "/Users/johnakinyele/go/src/github.com/zbi/data/etc/zbi"
	KUBECONFIG = "/Users/johnakinyele/.kube/config"
	NginxGVK   = schema.GroupVersionKind{Group: "apps", Kind: "Deployment", Version: "v1"}
	Group      = "apps"
	Kind       = "Deployment"
	Version    = "v1"
	Name       = "nginx"
	NginxYAML  = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  namespace: default
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:latest
        ports:
        - containerPort: 80
`
	// initialized = test.InitTestConfig()
	initialized = true
)

func getClients(t *testing.T, k8s_client, dyn_client, create_mapper bool) (kubernetes.Interface, dynamic.Interface, *restmapper.DeferredDiscoveryRESTMapper) {
	restConfig, err := clientcmd.BuildConfigFromFlags("", KUBECONFIG)
	assert.NoErrorf(t, err, "Failed to get rest config")

	var kubernetesClient kubernetes.Interface
	var dynamicClient dynamic.Interface
	var mapper *restmapper.DeferredDiscoveryRESTMapper

	if k8s_client {
		kubernetesClient, err = kubernetes.NewForConfig(restConfig)
		assert.NoErrorf(t, err, "Unable to create kuberntes client - %s", err)
		assert.NotNilf(t, kubernetesClient, "Unable to create kubernetes clientset")
	}

	if create_mapper {
		mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(kubernetesClient.Discovery()))
		assert.NotNilf(t, mapper, "Unable to create mapper object")
	}

	if dyn_client {
		dynamicClient, err := dynamic.NewForConfig(restConfig)
		assert.NoErrorf(t, err, "Unable to create kuberntes client - %s", err)
		assert.NotNilf(t, dynamicClient, "Unable to create kubernetes dynamic client")
	}

	return kubernetesClient, dynamicClient, mapper
}

func createResource(t *testing.T, ctx context.Context, dynamicClient dynamic.Interface, mapper *restmapper.DeferredDiscoveryRESTMapper,
	objJSON string) *unstructured.Unstructured {
	var obj unstructured.Unstructured
	err := json.Unmarshal([]byte(objJSON), &obj)
	assert.NoError(t, err, "Failed to generate pod object - %s", err)

	data, err := json.Marshal(obj)
	assert.NoError(t, err, "Failed to generate pod JSON object - %s", err)

	//dr, err := GetDynamicResourceInterface(dynamicClient, mapper, &obj)
	dr := GetDynamicResourceInterface(dynamicClient, &obj)
	_, err = dr.Patch(ctx, obj.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{FieldManager: "zbi-controller"})
	assert.NoError(t, err, "Failed to create pod object - %s", err)

	return &obj
}

func Test_FilterLabels(t *testing.T) {

	tests := []struct {
		Name   string
		Labels map[string]string
		Filter map[string]string
		Want   bool
	}{
		{Name: "Unfiltered", Labels: map[string]string{"type": "zcash", "kind": "deployment"}, Filter: nil, Want: true},
		{Name: "Unfiltered", Labels: map[string]string{"type": "zcash", "kind": "deployment"}, Filter: map[string]string{"type": "zcash"}, Want: true},
		{Name: "Unfiltered", Labels: map[string]string{"type": "lwd", "kind": "deployment"}, Filter: map[string]string{"type": "zcash"}, Want: false},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			exists := FilterLabels(test.Labels, test.Filter)
			if test.Want != exists {
				t.Errorf("got %v, want %v", exists, test.Want)
			}
		})
	}
}

func Test_DecodeFromYaml(t *testing.T) {

	obj, gvk, err := DecodeFromYaml(NginxYAML)
	if err != nil {
		t.Errorf("Got an error %s", err)
	}

	if gvk.Group != NginxGVK.Group || gvk.Kind != NginxGVK.Kind || gvk.Version != NginxGVK.Version {
		t.Logf("got %s.%s/%s, want %s.%s/%s", gvk.Group, gvk.Kind, gvk.Version, NginxGVK.Group, NginxGVK.Version, NginxGVK.Kind)
	}

	if obj == nil {
		t.Errorf("Got nil object")
	}

	if obj.GetName() != Name {
		t.Logf("got %s, want %s", obj.GetName(), Name)
	}
}

func Test_EncodeToJson(t *testing.T) {

	//proxyYaml := "apiVersion: projectcontour.io/v1\nkind: HTTPProxy\nmetadata:\n  name: mongo-express\n  namespace: mongodb\nspec:\n  virtualhost:\n    fqdn: api.zbitech.local\n    tls:\n      secretName: cert-manager/zbi-tls\n    includes:\n    - name: project1\n      namespace: project1\n      conditions:\n        - prefix: /project1\n    - name: project2\n      namespace: project2\n      conditions:\n        - prefix: /project2"
	//obj, _, _ := DecodeFromYaml(proxyYaml)
	//
	//var proxyObj = new(spec.ZBIIngress)
	//err := EncodeToJson(obj, proxyObj)
	//assert.NoError(t, err, "Failed to convert to YAML")
	//fmt.Printf("%s", yamlStr)

	//var proxyObj model.ZBIIngress
	//err = json.Unmarshal([]byte(yamlStr), &proxyObj)
	//assert.NoError(t, err, "Failed to convert object")
	//t.Logf("%v", proxyObj.Spec.VirtualHost.Includes)
}

//func Test_GetPodState(t *testing.T) {
//
//	ctx := context.Background()
//
//	kubernetesClient, _, _ := getClients(t, true, false, false)
//	//obj := createResource(t, ctx, dynamicClient, mapper, test.PodJSON)
//
//	pods, err := kubernetesClient.CoreV1().Pods("kube-system").List(ctx, metav1.ListOptions{})
//	assert.NoError(t, err, "Failed to get pod objects - %s", err)
//	assert.NotNilf(t, pods, "Failed to get pod objects")
//	assert.Greater(t, len(pods.Items), 0, "Failed to get at least 1 active pod")
//	pod := pods.Items[0]
//
//	//pod, err := kubernetesClient.CoreV1().Pods(obj.GetNamespace()).Get(ctx, obj.GetName(), metav1.GetOptions{})
//	//assert.NoError(t, err, "Failed to get pod object - %s", err)
//
//	podState := GetPodState(pod.Status)
//	t.Logf("POD State: %s", utils.MarshalObject(podState))
//}

//func Test_GetDeploymentStatus(t *testing.T) {
//
//	ctx := context.Background()
//
//	kubernetesClient, _, _ := getClients(t, true, false, false)
//	deployments, err := kubernetesClient.AppsV1().Deployments("kube-system").List(ctx, metav1.ListOptions{})
//	assert.NoError(t, err, "Failed to get pod objects - %s", err)
//	assert.NotNilf(t, deployments, "Failed to get pod objects")
//	assert.Greater(t, len(deployments.Items), 0, "Failed to get at least 1 active pod")
//	deployment := deployments.Items[0]
//
//	deploymentState := GetDeploymentStatus(deployment.Status)
//	t.Logf("Deployment State: %s", utils.MarshalObject(deploymentState))
//}

func Test_GetResourceStatus(t *testing.T) {

	//ctx := context.Background()
	//
	//restConfig, err := clientcmd.BuildConfigFromFlags("", vars.AppConfig.Kubernetes.KubeConfig)
	//if err != nil {
	//	t.Fatalf("Failed to create kubernetes client - %s", err)
	//}
	//
	//dynamicClient, err := dynamic.NewForConfig(restConfig)
	//if err != nil {
	//	t.Fatalf("Failed to create kubernetes client - %s", err)
	//}
	//
	//data, err := json.Marshal(object)
	//if err != nil {
	//	t.Fatalf("Failed to marshal object - %s", err)
	//}
	//
	//dr, err := GetDynamicResourceInterface(dynamicClient, k.Mapper, object)
	//if err != nil {
	//	t.Fatalf("Failed to get dynamic resource interface - %s", err)
	//}
	//
	//result, err := dr.Patch(ctx, object.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{FieldManager: "zbi-controller"})
	//state := GetResourceStatus(result)
	//t.Logf("%s Object: %s", result.GetKind(), utils.MarshalIndentObject(state))
}

func Test_GetResource(t *testing.T) {

	// object, _, err := DecodeFromYaml(NginxYAML)
	// if err != nil {
	// 	t.Errorf("Got an error instead of kubernets object %s", err)
	// }

	// vars.KubernetesKlient, err = NewKlientFactory(context.Background())
	// if err != nil {
	// 	t.Errorf("Got an error instead of kubernetes clint - %s", err)
	// }

	// result, err := vars.KubernetesKlient.ApplyResource(context.Background(), object)
	// if err != nil {
	// 	t.Errorf("Got an error instead of kubernetes deployment result - %s", err)
	// }
	// //t.Logf("Result - %s", fn.MarshalIndentObject(result))

	// gvr, err := GroupVersionResource(vars.KubernetesKlient.GetMapper(), object)
	// if err != nil {
	// 	t.Errorf("Got an error instead of GVR - %s", err)
	// }

	// obj, err := vars.KubernetesKlient.GetDynamicResource(context.Background(), result.GetNamespace(), result.GetName(), *gvr)
	// if err != nil {
	// 	t.Errorf("Got an error instead of kubernetes object - %s", err)
	// }

	// time.Sleep(time.Duration(15) * time.Second)
	// status, _, err := unstructured.NestedStringMap(obj.Object, "status")
	// if err != nil {
	// 	t.Errorf("Got an error while Getting deployment status - %s", err)
	// }
	// t.Logf("Deployment Object - %s", status)

	// dep, err := vars.KubernetesKlient.GetDeploymentByName(context.Background(), "default", Name)
	// if err != nil {
	// 	t.Errorf("Got an error while Getting deployment - %s", err)
	// }
	// t.Logf("Deployment - %s", fn.MarshalIndentObject(dep))

	// err = vars.KubernetesKlient.DeleteDynamicResource(context.Background(), result.GetNamespace(), result.GetName(), *gvr)
	// if err != nil {
	// 	t.Errorf("Got an error while deleting resource - %s", err)
	// }
}

func Test_GroupVersionResource(t *testing.T) {

}

func Test_GetDynamicResourceInterface(t *testing.T) {

}

func Test_GenerateKubernetesObjects(t *testing.T) {

}
