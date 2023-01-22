package manager

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/zbitech/controller/internal/helper"
	"github.com/zbitech/controller/internal/utils"
	"github.com/zbitech/controller/pkg/interfaces"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ProjectResourceManager struct {
	instances map[model.InstanceType]interfaces.InstanceResourceManagerIF
}

func NewProjectResourceManager(instanceManagers map[model.InstanceType]interfaces.InstanceResourceManagerIF) interfaces.ProjectResourceManagerIF {
	return &ProjectResourceManager{
		instances: instanceManagers,
	}
}

func (p ProjectResourceManager) CreateProjectResource(ctx context.Context, project *model.Project) ([]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "project.CreateProjectResource")
	defer func() { logger.LogServiceTime(log) }()

	fileTemplate := helper.Config.GetProjectTemplate()

	projectSpec := model.ProjectSpec{
		Name:         project.Name,
		Network:      project.Network,
		Owner:        project.Owner,
		TeamId:       project.TeamId,
		Namespace:    project.Name,
		Labels:       helper.CreateProjectLabels(project),
		InstancesMap: "",
	}

	var templates = []string{"NAMESPACE", "SERVICE"}
	specArr, err := fileTemplate.ExecuteTemplates(templates, projectSpec)
	if err != nil {
		log.Errorf("Project templates failed - %s", err)
		//logger.Errorf(ctx, "Project templates for version %s failed - %s", project.Version, err)
		//		return nil, errs.ErrProjectResourceFailed
		// return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	log.Debugf("Generated spec details - %s", specArr)
	return helper.CreateYAMLObjects(specArr)
}

func (p ProjectResourceManager) CreateProjectIngressResource(ctx context.Context, appIngress *unstructured.Unstructured, project *model.Project, action model.EventAction) ([]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "project.CreateProjectIngressResource")
	defer func() { logger.LogServiceTime(log) }()

	fileTemplate := helper.Config.GetProjectTemplate()

	projectSpec := model.ProjectSpec{
		Name:      project.Name,
		Network:   project.Network,
		Owner:     project.Owner,
		TeamId:    project.TeamId,
		Namespace: project.Name,
		Labels:    helper.CreateProjectLabels(project),
	}

	log.Debugf("Created project spec - %s", utils.MarshalObject(project))

	specObj, err := fileTemplate.ExecuteTemplates([]string{helper.INGRESS, helper.INGRESS_INCLUDE}, projectSpec)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Errorf("Project templates failed")
		return nil, err
	}

	log.Debugf("Generated spec details - %s", specObj)

	var ingressObj unstructured.Unstructured
	if err = helper.DecodeJSON(specObj[0], &ingressObj); err != nil {
		log.Errorf("Controller app template failed - %s", err)
		//		return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	var includeObj model.IngressInclude
	if err = json.Unmarshal([]byte(specObj[1]), &includeObj); err != nil {
		log.Errorf("Controller app template failed - %s", err)
		//		return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	//TODO - handle appIngress == nil - Return error?
	if err = helper.RemoveResourceField(appIngress, "metadata.managedFields"); err != nil {
		log.Errorf("unable to remove field metadata.managedFields")
	}

	if err = helper.RemoveResourceField(appIngress, "spec.status"); err != nil {
		log.Errorf("unable to remove field spec.status")
	}

	//TODO - add logic to update ingress if existing
	var includes []model.IngressInclude
	includeData := helper.GetResourceField(appIngress, "spec.includes")
	if includeData == nil {
		includes = make([]model.IngressInclude, 0)
	} else {
		data := includeData.([]interface{})
		dataBytes := new(bytes.Buffer)
		json.NewEncoder(dataBytes).Encode(data)
		if err = json.Unmarshal(dataBytes.Bytes(), &includes); err != nil {
			log.Errorf("Error unmarshalling ingress routes - %s", err)
		}
	}

	var updated = false
	for index, include := range includes {
		if include.Namespace == includeObj.Namespace {
			if action == model.EventActionDelete {
				includes = append(includes[:index], includes[index+1:]...)
			} else {
				includes = append(includes[:index], includeObj)
				includes = append(includes, includes[index+1:]...)
			}
			updated = true
		}
	}

	if !updated && action != "remove" {
		includes = append(includes, includeObj)
	}
	log.Debugf("Ingress includes - %s", utils.MarshalObject(includes))
	helper.SetResourceField(appIngress, "spec.includes", includes)

	appIng := *appIngress
	objIng := ingressObj

	return []unstructured.Unstructured{appIng, objIng}, nil
}

func (p ProjectResourceManager) CreateInstanceResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, peers ...model.Instance) ([]unstructured.Unstructured, [][]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "project.CreateInstanceResource")
	defer func() { logger.LogServiceTime(log) }()

	instanceManager, ok := p.instances[instance.InstanceType]
	if !ok {
		return nil, nil, errors.New("resource retrieval error")
	}

	project.Instances = append(project.Instances, model.Instance{
		Name:         instance.Name,
		InstanceType: instance.InstanceType,
		Request:      instance.Request,
	})

	fileTemplate := helper.Config.GetProjectTemplate()
	projectSpec := model.ProjectSpec{
		Namespace:    project.Name,
		InstancesMap: utils.MarshalObject(project.Instances),
		Labels:       helper.CreateProjectLabels(project),
	}

	specArr, err := fileTemplate.ExecuteTemplates([]string{helper.INSTANCE_LIST}, projectSpec)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Errorf("Project templates failed")
		return nil, nil, err
	}

	presources, err := helper.CreateYAMLObjects(specArr)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Errorf("failed to create instance list template")
		return nil, nil, err
	}

	iresources, err := instanceManager.CreateInstanceResource(ctx, projIngress, project, instance, peers...)
	if err != nil {
		return nil, nil, err
	}

	return presources, iresources, nil
}

func (p ProjectResourceManager) CreateUpdateResource(ctx context.Context, project *model.Project, instance *model.Instance, peers ...model.Instance) ([][]unstructured.Unstructured, error) {
	instanceManager, ok := p.instances[instance.InstanceType]
	if !ok {
		return nil, errors.New("resource retrieval error")
	}

	return instanceManager.CreateUpdateResource(ctx, project, instance, peers...)
}

func (p ProjectResourceManager) CreateStartResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error) {
	instanceManager, ok := p.instances[instance.InstanceType]
	if !ok {
		return nil, errors.New("resource retrieval error")
	}

	return instanceManager.CreateStartResource(ctx, projIngress, project, instance)
}

func (p ProjectResourceManager) CreateStopResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance) ([]model.KubernetesResource, []unstructured.Unstructured, error) {
	instanceManager, ok := p.instances[instance.InstanceType]
	if !ok {
		return nil, nil, errors.New("resource retrieval error")
	}

	return instanceManager.CreateStopResource(ctx, projIngress, project, instance)
}

func (p ProjectResourceManager) CreateRepairResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, peers ...model.Instance) ([]unstructured.Unstructured, error) {
	instanceManager, ok := p.instances[instance.InstanceType]
	if !ok {
		return nil, errors.New("resource retrieval error")
	}

	return instanceManager.CreateRepairResource(ctx, projIngress, project, instance, peers...)
}

func (p ProjectResourceManager) CreateIngressResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, action model.EventAction) (*unstructured.Unstructured, error) {
	instanceManager, ok := p.instances[instance.InstanceType]
	if !ok {
		return nil, errors.New("resource retrieval error")
	}

	return instanceManager.CreateIngressResource(ctx, projIngress, project, instance, action)
}

func (p ProjectResourceManager) CreateSnapshotResource(ctx context.Context, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error) {
	instanceManager, ok := p.instances[instance.InstanceType]
	if !ok {
		return nil, errors.New("resource retrieval error")
	}

	return instanceManager.CreateSnapshotResource(ctx, project, instance)
}

func (p ProjectResourceManager) CreateSnapshotScheduleResource(ctx context.Context, project *model.Project, instance *model.Instance, schedule model.SnapshotScheduleType) ([]unstructured.Unstructured, error) {
	instanceManager, ok := p.instances[instance.InstanceType]
	if !ok {
		return nil, errors.New("resource retrieval error")
	}

	return instanceManager.CreateSnapshotScheduleResource(ctx, project, instance, schedule)
}

func (p ProjectResourceManager) CreateRotationResource(ctx context.Context, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error) {
	instanceManager, ok := p.instances[instance.InstanceType]
	if !ok {
		return nil, errors.New("resource retrieval error")
	}

	return instanceManager.CreateRotationResource(ctx, project, instance)
}

func (p ProjectResourceManager) CreateDeleteResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, resources []model.KubernetesResource) ([]model.KubernetesResource, []unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "project.CreateDeleteResource")
	defer func() { logger.LogServiceTime(log) }()

	instanceManager, ok := p.instances[instance.InstanceType]
	if !ok {
		return nil, nil, errors.New("resource retrieval error")
	}

	project.Instances = append(project.Instances, *instance)
	for index, inst := range project.Instances {
		if inst.Name == instance.Name && inst.InstanceType == instance.InstanceType {
			log.WithFields(logrus.Fields{"instance": instance.Name}).Infof("removed from project instance list")
			project.Instances = append(project.Instances[:index], project.Instances[index+1:]...)
			break
		}
	}

	fileTemplate := helper.Config.GetProjectTemplate()
	projectSpec := model.ProjectSpec{
		Namespace:    project.Name,
		InstancesMap: utils.MarshalObject(project.Instances),
		Labels:       helper.CreateProjectLabels(project),
	}

	specArr, err := fileTemplate.ExecuteTemplates([]string{helper.INSTANCE_LIST}, projectSpec)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Errorf("Project templates failed")
		return nil, nil, err
	}

	newResources, err := helper.CreateYAMLObjects(specArr)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Errorf("failed to create instance list template")
	}

	delResources, newiResources, err := instanceManager.CreateDeleteResource(ctx, projIngress, project, instance, resources)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Errorf("failed to create resources to delete")
		return nil, nil, err
	}

	newResources = append(newResources, newiResources...)

	return delResources, newResources, nil
}
