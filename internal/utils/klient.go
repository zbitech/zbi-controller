package utils

//var (
//	JSONSerializer = k8sjson.NewSerializerWithOptions(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, k8sjson.SerializerOptions{Pretty: true})
//	YAMLSerializer = k8sjson.NewSerializerWithOptions(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, k8sjson.SerializerOptions{Yaml: true})
//)

//func DecodeYAML(yaml string, object *unstructured.Unstructured) error {
//	_, _, err := YAMLSerializer.Decode([]byte(yaml), nil, object)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func EncodeYAML(object *unstructured.Unstructured) (string, error) {
//	var buffer = new(bytes.Buffer)
//	if err := YAMLSerializer.Encode(object, buffer); err != nil {
//		return "", err
//	}
//
//	return buffer.String(), nil
//}
//
//func DecodeJSON(data string, object *unstructured.Unstructured) error {
//	_, _, err := JSONSerializer.Decode([]byte(data), nil, object)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func EncodeJSON(object *unstructured.Unstructured) (string, error) {
//	var buffer = new(bytes.Buffer)
//	if err := JSONSerializer.Encode(object, buffer); err != nil {
//		return "", err
//	}
//
//	return buffer.String(), nil
//}
//
//func CreateYAMLObjects(specArr []string /*, project, instance string, propertiesArr ...map[string]interface{}*/) ([]model.ResourceObject, error) {
//	var objects = make([]model.ResourceObject, len(specArr))
//	for index, yamlString := range specArr {
//		var object unstructured.Unstructured
//		/*var properties map[string]interface{}
//		if propertiesArr != nil && len(propertiesArr) > index {
//			properties = propertiesArr[index]
//		}*/
//		if err := DecodeYAML(yamlString, &object); err != nil {
//			return nil, err
//		}
//		objects[index] = model.ResourceObject{Unstructured: object /*, Project: project, Instance: instance, Properties: properties*/}
//	}
//
//	return objects, nil
//}
//
//func CreateYAMLObject(yamlString, project, instance string, properties map[string]interface{}) (*model.ResourceObject, error) {
//	var object unstructured.Unstructured
//	if err := DecodeYAML(yamlString, &object); err != nil {
//		return nil, err
//	}
//
//	return &model.ResourceObject{Unstructured: object /*, Project: project, Instance: instance, Properties: properties*/}, nil
//}
//
//func CreateProjectLabels(project *model.Project) map[string]string {
//	return map[string]string{
//		"platform": "zbi",
//		"project":  project.Name,
//		"owner":    project.Owner,
//		"team":     project.TeamId,
//		"network":  string(project.Network),
//	}
//}
//
//func CreateInstanceLabels(instance *model.Instance) map[string]string {
//	return map[string]string{
//		"platform": "zbi",
//		"project":  instance.Project,
//		"instance": instance.Name,
//		"type":     string(instance.InstanceType),
//	}
//}
//
//func CreateEnvoySpec(envoyServicePort int32) model.EnvoySpec {
//	envoy := helper.Config.GetPolicyConfig().Envoy
//
//	return model.EnvoySpec{
//		Image:                 envoy.Image,
//		Command:               MarshalObject(envoy.Command),
//		Port:                  envoyServicePort,
//		Timeout:               envoy.Timeout,
//		AccessAuthorization:   envoy.AccessAuthorization,
//		AuthServerURL:         envoy.AuthServerURL,
//		AuthServerPort:        envoy.AuthServerPort,
//		AuthenticationEnabled: envoy.AuthenticationEnabled,
//	}
//}
//
//func CreateSnapshotSchedule(schedule model.SnapshotScheduleType) string {
//	if schedule == model.DailySnapshotSchedule {
//		hour := 5
//		min := 1
//		return fmt.Sprintf("%d %d * * *", min, hour)
//	} else if schedule == model.WeeklySnapshotSchedule {
//		weekDay := 1
//		return fmt.Sprintf("* * * * %d", weekDay)
//	} else if schedule == model.MonthlySnapshotSchedule {
//		day := 1
//		month := 1
//		return fmt.Sprintf("* * %d %d *", day, month)
//	}
//
//	return ""
//}
//
//func CreateIngressRoute(ctx context.Context, specObj string) (*model.IngressRoute, error) {
//
//	var route model.IngressRoute
//	if err := json.Unmarshal([]byte(specObj), &route); err != nil {
//		//		logger.Errorf(ctx, "Zcash route marshal failed - %s", err)
//		return nil, err
//	}
//
//	//	logger.Debugf(ctx, "Route: %s", utils.MarshalObject(route))
//	return &route, nil
//}
//
//func UpdateIngressRoute(ctx context.Context, projIngress *unstructured.Unstructured, route *model.IngressRoute, remove bool) error {
//
//	var err error
//
//	if err = RemoveResourceField(projIngress, "metadata.managedFields"); err != nil {
//		//		logger.Errorf(ctx, "Error removing metadata.managedFields - %s", err)
//		//		return errs.NewApplicationError(errs.ResourceGenerationError, err)
//		return err
//	}
//
//	if err = RemoveResourceField(projIngress, "spec.status"); err != nil {
//		//		logger.Errorf(ctx, "Error removing spec.status - %s", err)
//		//		return errs.NewApplicationError(errs.ResourceGenerationError, err)
//		return err
//	}
//
//	routeData := MarshalObject(GetResourceField(projIngress, "spec.routes"))
//	var routes []model.IngressRoute
//	if err = json.Unmarshal([]byte(routeData), &routes); err != nil {
//		//		logger.Errorf(ctx, "Error unmarshalling ingress routes - %s", err)
//		//		return errs.NewApplicationError(errs.ResourceGenerationError, err)
//		return err
//	}
//
//	var updated = false
//	for index, r := range routes {
//		for _, condition := range r.Conditions {
//			//			logger.Infof(ctx, "Comparing %s and %s at index %d ...", condition.Prefix, route.Conditions[0].Prefix, index)
//			if condition.Prefix == route.Conditions[0].Prefix {
//				if remove {
//					routes = append(routes[:index], routes[index+1:]...)
//				} else {
//					routes = append(routes[:index], *route)
//					routes = append(routes, routes[index+1:]...)
//				}
//				updated = true
//				break
//			}
//		}
//	}
//
//	if !updated {
//		routes = append(routes, *route)
//	}
//
//	//	logger.Debugf(ctx, "Ingress routes: %s", utils.MarshalObject(routes))
//	if err = SetResourceField(projIngress, "spec.routes", routes); err != nil {
//		//		logger.Errorf(ctx, "Error setting spec.status - %s", err)
//		//		return errs.NewApplicationError(errs.ResourceGenerationError, err)
//		return err
//	}
//
//	return nil
//}
//
//func GetResourceField(obj *unstructured.Unstructured, path string) interface{} {
//	var content = obj.UnstructuredContent()
//	parts := strings.Split(path, ".")
//	for index, part := range parts {
//		if index == len(parts)-1 {
//			return content[part]
//		} else {
//			entry := content[part]
//			if entry == nil {
//				return nil
//			}
//			content = entry.(map[string]interface{})
//		}
//	}
//	return nil
//}
//
//func ReadResourceField(obj *unstructured.Unstructured, path string, data interface{}) error {
//	value := GetResourceField(obj, path)
//	if value != nil {
//
//		valueBytes, err := json.Marshal(value)
//		if err != nil {
//			return err
//		}
//
//		return json.Unmarshal(valueBytes, &data)
//	}
//	return nil
//}
//
//func SetResourceField(obj *unstructured.Unstructured, path string, value interface{}) error {
//	var content = obj.UnstructuredContent()
//	parts := strings.Split(path, ".")
//	for index, part := range parts {
//		if index == len(parts)-1 {
//			content[part] = value
//		} else {
//			entry := content[part]
//			if entry == nil {
//				entry = make(map[string]interface{})
//				content[part] = entry
//			}
//			content = entry.(map[string]interface{})
//		}
//	}
//	return nil
//}
//
//func AddResourceField(obj *unstructured.Unstructured, path string, value interface{}) error {
//	var content = obj.UnstructuredContent()
//	parts := strings.Split(path, ".")
//	for index, part := range parts {
//		if index == len(parts)-1 {
//			var array []interface{}
//			entry := content[part]
//			if entry != nil {
//				array = entry.([]interface{})
//			}
//			array = append(array, value)
//			content[part] = array
//		} else {
//			entry := content[part]
//			if entry == nil {
//				entry = make(map[string]interface{})
//				content[part] = entry
//			}
//			content = entry.(map[string]interface{})
//		}
//	}
//	return nil
//}
//
//func RemoveResourceField(obj *unstructured.Unstructured, path string) error {
//	var content = obj.UnstructuredContent()
//	parts := strings.Split(path, ".")
//	for index, part := range parts {
//		if index == len(parts)-1 {
//			if content[part] != nil {
//				delete(content, part)
//			}
//		} else {
//			entry := content[part]
//			if entry == nil {
//				entry = make(map[string]interface{})
//				content[part] = entry
//			}
//			content = entry.(map[string]interface{})
//		}
//	}
//	return nil
//}
//
//// GetResourceProperties returns the corresponding property type for a kubernetes resource
//// returns map of data entries for ConfigMap
//// returns map of data entries (base64 decoded) for Secret
//// returns map of requested size, actual size, storage class name and volume name for PersistentVolumeClaim
//// returns an empty map for all other resources
//func GetResourceProperties(obj *unstructured.Unstructured) map[string]interface{} {
//	kind := obj.GetKind()
//
//	switch model.ResourceObjectType(kind) {
//	case model.ResourceConfigMap:
//		return GetResourceField(obj, "data").(map[string]interface{})
//
//	case model.ResourceSecret:
//		data := make(map[string]interface{}, 0)
//		for key, value := range GetResourceField(obj, "data").(map[string]interface{}) {
//			data[key] = Base64DecodeString(value.(string))
//		}
//		return data
//
//	case model.ResourcePersistentVolumeClaim:
//		return map[string]interface{}{
//			"requestedStorage": GetResourceField(obj, "spec.resources.requests.storage"),
//			"actualStorage":    GetResourceField(obj, "status.capacity.storage"),
//			"storageClassName": GetResourceField(obj, "spec.storageClassName"),
//			"volumeName":       GetResourceField(obj, "spec.volumeName"),
//		}
//
//	case model.ResourceDeployment:
//		_containers := GetResourceField(obj, "spec.template.spec.containers")
//		if _containers != nil {
//			containers := _containers.([]map[string]interface{})
//
//			return map[string]interface{}{
//				"image":     containers[0]["image"].(string),
//				"resources": containers[0]["resources"],
//			}
//		}
//	}
//
//	return map[string]interface{}{}
//}
//
//// GetResourceStatusField returns the corresponding status for a kubernetes resource
//// returns the status based on available replicas for a Deployment resource (active, partial or pending)
//// returns the phase for a PersistentVolumeClaim resource
//// returns active for all other resources
//func GetResourceStatusField(obj *unstructured.Unstructured) string {
//
//	kind := obj.GetKind()
//
//	status := "active"
//
//	switch model.ResourceObjectType(kind) {
//	case model.ResourceDeployment:
//		_availableReplicas := GetResourceField(obj, "status.availableReplicas")
//		_readyReplicas := GetResourceField(obj, "status.readyReplicas")
//		_replicas := GetResourceField(obj, "status.replicas")
//
//		var availableReplicas = 0
//		var readyReplicas = 0
//		var replicas = 1
//
//		if _availableReplicas != nil {
//			availableReplicas = _availableReplicas.(int)
//		}
//
//		if _readyReplicas != nil {
//			readyReplicas = _readyReplicas.(int)
//		}
//
//		if _replicas != nil {
//			replicas = _replicas.(int)
//		}
//
//		if replicas == readyReplicas && availableReplicas == replicas {
//			status = "active"
//		} else if readyReplicas > 0 && availableReplicas > 0 {
//			status = "partial"
//		} else {
//			status = "pending"
//		}
//
//	case model.ResourcePersistentVolumeClaim:
//		status = strings.ToLower(GetResourceField(obj, "status.phase").(string))
//
//	case model.ResourceVolumeSnapshot:
//		_status := GetResourceField(obj, "status.readyToUse")
//		if _status != nil {
//			if _status.(bool) {
//				status = "active"
//			}
//		} else {
//			status = "inactive"
//		}
//
//	case model.ResourceSnapshotSchedule:
//		status = ""
//
//	case model.ResourceHTTPProxy:
//		status = strings.ToLower(GetResourceField(obj, "status.currentStatus").(string))
//	}
//
//	return status
//}
