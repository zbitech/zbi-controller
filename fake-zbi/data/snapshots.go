package data

//func CreateInstantSnapshot(instance *model.Instance, props map[string]interface{}) *model.SnapshotResource {
//
//	return &model.SnapshotResource{
//		Name:     getProperty(props, "name", randomString(10)).(string),
//		Project:  instance.GetProject(),
//		Instance: instance.GetName(),
//		Created:  getProperty(props, "created", time.Now()).(time.Time),
//		Status:   getProperty(props, "status", randomValue(resourceStatus)).(string),
//	}
//}
//
//func CreateScheduledSnapshot(instance *model.Instance, props map[string]interface{}) *model.SnapshotResource {
//
//	return &model.SnapshotResource{
//		Name:     getProperty(props, "name", randomString(10)).(string),
//		Project:  instance.GetProject(),
//		Instance: instance.GetName(),
//		Created:  getProperty(props, "created", time.Now()).(time.Time),
//		Status:   getProperty(props, "status", randomValue(resourceStatus)).(string),
//	}
//}
