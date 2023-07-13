package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/zbitech/controller/internal/utils"
	"github.com/zbitech/controller/pkg/interfaces"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"strings"

	//	appsv1 "k8s.io/api/apps/v1"
	//	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	"time"
	// dynamicfake "k8s.io/client-go/dynamic/fake"
	// kubernetesfake "k8s.io/client-go/kubernetes/fake"
)

var (
	decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	GvrMap          = map[model.ResourceObjectType]schema.GroupVersionResource{
		model.ResourceNamespace:             {Group: "", Version: "v1", Resource: "namespaces"},
		model.ResourceConfigMap:             {Group: "", Version: "v1", Resource: "configmaps"},
		model.ResourceSecret:                {Group: "", Version: "v1", Resource: "secrets"},
		model.ResourceDeployment:            {Group: "apps", Version: "v1", Resource: "deployments"},
		model.ResourcePod:                   {Group: "", Version: "v1", Resource: "pods"},
		model.ResourceService:               {Group: "", Version: "v1", Resource: "services"},
		model.ResourcePersistentVolume:      {Group: "", Version: "v1", Resource: "persistentvolumes"},
		model.ResourcePersistentVolumeClaim: {Group: "", Version: "v1", Resource: "persistentvolumeclaims"},
		model.ResourceVolumeSnapshot:        {Group: "snapshot.storage.k8s.io", Version: "v1", Resource: "volumesnapshots"},
		model.ResourceVolumeSnapshotClass:   {Group: "snapshot.storage.k8s.io", Version: "v1", Resource: "volumesnapshotclasses"},
		model.ResourceSnapshotSchedule:      {Group: "snapscheduler.backube", Version: "v1", Resource: "snapshotschedules"},
		model.ResourceHTTPProxy:             {Group: "projectcontour.io", Version: "v1", Resource: "httpproxies"},
	}

	JSONSerializer = k8sjson.NewSerializerWithOptions(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, k8sjson.SerializerOptions{Pretty: true})
	YAMLSerializer = k8sjson.NewSerializerWithOptions(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, k8sjson.SerializerOptions{Yaml: true})

	TIME_LAYOUT = "2006-01-02T15:04:05Z"
)

//func CalculateAge(ts time.Time) string {
//	return time.Now().Sub(ts).String()
//}

func FilterLabels(objLabels map[string]string, filterLabels map[string]string) bool {
	if filterLabels == nil {
		return true
	}

	found := 0
	for key, value := range filterLabels {
		if objValue, ok := objLabels[key]; ok && objValue == value {
			found++
		}
	}

	return len(filterLabels) == found
}

func DecodeFromYaml(resourceYaml string) (*unstructured.Unstructured, *schema.GroupVersionKind, error) {
	object := &unstructured.Unstructured{}

	s := k8sjson.NewSerializerWithOptions(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, k8sjson.SerializerOptions{false, false, false})
	s.Decode([]byte(resourceYaml), nil, object)

	_, gvk, err := decUnstructured.Decode([]byte(resourceYaml), nil, object)
	return object, gvk, err
}

func CreateUnstructuredResource(result *unstructured.Unstructured) *model.KubernetesResource {
	var status = GetResourceStatusField(result)
	var properties = GetResourceProperties(result)
	created := time.Now()
	objType := model.ResourceObjectType(result.GetKind())

	return &model.KubernetesResource{
		Name: result.GetName(),
		//		Project:    object.Project,
		//		Instance:   object.Instance,
		Namespace: result.GetNamespace(),
		Type:      objType,
		Status:    status,
		Created:   &created,
		//		Updated:    &created,
		Properties: properties,
	}
}

func GetResourceRequest(ctx context.Context, instance *model.Instance) (*model.ResourceRequest, bool) {
	for _, resource := range instance.Resources.Resources {
		if resource.Type == model.ResourceConfigMap {
			properties, ok := resource.Properties["request"]
			if ok {
				var request model.ResourceRequest
				if err := json.Unmarshal([]byte(properties.(string)), &request); err != nil {
					logger.GetLogger(ctx).WithFields(logrus.Fields{"error": err}).Errorf("failed to marshal resource request")
					return nil, false
				}

				return &request, true
			}
		}
	}

	return nil, false
}

func createResourceRequest(ctx context.Context, resource *model.KubernetesResource) (*model.ResourceRequest, bool) {
	properties, ok := resource.Properties["request"]
	if ok {
		var request model.ResourceRequest
		if err := json.Unmarshal([]byte(properties.(string)), &request); err != nil {
			logger.GetLogger(ctx).WithFields(logrus.Fields{"error": err}).Errorf("failed to marshal resource request")
			return nil, false
		}

		return &request, true
	}

	return nil, false
}

func CreateCoreResource(ctx context.Context, rType model.ResourceObjectType, result runtime.Object, client interfaces.KlientIF) *model.KubernetesResource {

	var name, status, namespace string
	//	var labels map[string]string
	var properties = make(map[string]interface{}, 0)
	var created time.Time
	//var updated time.Time
	//var err error

	if rType == model.ResourceConfigMap {

		cmap := result.(*corev1.ConfigMap)
		name = cmap.Name
		namespace = cmap.Namespace
		status = "active"
		created = cmap.ObjectMeta.CreationTimestamp.Time
		for key, value := range cmap.Data {
			properties[key] = value
		}

	} else if rType == model.ResourceSecret {

		secret := result.(*corev1.Secret)
		name = secret.Name
		namespace = secret.Namespace
		status = "active"
		created = secret.ObjectMeta.CreationTimestamp.Time
		for key, value := range secret.Data {
			properties[key] = string(value)
		}

	} else if rType == model.ResourceDeployment {

		dep := result.(*appsv1.Deployment)
		name = dep.Name
		namespace = dep.Namespace
		created = dep.ObjectMeta.CreationTimestamp.Time
		status = GetDeploymentStatus(dep)
		properties = GetDeploymentProperties(dep)

		pods := client.GetPods(ctx, dep.Namespace, dep.Spec.Template.Labels)
		podMap := make(map[string]interface{})
		for _, p := range pods {
			startTime := time.Now()
			pod_status := strings.ToLower(string(p.Status.Phase))

			if pod_status != "running" {
				status = "progressing"
			}

			if p.Status.StartTime != nil {
				startTime = p.Status.StartTime.Time
			}
			podMap[p.Name] = map[string]interface{}{
				"status":    pod_status,
				"startTime": startTime,
			}
		}
		properties["pods"] = podMap

	} else if rType == model.ResourcePod {
		pod := result.(*corev1.Pod)
		name = pod.Name
		namespace = pod.Namespace
		created = pod.ObjectMeta.CreationTimestamp.Time
		status = GetPodStatus(pod)
		properties = GetPodProperties(pod)

	} else if rType == model.ResourceService {

		svc := result.(*corev1.Service)
		name = svc.Name
		namespace = svc.Namespace
		status = "active"
		created = svc.ObjectMeta.CreationTimestamp.Time

	} else if rType == model.ResourcePersistentVolumeClaim {

		pvc := result.(*corev1.PersistentVolumeClaim)
		name = pvc.Name
		namespace = pvc.Namespace
		created = pvc.ObjectMeta.CreationTimestamp.Time
		status = GetPersistentVolumeClaimStatus(pvc)
		properties = GetPersistenVolumeClaimProperties(pvc)

	} else if rType == model.ResourceHTTPProxy {

		ing := result.(*unstructured.Unstructured)
		name = ing.GetName()
		created = GetResourceCreationTime(ing)
		namespace = ing.GetNamespace()
		status = GetResourceStatusField(ing)

	} else if rType == model.ResourceVolumeSnapshot {
		vs := result.(*unstructured.Unstructured)
		name = vs.GetName()

		created = GetResourceCreationTime(vs)
		namespace = vs.GetNamespace()
		status = GetResourceStatusField(vs)
		properties = GetResourceProperties(vs)

	} else if rType == model.ResourceSnapshotSchedule {

	}

	return &model.KubernetesResource{
		Name:       name,
		Namespace:  namespace,
		Type:       rType,
		Status:     status,
		Created:    &created,
		Updated:    nil,
		Properties: properties,
	}
}

//func GetConfigMapProperties(cm *corev1.ConfigMap) map[string]interface{} {
//
//}
//
//func GetSecretProperties(secret *corev1.Secret) map[string]interface{} {
//
//}

func GetPodStatus(pod *corev1.Pod) string {
	return strings.ToLower(string(pod.Status.Phase))
}

func GetPodProperties(pod *corev1.Pod) map[string]interface{} {
	startTime := time.Now()
	if pod.Status.StartTime != nil {
		startTime = pod.Status.StartTime.Time
	}
	replicasetName := pod.ObjectMeta.OwnerReferences[0].Name
	//	podTemplateHash := pod.ObjectMeta.Labels["pod-template-hash"]
	lastIndex := strings.LastIndex(replicasetName, "-")

	return map[string]interface{}{
		"deployment": replicasetName[:lastIndex],
		"startTime":  startTime,
	}
}

func GetDeploymentStatus(dep *appsv1.Deployment) string {

	var status string

	replicas := dep.Status.Replicas
	available := dep.Status.AvailableReplicas
	ready := dep.Status.ReadyReplicas

	if replicas == available && replicas == ready {
		status = "active"
	} else if available > 0 && ready > 0 {
		status = "progressing"
	} else {
		status = "pending"
	}

	return status
}

func GetDeploymentProperties(dep *appsv1.Deployment) map[string]interface{} {

	properties := make(map[string]interface{})

	properties["image"] = dep.Spec.Template.Spec.Containers[0].Image
	res := dep.Spec.Template.Spec.Containers[0].Resources
	limits := make(map[string]string)
	requests := make(map[string]string)

	if res.Limits.Cpu() != nil {
		limits["cpu"] = res.Limits.Cpu().String()
	}

	if res.Limits.Memory() != nil {
		limits["memory"] = res.Limits.Memory().String()
	}

	if res.Requests.Cpu() != nil {
		requests["cpu"] = res.Requests.Cpu().String()
	}

	if res.Requests.Memory() != nil {
		requests["memory"] = res.Requests.Memory().String()
	}

	properties["resources"] = map[string]interface{}{"limits": limits, "requests": requests}
	return properties
}

func GetPersistentVolumeClaimStatus(pvc *corev1.PersistentVolumeClaim) string {
	return strings.ToLower(string(pvc.Status.Phase))
}

func GetPersistenVolumeClaimProperties(pvc *corev1.PersistentVolumeClaim) map[string]interface{} {

	properties := make(map[string]interface{})

	properties["requestedStorage"] = pvc.Spec.Resources.Requests.Storage().String()
	properties["actualStorage"] = pvc.Status.Capacity.Storage().String()
	properties["storageClassName"] = *pvc.Spec.StorageClassName
	properties["volumeName"] = pvc.Spec.VolumeName

	return properties
}

func AddKubernetesResources(ctx context.Context, client interfaces.KlientIF, resources []model.KubernetesResource, rType model.ResourceObjectType, objects ...interface{}) []model.KubernetesResource {

	if objects != nil {
		for _, object := range objects {
			resources = append(resources, *CreateCoreResource(ctx, rType, object.(runtime.Object), client))
		}
	}

	return resources
}

func GetDynamicResourceInterface(dynamicClient dynamic.Interface, object *unstructured.Unstructured) dynamic.ResourceInterface {
	var dr dynamic.ResourceInterface
	kind := object.GroupVersionKind().Kind
	kindType := model.ResourceObjectType(kind)
	gvr := GvrMap[kindType]
	if kindType == model.ResourceNamespace {
		dr = dynamicClient.Resource(gvr)
	} else {
		dr = dynamicClient.Resource(gvr).Namespace(object.GetNamespace())
	}

	return dr
}

func DecodeYAML(yaml string, object *unstructured.Unstructured) error {
	_, _, err := YAMLSerializer.Decode([]byte(yaml), nil, object)
	if err != nil {
		return err
	}

	return nil
}

func EncodeYAML(object *unstructured.Unstructured) (string, error) {
	var buffer = new(bytes.Buffer)
	if err := YAMLSerializer.Encode(object, buffer); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func DecodeJSON(data string, object *unstructured.Unstructured) error {
	_, _, err := JSONSerializer.Decode([]byte(data), nil, object)
	if err != nil {
		return err
	}

	return nil
}

func EncodeJSON(object *unstructured.Unstructured) (string, error) {
	var buffer = new(bytes.Buffer)
	if err := JSONSerializer.Encode(object, buffer); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func CreateYAMLObjects(specArr []string) ([]unstructured.Unstructured, error) {
	var objects = make([]unstructured.Unstructured, len(specArr))
	for index, yamlString := range specArr {
		var object unstructured.Unstructured
		if err := DecodeYAML(yamlString, &object); err != nil {
			return nil, err
		}
		objects[index] = object
	}

	return objects, nil
}

func CreateYAMLObject(yamlString string) (*unstructured.Unstructured, error) {
	var object unstructured.Unstructured
	if err := DecodeYAML(yamlString, &object); err != nil {
		return nil, err
	}

	return &object, nil
}

func CreateProjectLabels(project *model.Project) map[string]string {
	return map[string]string{
		"platform": "zbi",
		"project":  project.Name,
		"owner":    project.Owner,
		"team":     project.TeamId,
		"network":  string(project.Network),
		"level":    "project",
	}
}

func CreateInstanceLabels(instance *model.Instance) map[string]string {
	return map[string]string{
		"platform": "zbi",
		"project":  instance.Project,
		"instance": instance.Name,
		"type":     string(instance.InstanceType),
		"level":    "instance",
	}
}

func CreateEnvoySpec(envoyServicePort int32) model.EnvoySpec {
	envoy := Config.GetPolicyConfig().Envoy

	return model.EnvoySpec{
		Image:                 envoy.Image,
		Command:               utils.MarshalObject(envoy.Command),
		Port:                  envoyServicePort,
		Timeout:               envoy.Timeout,
		AccessAuthorization:   envoy.AccessAuthorization,
		AuthServerURL:         envoy.AuthServerURL,
		AuthServerPort:        envoy.AuthServerPort,
		AuthenticationEnabled: envoy.AuthenticationEnabled,
	}
}

func CreateSnapshotSchedule(schedule model.SnapshotScheduleType) string {
	if schedule == model.DailySnapshotSchedule {
		hour := 5
		min := 1
		return fmt.Sprintf("%d %d * * *", min, hour)
	} else if schedule == model.WeeklySnapshotSchedule {
		weekDay := 1
		return fmt.Sprintf("* * * * %d", weekDay)
	} else if schedule == model.MonthlySnapshotSchedule {
		day := 1
		month := 1
		return fmt.Sprintf("* * %d %d *", day, month)
	}

	return ""
}

func CreateIngressRoute(ctx context.Context, specObj string) (*model.IngressRoute, error) {

	var route model.IngressRoute
	if err := json.Unmarshal([]byte(specObj), &route); err != nil {
		//		logger.Errorf(ctx, "Zcash route marshal failed - %s", err)
		return nil, err
	}

	return &route, nil
}

func UpdateIngressRoute(ctx context.Context, projIngress *unstructured.Unstructured, route *model.IngressRoute, remove bool) error {

	var err error

	if err = RemoveResourceField(projIngress, "metadata.managedFields"); err != nil {
		//		logger.Errorf(ctx, "Error removing metadata.managedFields - %s", err)
		//		return errs.NewApplicationError(errs.ResourceGenerationError, err)
		return err
	}

	if err = RemoveResourceField(projIngress, "spec.status"); err != nil {
		//		logger.Errorf(ctx, "Error removing spec.status - %s", err)
		//		return errs.NewApplicationError(errs.ResourceGenerationError, err)
		return err
	}

	routeData := utils.MarshalObject(GetResourceField(projIngress, "spec.routes"))
	var routes []model.IngressRoute
	if err = json.Unmarshal([]byte(routeData), &routes); err != nil {
		//		logger.Errorf(ctx, "Error unmarshalling ingress routes - %s", err)
		//		return errs.NewApplicationError(errs.ResourceGenerationError, err)
		return err
	}

	var updated = false
	for index, r := range routes {
		for _, condition := range r.Conditions {
			//			logger.Infof(ctx, "Comparing %s and %s at index %d ...", condition.Prefix, route.Conditions[0].Prefix, index)
			if condition.Prefix == route.Conditions[0].Prefix {
				if remove {
					routes = append(routes[:index], routes[index+1:]...)
				} else {
					routes = append(routes[:index], *route)
					routes = append(routes, routes[index+1:]...)
				}
				updated = true
				break
			}
		}
	}

	if !updated {
		routes = append(routes, *route)
	}

	//	logger.Debugf(ctx, "Ingress routes: %s", utils.MarshalObject(routes))
	if err = SetResourceField(projIngress, "spec.routes", routes); err != nil {
		//		logger.Errorf(ctx, "Error setting spec.status - %s", err)
		//		return errs.NewApplicationError(errs.ResourceGenerationError, err)
		return err
	}

	return nil
}

func GetResourceLabels(object *unstructured.Unstructured) map[string]string {
	labels := make(map[string]string)

	for key, value := range GetResourceField(object, "metadata.labels").(map[string]interface{}) {
		labels[key] = value.(string)
	}

	return labels
}

func GetResourceField(obj *unstructured.Unstructured, path string) interface{} {
	var content = obj.UnstructuredContent()
	parts := strings.Split(path, ".")
	for index, part := range parts {
		if index == len(parts)-1 {
			return content[part]
		} else {
			entry := content[part]
			if entry == nil {
				return nil
			}
			content = entry.(map[string]interface{})
		}
	}
	return nil
}

func GetResourceCreationTime(obj *unstructured.Unstructured) time.Time {
	created, err := time.Parse(TIME_LAYOUT, GetResourceField(obj, "metadata.creationTimestamp").(string))
	if err != nil {
		return time.Time{}
	}

	return created
}

func ReadResourceField(obj *unstructured.Unstructured, path string, data interface{}) error {
	value := GetResourceField(obj, path)
	if value != nil {

		valueBytes, err := json.Marshal(value)
		if err != nil {
			return err
		}

		return json.Unmarshal(valueBytes, &data)
	}
	return nil
}

func SetResourceField(obj *unstructured.Unstructured, path string, value interface{}) error {
	var content = obj.UnstructuredContent()
	parts := strings.Split(path, ".")
	for index, part := range parts {
		if index == len(parts)-1 {
			content[part] = value
		} else {
			entry := content[part]
			if entry == nil {
				entry = make(map[string]interface{})
				content[part] = entry
			}
			content = entry.(map[string]interface{})
		}
	}
	return nil
}

func AddResourceField(obj *unstructured.Unstructured, path string, value interface{}) error {
	var content = obj.UnstructuredContent()
	parts := strings.Split(path, ".")
	for index, part := range parts {
		if index == len(parts)-1 {
			var array []interface{}
			entry := content[part]
			if entry != nil {
				array = entry.([]interface{})
			}
			array = append(array, value)
			content[part] = array
		} else {
			entry := content[part]
			if entry == nil {
				entry = make(map[string]interface{})
				content[part] = entry
			}
			content = entry.(map[string]interface{})
		}
	}
	return nil
}

func RemoveResourceField(obj *unstructured.Unstructured, path string) error {
	var content = obj.UnstructuredContent()
	parts := strings.Split(path, ".")
	for index, part := range parts {
		if index == len(parts)-1 {
			if content[part] != nil {
				delete(content, part)
			}
		} else {
			entry := content[part]
			if entry == nil {
				entry = make(map[string]interface{})
				content[part] = entry
			}
			content = entry.(map[string]interface{})
		}
	}
	return nil
}

// GetResourceProperties returns the corresponding property type for a kubernetes resource
// returns map of data entries for ConfigMap
// returns map of data entries (base64 decoded) for Secret
// returns map of requested size, actual size, storage class name and volume name for PersistentVolumeClaim
// returns an empty map for all other resources
func GetResourceProperties(obj *unstructured.Unstructured) map[string]interface{} {
	kind := obj.GetKind()

	switch model.ResourceObjectType(kind) {
	case model.ResourceConfigMap:

		return GetResourceField(obj, "data").(map[string]interface{})

	case model.ResourceSecret:

		data := make(map[string]interface{}, 0)
		for key, value := range GetResourceField(obj, "data").(map[string]interface{}) {
			data[key] = utils.Base64DecodeString(value.(string))
		}
		return data

	case model.ResourcePersistentVolumeClaim:

		return map[string]interface{}{
			"requestedStorage": GetResourceField(obj, "spec.resources.requests.storage"),
			"actualStorage":    GetResourceField(obj, "status.capacity.storage"),
			"storageClassName": GetResourceField(obj, "spec.storageClassName"),
			"volumeName":       GetResourceField(obj, "spec.volumeName"),
		}

	case model.ResourceVolumeSnapshot:

		return map[string]interface{}{
			"contentName":   GetResourceField(obj, "status.boundVolumeSnapshotContentName"),
			"snapshotClass": GetResourceField(obj, "spec.volumeSnapshotClassName"),
			"volumeName":    GetResourceField(obj, "spec.source.persistentVolumeClaimName"),
			"creationTime":  GetResourceField(obj, "status.creationTime"),
			"restoreSize":   GetResourceField(obj, "status.restoreSize"),
		}

	case model.ResourceDeployment:
		_containers := GetResourceField(obj, "spec.template.spec.containers")
		if _containers != nil {
			containers := _containers.([]interface{})
			container := containers[0].(map[string]interface{})
			return map[string]interface{}{
				"image":     container["image"].(string),
				"resources": container["resources"],
			}
		}

	case model.ResourcePod:
		ownerReferences := GetResourceField(obj, "metadata.ownerReferences").([]interface{})
		ownerReference := ownerReferences[0].(map[string]interface{})

		startTime, _ := time.Parse(TIME_LAYOUT, GetResourceField(obj, "status.startTime").(string))
		return map[string]interface{}{
			"deployment": ownerReference["name"],
			"startTime":  startTime,
		}
	}

	return map[string]interface{}{}
}

// GetResourceStatusField returns the corresponding status for a kubernetes resource
// returns the status based on available replicas for a Deployment resource (active, partial or pending)
// returns the phase for a PersistentVolumeClaim resource
// returns active for all other resources
func GetResourceStatusField(obj *unstructured.Unstructured) string {

	kind := obj.GetKind()

	status := "active"

	switch model.ResourceObjectType(kind) {
	case model.ResourceDeployment:
		_availableReplicas := GetResourceField(obj, "status.availableReplicas")
		_readyReplicas := GetResourceField(obj, "status.readyReplicas")
		_replicas := GetResourceField(obj, "status.replicas")

		var availableReplicas = 0
		var readyReplicas = 0
		var replicas = 1

		if _availableReplicas != nil {
			availableReplicas = _availableReplicas.(int)
		}

		if _readyReplicas != nil {
			readyReplicas = _readyReplicas.(int)
		}

		if _replicas != nil {
			replicas = _replicas.(int)
		}

		if replicas == readyReplicas && availableReplicas == replicas {
			status = "active"
		} else if readyReplicas > 0 && availableReplicas > 0 {
			status = "partial"
		} else {
			status = "pending"
		}

	case model.ResourcePod:
		status = strings.ToLower(GetResourceField(obj, "status.phase").(string))

	case model.ResourcePersistentVolumeClaim:
		status = strings.ToLower(GetResourceField(obj, "status.phase").(string))

	case model.ResourceVolumeSnapshot:
		_status := GetResourceField(obj, "status.readyToUse")
		if _status != nil {
			if _status.(bool) {
				status = "ready"
			} else {
				status = "notReady"
			}
		} else {
			status = "notReady"
		}

	case model.ResourceSnapshotSchedule:
		status = ""

	case model.ResourceHTTPProxy:
		status = strings.ToLower(GetResourceField(obj, "status.currentStatus").(string))
	}

	return status
}

func ParseTime(t string) time.Time {
	pt, err := time.Parse(TIME_LAYOUT, t)
	if err != nil {
		return time.Time{}
	}

	return pt
}
