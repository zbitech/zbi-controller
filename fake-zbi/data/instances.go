package data

import (
	"github.com/zbitech/controller/pkg/model"
)

var (
	Instance1    = &model.Instance{Project: Project1.Name, Name: "instance1", InstanceType: model.InstanceTypeZCASH}
	Instance2    = &model.Instance{Project: Project1.Name, Name: "instance2", InstanceType: model.InstanceTypeZCASH}
	Instance3    = &model.Instance{Project: Project2.Name, Name: "instance3", InstanceType: model.InstanceTypeZCASH}
	Instance4    = &model.Instance{Project: Project2.Name, Name: "instance4", InstanceType: model.InstanceTypeZCASH}
	Instance5    = &model.Instance{Project: Project3.Name, Name: "instance5", InstanceType: model.InstanceTypeZCASH}
	Instance6    = &model.Instance{Project: Project3.Name, Name: "instance6", InstanceType: model.InstanceTypeZCASH}
	Instance7    = &model.Instance{Project: Project4.Name, Name: "instance7", InstanceType: model.InstanceTypeZCASH}
	Instance8    = &model.Instance{Project: Project4.Name, Name: "instance8", InstanceType: model.InstanceTypeZCASH}
	LwdInstance1 = &model.Instance{Project: Project1.Name, Name: "instance10", InstanceType: model.InstanceTypeLWD}
	LwdInstance2 = &model.Instance{Project: Project1.Name, Name: "instance11", InstanceType: model.InstanceTypeLWD}

	Instances       = []*model.Instance{Instance1, Instance2, Instance3, Instance4, Instance5, Instance6, Instance7, Instance8}
	Owner1Instances = []*model.Instance{Instance1, Instance2, Instance3, Instance4}
	Owner2Instances = []*model.Instance{Instance5, Instance6, Instance7, Instance8}

	Project1Instances = []*model.Instance{Instance1, Instance2}
	Project2Instances = []*model.Instance{Instance3, Instance4}
	Project3Instances = []*model.Instance{Instance6, Instance6}
	Project4Instances = []*model.Instance{Instance7, Instance8}
)

func AppendInstances(instances []*model.Instance, _instances ...*model.Instance) []*model.Instance {
	return append(instances, _instances...)
}

func CreateInstances(project *model.Project, count int, iType model.InstanceType, props map[string]interface{}) []*model.Instance {
	var instances = make([]*model.Instance, count)
	for index := range instances {
		instances[index] = CreateInstance(project, iType, props)
	}
	return instances
}

func CreateInstance(project *model.Project, iType model.InstanceType, props map[string]interface{}) *model.Instance {

	switch iType {
	case model.InstanceTypeZCASH:
		return CreateZcashInstance(project, props)
	case model.InstanceTypeLWD:
		return CreateLWDInstance(project, props)
	}

	return nil
}

func createInstance(project *model.Project, iType model.InstanceType, props map[string]interface{}) model.Instance {
	return model.Instance{
		Project:      project.Name,
		Name:         getProperty(props, "name", randomString(10)).(string),
		InstanceType: iType,
	}
}

func CreateZcashInstance(project *model.Project, props map[string]interface{}) *model.Instance {

	zcash := createInstance(project, model.InstanceTypeZCASH, props)

	return &zcash
}

func CreateLWDInstance(project *model.Project, props map[string]interface{}) *model.Instance {

	lwd := createInstance(project, model.InstanceTypeLWD, props)

	return &lwd
}
