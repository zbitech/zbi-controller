package zbi

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zbitech/controller/internal/helper"
	client "github.com/zbitech/controller/internal/klient/client"
	"github.com/zbitech/controller/internal/vars"
	"github.com/zbitech/controller/pkg/interfaces"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
	"strings"
	"time"
)

type ZBIClient struct {
	client interfaces.KlientIF
	//	informer interfaces.KlientInformerControllerIF
}

func NewZBIClient(client *client.Klient) interfaces.ZBIClientIF {
	return &ZBIClient{client: client}
}

func (z *ZBIClient) GetProjects(ctx context.Context) ([]model.Project, error) {
	var log = logger.GetServiceLogger(ctx, "zbi.GetProjects")
	defer func() { logger.LogServiceTime(log) }()

	labels := map[string]string{"platform": "zbi"}

	projects := make([]model.Project, 0)

	namespaces, err := z.client.GetNamespaces(ctx, labels)
	if err != nil {
		//TODO - check if err indicates server failure
		return nil, err
	}

	if namespaces != nil {
		for _, namespace := range namespaces {
			project := model.Project{
				Name:    namespace.Labels["project"],
				Network: model.NetworkType(namespace.Labels["network"]),
				Owner:   namespace.Labels["owner"],
				TeamId:  namespace.Labels["team"],
			}
			projects = append(projects, project)
		}
	}

	return projects, nil
}

func (z *ZBIClient) GetProject(ctx context.Context, project string) (*model.Project, error) {
	var log = logger.GetServiceLogger(ctx, "zbi.GetProjects")
	defer func() { logger.LogServiceTime(log) }()

	namespace, err := z.client.GetNamespace(ctx, project)
	if err != nil {
		//TODO - check if err indicates server failure
		return nil, err
	}

	return &model.Project{
		Name:    namespace.Labels["project"],
		Network: model.NetworkType(namespace.Labels["network"]),
		Owner:   namespace.Labels["owner"],
		TeamId:  namespace.Labels["team"],
	}, nil
}

// CreateProject creates the resources for a project
func (z *ZBIClient) CreateProject(ctx context.Context, project *model.Project) error {
	var log = logger.GetServiceLogger(ctx, "zbi.CreateProject")
	defer func() { logger.LogServiceTime(log) }()

	rscMgr := vars.ManagerFactory.GetProjectDataManager(ctx)

	objects, err := rscMgr.CreateProjectResource(ctx, project)
	if err != nil {
		log.Errorf("project kubernetes resource generation failed - %s", err)
		return err
	}

	resources, err := z.client.ApplyResources(ctx, objects)
	if err != nil {
		log.Errorf("project kubernetes resource creation failed - %s", err)
		return err
	}

	log.Infof("created %d resources for project %s", len(resources), project.Name)

	appIngress, _ := z.client.GetIngress(ctx, "zbi", "zbi-proxy")
	objects, err = rscMgr.CreateProjectIngressResource(ctx, appIngress, project, "add")
	if err != nil {
		log.Errorf("project ingress object creation failed - %s", err)
		return err
	} else {
		resources, err = z.client.ApplyResources(ctx, objects)
		if err != nil {
			log.Errorf("Project ingress resource creation failed - %s", err)
			return err
		}
	}

	project.Resources = make([]model.KubernetesResource, 0)
	project.AddResources(resources...)

	return nil
}

func (z *ZBIClient) RepairProject(ctx context.Context, project *model.Project) error {
	var log = logger.GetServiceLogger(ctx, "zbi.RepairProject")
	defer func() { logger.LogServiceTime(log) }()

	//TODO implement me
	panic("implement me")
}

func (z *ZBIClient) DeleteProject(ctx context.Context, project *model.Project, instances []model.Instance) error {
	var log = logger.GetServiceLogger(ctx, "zbi.DeleteProject")
	defer func() { logger.LogServiceTime(log) }()

	log.Infof("Delete project: %s", project.Name)
	if instances != nil {
		for _, instance := range instances {
			if err := z.DeleteInstance(ctx, project, &instance); err != nil {
				log.WithFields(logrus.Fields{"error": err, "instance": instance}).Errorf("failed to delete instance")
			}
		}
	}

	err := z.client.DeleteNamespace(ctx, project.Name)
	if err != nil {
		return err
	}

	return nil
}

func (z *ZBIClient) GetProjectResources(ctx context.Context, project string) ([]model.KubernetesResource, error) {
	var log = logger.GetServiceLogger(ctx, "zbi.GetProjectResources")
	defer func() { logger.LogServiceTime(log) }()

	labels := map[string]string{"platform": "zbi", "project": project}

	resources := make([]model.KubernetesResource, 0)

	cmaps, _ := z.client.GetConfigMaps(ctx, project, labels)
	if cmaps != nil {
		for _, object := range cmaps {

			cmap := helper.CreateCoreResource(ctx, model.ResourceConfigMap, &object, z.client)

			resources = append(resources, *cmap)
		}
	}

	secrets, _ := z.client.GetSecrets(ctx, project, labels)
	if secrets != nil {
		for _, object := range secrets {
			resources = append(resources, *helper.CreateCoreResource(ctx, model.ResourceSecret, &object, z.client))
		}
	}

	pvcs, _ := z.client.GetPersistentVolumeClaims(ctx, project, labels)
	if pvcs != nil {
		for _, object := range pvcs {
			resources = append(resources, *helper.CreateCoreResource(ctx, model.ResourcePersistentVolumeClaim, &object, z.client))
		}
	}

	deps := z.client.GetDeployments(ctx, project, labels)
	if deps != nil {
		for _, object := range deps {

			deployment := helper.CreateCoreResource(ctx, model.ResourceDeployment, &object, z.client)

			//pods := z.client.GetPods(ctx, project, object.Spec.Template.Labels)
			//podMap := make(map[string]interface{})
			//for _, p := range pods {
			//	podMap[p.Name] = map[string]interface{}{
			//		"status":    strings.ToLower(string(p.Status.Phase)),
			//		"startTime": p.Status.StartTime.Time,
			//	}
			//}
			//deployment.Properties["pods"] = podMap

			resources = append(resources, *deployment)
		}
	}

	svcs, _ := z.client.GetServices(ctx, project, labels)
	if svcs != nil {
		for _, object := range svcs {
			resources = append(resources, *helper.CreateCoreResource(ctx, model.ResourceService, &object, z.client))
		}
	}

	vs := z.client.GetVolumeSnapshots(ctx, project, labels)
	if vs != nil {
		for _, object := range vs {
			resources = append(resources, *helper.CreateCoreResource(ctx, model.ResourceVolumeSnapshot, &object, z.client))
		}
	}

	schs := z.client.GetSnapshotSchedules(ctx, project, labels)
	if schs != nil {
		for _, object := range schs {
			resources = append(resources, *helper.CreateCoreResource(ctx, model.ResourceSnapshotSchedule, &object, z.client))
		}
	}

	ingresses := z.client.GetIngresses(ctx, project, labels)
	if ingresses != nil {
		for _, object := range ingresses {
			resources = append(resources, *helper.CreateCoreResource(ctx, model.ResourceHTTPProxy, &object, z.client))
		}
	}

	return resources, nil
}

func (z *ZBIClient) GetProjectResource(ctx context.Context, project, resourceName string, resourceType model.ResourceObjectType) (*model.KubernetesResource, error) {

	var log = logger.GetServiceLogger(ctx, "zbi.GetProjectResource")
	defer func() { logger.LogServiceTime(log) }()

	object, err := z.client.GetDynamicResource(ctx, project, resourceName, helper.GvrMap[resourceType])
	if err != nil {
		return nil, err
	}

	var status = helper.GetResourceStatusField(object) //TODO - get status of resource
	var properties = helper.GetResourceProperties(object)
	created := time.Now()
	objType := model.ResourceObjectType(object.GetKind())

	resource := &model.KubernetesResource{
		Name: resourceName,
		//		Project:    project,
		Namespace: object.GetNamespace(),
		Type:      objType,
		Status:    status,
		Created:   &created,
		//		Updated:    &created,
		Properties: properties,
	}

	return resource, nil
}

func (z *ZBIClient) CreateInstance(ctx context.Context, project *model.Project, instance *model.Instance, request *model.ResourceRequest, peers ...*model.Instance) error {

	var log = logger.GetServiceLogger(ctx, "zbi.CreateInstance")
	defer func() { logger.LogServiceTime(log) }()

	dataMgr := vars.ManagerFactory.GetProjectDataManager(ctx)

	projectIngress, err := z.client.GetIngress(ctx, project.GetNamespace(), "project-ingress")
	if err != nil {
		log.WithFields(logrus.Fields{"error": err, "instance": instance}).Errorf("failed to create instance")
		return err
	}

	objects, err := dataMgr.CreateInstanceResource(ctx, projectIngress, project, instance, request, peers...)
	if err != nil {
		log.Errorf("instance kubernetes resource generation failed - %s", err)
		return err
	}

	resources, err := z.client.ApplyResources(ctx, objects[0])
	if err != nil {
		log.WithFields(logrus.Fields{"error": err, "instance": instance}).Errorf("instance kubernetes resource creation failed")
		return err
	}

	instance.Resources = make([]model.KubernetesResource, 0)
	instance.AddResources(resources...)

	if peers != nil && len(objects) > 1 {
		for index, _ := range peers {
			resources, err = z.client.ApplyResources(ctx, objects[index])
			if err != nil {
				log.WithFields(logrus.Fields{"error": err, "peer": peers[index]}).Errorf("peer kubernetes resource creation failed")
			} else {
				peers[index].AddResources(resources...)
			}
		}
	} else {
		log.Infof("no peer update needed")
	}

	return nil
}

func (z *ZBIClient) GetInstances(ctx context.Context, project string) ([]model.Instance, error) {

	log := logger.GetServiceLogger(ctx, "zbi.GetInstances")
	defer logger.LogServiceTime(log)

	labels := map[string]string{"platform": "zbi", "project": project}

	var instance *model.Instance
	var exists bool
	instanceMap := make(map[string]model.Instance)

	cmaps, _ := z.client.GetConfigMaps(ctx, project, labels)
	if cmaps != nil {
		for _, object := range cmaps {
			instance, exists, instanceMap = helper.GetInstance(instanceMap, object.Labels)
			if exists {
				instance.Resources = append(instance.Resources, *helper.CreateCoreResource(ctx, model.ResourceConfigMap, &object, z.client))
				instanceMap[instance.Name] = *instance
			}
		}
	}

	secrets, _ := z.client.GetSecrets(ctx, project, labels)
	if secrets != nil {
		for _, object := range secrets {
			instance, exists, instanceMap = helper.GetInstance(instanceMap, object.Labels)
			if exists {
				instance.Resources = append(instance.Resources, *helper.CreateCoreResource(ctx, model.ResourceSecret, &object, z.client))
				instanceMap[instance.Name] = *instance
			}
		}
	}

	pvcs, _ := z.client.GetPersistentVolumeClaims(ctx, project, labels)
	if pvcs != nil {
		for _, object := range pvcs {
			instance, exists, instanceMap = helper.GetInstance(instanceMap, object.Labels)
			if exists {
				instance.Resources = append(instance.Resources, *helper.CreateCoreResource(ctx, model.ResourcePersistentVolumeClaim, &object, z.client))
				instanceMap[instance.Name] = *instance
			}
		}
	}

	deps := z.client.GetDeployments(ctx, project, labels)
	if deps != nil {
		for _, object := range deps {

			instance, exists, instanceMap = helper.GetInstance(instanceMap, object.Labels)
			if exists {
				deployment := helper.CreateCoreResource(ctx, model.ResourceDeployment, &object, z.client)
				instance.Resources = append(instance.Resources, *deployment)
				instanceMap[instance.Name] = *instance
			}
		}
	}

	svcs, _ := z.client.GetServices(ctx, project, labels)
	if svcs != nil {
		for _, object := range svcs {
			instance, exists, instanceMap = helper.GetInstance(instanceMap, object.Labels)
			if exists {
				instance.Resources = append(instance.Resources, *helper.CreateCoreResource(ctx, model.ResourceService, &object, z.client))
				instanceMap[instance.Name] = *instance
			}
		}
	}

	vs := z.client.GetVolumeSnapshots(ctx, project, labels)
	if vs != nil {
		for _, object := range vs {
			instance, exists, instanceMap = helper.GetInstance(instanceMap, helper.GetResourceLabels(&object))
			if exists {
				instance.Resources = append(instance.Resources, *helper.CreateCoreResource(ctx, model.ResourceVolumeSnapshot, &object, z.client))
				instanceMap[instance.Name] = *instance
			}
		}
	}

	schs := z.client.GetSnapshotSchedules(ctx, project, labels)
	if schs != nil {
		for _, object := range schs {
			instance, exists, instanceMap = helper.GetInstance(instanceMap, helper.GetResourceLabels(&object))
			if exists {
				instance.Resources = append(instance.Resources, *helper.CreateCoreResource(ctx, model.ResourceSnapshotSchedule, &object, z.client))
				instanceMap[instance.Name] = *instance
			}
		}
	}

	ingresses := z.client.GetIngresses(ctx, project, labels)
	if ingresses != nil {
		for _, object := range ingresses {
			instance, exists, instanceMap = helper.GetInstance(instanceMap, helper.GetResourceLabels(&object))
			if exists {
				instance.Resources = append(instance.Resources, *helper.CreateCoreResource(ctx, model.ResourceHTTPProxy, &object, z.client))
				instanceMap[instance.Name] = *instance
			}
		}
	}

	instances := make([]model.Instance, 0)
	for _, instance := range instanceMap {
		instances = append(instances, instance)
	}

	return instances, nil
}

func (z *ZBIClient) GetInstance(ctx context.Context, project, instance string) (*model.Instance, error) {

	return nil, nil
}

func (z *ZBIClient) GetInstanceResources(ctx context.Context, project, instance string) ([]model.KubernetesResource, error) {

	var log = logger.GetServiceLogger(ctx, "zbi.GetInstanceResources")
	defer func() { logger.LogServiceTime(log) }()

	labels := map[string]string{"project": project, "instance": instance}

	resources := make([]model.KubernetesResource, 0)

	cmaps, _ := z.client.GetConfigMaps(ctx, project, labels)
	if cmaps != nil {
		for _, object := range cmaps {
			resources = append(resources, *helper.CreateCoreResource(ctx, model.ResourceConfigMap, &object, z.client))
		}
	}

	secrets, _ := z.client.GetSecrets(ctx, project, labels)
	if secrets != nil {
		for _, object := range secrets {
			resources = append(resources, *helper.CreateCoreResource(ctx, model.ResourceSecret, &object, z.client))
		}
	}

	pvcs, _ := z.client.GetPersistentVolumeClaims(ctx, project, labels)
	if pvcs != nil {
		for _, object := range pvcs {
			resources = append(resources, *helper.CreateCoreResource(ctx, model.ResourcePersistentVolumeClaim, &object, z.client))
		}
	}

	deps := z.client.GetDeployments(ctx, project, labels)
	if deps != nil {
		for _, object := range deps {

			deployment := helper.CreateCoreResource(ctx, model.ResourceDeployment, &object, z.client)

			//pods := z.client.GetPods(ctx, project, object.Spec.Template.Labels)
			//podMap := make(map[string]interface{})
			//for _, p := range pods {
			//	podMap[p.Name] = map[string]interface{}{
			//		"status":    strings.ToLower(string(p.Status.Phase)),
			//		"startTime": p.Status.StartTime.Time,
			//	}
			//}
			//deployment.Properties["pods"] = podMap

			resources = append(resources, *deployment)
		}
	}

	svcs, _ := z.client.GetServices(ctx, project, labels)
	if svcs != nil {
		for _, object := range svcs {
			resources = append(resources, *helper.CreateCoreResource(ctx, model.ResourceService, &object, z.client))
		}
	}

	vs := z.client.GetVolumeSnapshots(ctx, project, labels)
	if vs != nil {
		for _, object := range vs {
			resources = append(resources, *helper.CreateCoreResource(ctx, model.ResourceVolumeSnapshot, &object, z.client))
		}
	}

	schs := z.client.GetSnapshotSchedules(ctx, project, labels)
	if schs != nil {
		for _, object := range schs {
			resources = append(resources, *helper.CreateCoreResource(ctx, model.ResourceSnapshotSchedule, &object, z.client))
		}
	}

	ingresses := z.client.GetIngresses(ctx, project, labels)
	if ingresses != nil {
		for _, object := range ingresses {
			resources = append(resources, *helper.CreateCoreResource(ctx, model.ResourceHTTPProxy, &object, z.client))
		}
	}

	return resources, nil
}

func (z *ZBIClient) GetInstanceResource(ctx context.Context, project, instance, resourceName string, resourceType model.ResourceObjectType) (*model.KubernetesResource, error) {

	var log = logger.GetServiceLogger(ctx, "zbi.GetInstanceResource")
	defer func() { logger.LogServiceTime(log) }()

	object, err := z.client.GetDynamicResource(ctx, project, resourceName, helper.GvrMap[resourceType])
	if err != nil {
		return nil, err
	}

	resource, err := helper.CreateUnstructuredResource(object), nil
	if err != nil {
		return nil, err
	}

	if resourceType == model.ResourceDeployment {
		labels := helper.GetResourceField(object, "spec.template.labels").(map[string]string)
		pods := z.client.GetPods(ctx, project, labels)
		podMap := make(map[string]interface{})
		for _, p := range pods {
			podMap[p.Name] = map[string]interface{}{
				"status":    strings.ToLower(string(p.Status.Phase)),
				"startTime": p.Status.StartTime.Time,
			}
		}
		resource.Properties["pods"] = podMap
	}

	return resource, nil
}

func (z *ZBIClient) DeleteInstanceResource(ctx context.Context, project *model.Project, instance *model.Instance, resourceName string, resourceType model.ResourceObjectType) error {

	var log = logger.GetServiceLogger(ctx, "zbi.DeleteInstanceResource")
	defer func() { logger.LogServiceTime(log) }()

	if err := z.client.DeleteDynamicResource(ctx, project.GetNamespace(), resourceName, helper.GvrMap[resourceType]); err != nil {
		return err
	}

	return nil
}

func (z *ZBIClient) UpdateInstance(ctx context.Context, project *model.Project, instance *model.Instance, request *model.ResourceRequest, peers ...*model.Instance) error {

	var log = logger.GetServiceLogger(ctx, "zbi.UpdateInstance")
	defer func() { logger.LogServiceTime(log) }()

	dataMgr := vars.ManagerFactory.GetProjectDataManager(ctx)

	objects, err := dataMgr.CreateUpdateResource(ctx, project, instance, request, peers...)
	if err != nil {
		log.Errorf("instance kubernetes resource generation failed - %s", err)
		return err
	}

	resources, err := z.client.ApplyResources(ctx, objects[0])
	if err != nil {
		log.WithFields(logrus.Fields{"error": err, "instance": instance}).Errorf("instance kubernetes resource creation failed")
		return err
	}

	instance.Resources = make([]model.KubernetesResource, 0)
	instance.AddResources(resources...)

	if peers != nil && len(objects) > 1 {
		for index, _ := range peers {
			resources, err = z.client.ApplyResources(ctx, objects[index])
			if err != nil {
				log.WithFields(logrus.Fields{"error": err, "peer": peers[index]}).Errorf("peer kubernetes resource creation failed")
			} else {
				peers[index].AddResources(resources...)
			}
		}
	} else {
		log.Infof("no peer update needed")
	}

	return nil
}

func (z *ZBIClient) DeleteInstance(ctx context.Context, project *model.Project, instance *model.Instance) error {

	var log = logger.GetServiceLogger(ctx, "zbi.DeleteInstance")
	defer func() { logger.LogServiceTime(log) }()

	projIngress, err := z.client.GetIngress(ctx, project.GetNamespace(), "project-ingress")
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Errorf("failed to get project ingress")
		return err
	}

	if instance.Resources == nil || len(instance.Resources) == 0 {
		instance.Resources, err = z.GetInstanceResources(ctx, instance.Project, instance.Name)
		if err != nil {
			log.WithFields(logrus.Fields{"error": err}).Errorf("failed to get instance resources")
			return err
		}
	}

	resources, newResources, err := vars.ManagerFactory.GetProjectDataManager(ctx).CreateDeleteResource(ctx, projIngress, project, instance, instance.Resources)

	_, err = z.client.ApplyResources(ctx, newResources)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Errorf("failed to create new resources")
	}

	resources, err = z.client.DeleteResources(ctx, resources)
	instance.AddResources(resources...)
	return err
}

func (z *ZBIClient) RepairInstance(ctx context.Context, project *model.Project, instance *model.Instance) error {

	var log = logger.GetServiceLogger(ctx, "zbi.RepairInstance")
	defer func() { logger.LogServiceTime(log) }()

	projectIngress, err := z.client.GetIngress(ctx, project.GetNamespace(), "project-ingress")
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Errorf("failed to get project-ingress")
		return err
	}

	if instance.Resources == nil || len(instance.Resources) == 0 {
		instance.Resources, err = z.GetInstanceResources(ctx, instance.Project, instance.Name)
		if err != nil {
			log.WithFields(logrus.Fields{"error": err}).Errorf("failed to get instance resources")
			return err
		}
	}

	rscMgr := vars.ManagerFactory.GetProjectDataManager(ctx)

	objects, err := rscMgr.CreateRepairResource(ctx, projectIngress, project, instance)
	if err != nil {
		log.Errorf("instance kubernetes resource generation failed - %s", err)
		return err
	}

	resources, err := z.client.ApplyResources(ctx, objects)
	if err != nil {
		log.Errorf("instance kubernetes resource creation failed - %s", err)
		return err
	}

	instance.AddResources(resources...)
	log.Infof("created %d resources for instance %s in project %s", len(instance.Resources), instance.Name, instance.Project)

	return nil
}

func (z *ZBIClient) StopInstance(ctx context.Context, project *model.Project, instance *model.Instance) error {

	var log = logger.GetServiceLogger(ctx, "zbi.StopInstance")
	defer func() { logger.LogServiceTime(log) }()

	//	var evtErr error

	projectIngress, err := z.client.GetIngress(ctx, project.GetNamespace(), "project-ingress")
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Errorf("failed to get project-ingress")
		return err
	}

	if instance.Resources == nil || len(instance.Resources) == 0 {
		instance.Resources, err = z.GetInstanceResources(ctx, instance.Project, instance.Name)
		if err != nil {
			log.WithFields(logrus.Fields{"error": err}).Errorf("failed to get instance resources")
			return err
		}
	}

	rscMgr := vars.ManagerFactory.GetProjectDataManager(ctx)
	resources, objects, err := rscMgr.CreateStopResource(ctx, projectIngress, project, instance)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err, "instance": instance}).Errorf("failed to get stop resources")
		return err
	}

	resources, err = z.client.DeleteResources(ctx, resources)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err, "instance": instance}).Errorf("failed to delete resources for stop event")
	}
	instance.AddResources(resources...)

	resources, err = z.client.ApplyResources(ctx, objects)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err, "instance": instance}).Errorf("failed to create resources for stop event")
	}
	instance.AddResources(resources...)

	//objects, err := rscMgr.CreateStartResource(ctx, instance)
	//if err != nil {
	//	log.Errorf("instance kubernetes resource generation failed - %s", err)
	//	return err
	//}

	//for _, obj := range objects {
	//	objType := model.ResourceObjectType(obj.GetKind())
	//	if objType == model.ResourceDeployment {
	//		if err = z.client.DeleteDynamicResource(ctx, obj.GetNamespace(), obj.GetName(), helper.GvrMap[objType]); err != nil {
	//			log.Errorf("Failed to delete %s - %s", objType, err)
	//			if evtErr == nil {
	//				evtErr = err
	//			} else {
	//				evtErr = errors.Wrap(evtErr, err.Error())
	//			}
	//		}
	//
	//	} else if objType == model.ResourceService {
	//		// TODO map service to Stopped display by updating label match
	//	}
	//}

	//if evtErr != nil {
	//	return evtErr
	//}

	//obj, err := rscMgr.CreateIngressResource(ctx, projectIngress, instance, model.EventActionStopInstance)
	//if err != nil {
	//	log.Errorf("Instance ingress object creation failed - %s", err)
	//	return err
	//} else {
	//	_, err = z.client.ApplyResource(ctx, obj)
	//	if err != nil {
	//		log.Errorf("Instance ingress resource creation failed - %s", err)
	//		return err
	//	}
	//}

	return nil
}

// StartInstance creates the resource necessary to start the zcash instance and returns the resources
// that were created.
func (z *ZBIClient) StartInstance(ctx context.Context, project *model.Project, instance *model.Instance) error {

	var log = logger.GetServiceLogger(ctx, "zbi.StartInstance")
	defer func() { logger.LogServiceTime(log) }()

	projectIngress, err := z.client.GetIngress(ctx, project.GetNamespace(), "project-ingress")
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Errorf("failed to get project-ingress")
		return err
	}

	if instance.Resources == nil || len(instance.Resources) == 0 {
		instance.Resources, err = z.GetInstanceResources(ctx, instance.Project, instance.Name)
		if err != nil {
			log.WithFields(logrus.Fields{"error": err}).Errorf("failed to get instance resources")
			return err
		}
	}

	rscMgr := vars.ManagerFactory.GetProjectDataManager(ctx)
	objects, err := rscMgr.CreateStartResource(ctx, projectIngress, project, instance)
	if err != nil {
		log.Errorf("instance kubernetes resource generation failed - %s", err)
		return err
	}

	//obj, err := rscMgr.CreateIngressResource(ctx, projectIngress, project, instance, model.EventActionStartInstance)
	//if err != nil {
	//	log.Errorf("instance ingress object creation failed - %s", err)
	//	return err
	//}
	//
	//objects = append(objects, *obj)
	resources, err := z.client.ApplyResources(ctx, objects)

	if err != nil {
		log.Errorf("Instance kubernetes resource creation failed - %s", err)
		return err
	}

	instance.AddResources(resources...)
	log.Infof("Created %d resources for instance %s in project %s", len(resources), instance.Name, instance.Project)
	return nil
}

// RotateInstanceCredentials creates new resources to rotate the credentials associated with a Zcash instance and returns
// the resources that were created.
func (z *ZBIClient) RotateInstanceCredentials(ctx context.Context, project *model.Project, instance *model.Instance) error {

	var log = logger.GetServiceLogger(ctx, "zbi.RotateInstanceCredentials")
	defer func() { logger.LogServiceTime(log) }()

	var err error
	if instance.Resources == nil || len(instance.Resources) == 0 {
		instance.Resources, err = z.GetInstanceResources(ctx, instance.Project, instance.Name)
		if err != nil {
			log.WithFields(logrus.Fields{"error": err}).Errorf("failed to get instance resources")
			return err
		}
	}

	projMgr := vars.ManagerFactory.GetProjectDataManager(ctx)
	objects, err := projMgr.CreateRotationResource(ctx, project, instance)
	if err != nil {
		log.Errorf("instance rotation resource generation failed - %s", err)
		return err
	}

	resources, err := z.client.ApplyResources(ctx, objects)
	if err != nil {
		log.Errorf("instance rotation resource creation failed - %s", err)
		return err
	}

	instance.AddResources(resources...)
	log.Infof("created %d resources for instance %s in project %s", len(resources), instance.Name, instance.Project)

	return nil
}

// CreateSnapshot creates a new snapshot of the Zcash instance data volume.
func (z *ZBIClient) CreateSnapshot(ctx context.Context, project *model.Project, instance *model.Instance) error {

	var log = logger.GetServiceLogger(ctx, "zbi.CreateSnapshot")
	defer func() { logger.LogServiceTime(log) }()

	var err error
	if instance.Resources == nil || len(instance.Resources) == 0 {
		instance.Resources, err = z.GetInstanceResources(ctx, instance.Project, instance.Name)
		if err != nil {
			log.WithFields(logrus.Fields{"error": err}).Errorf("failed to get instance resources")
			return err
		}
	}

	projMgr := vars.ManagerFactory.GetProjectDataManager(ctx)
	objects, err := projMgr.CreateSnapshotResource(ctx, project, instance)

	if err != nil {
		log.Errorf("Failed to create snapshot assets - %s", err)
		return err
	}

	resources, err := z.client.ApplyResources(ctx, objects)
	if err != nil {
		log.Errorf("instance snapshot resource creation failed - %s", err)
		return err
	}

	instance.AddResources(resources...)
	log.Infof("created %d snapshot resources for instance %s in project %s", len(resources), instance.Name, instance.Project)

	return nil
}

// CreateSnapshotSchedule creates a new snapshot schedule for the Zcash instance data volume and returns the resources that were created.
func (z *ZBIClient) CreateSnapshotSchedule(ctx context.Context, project *model.Project, instance *model.Instance, schedule model.SnapshotScheduleType) error {

	var log = logger.GetServiceLogger(ctx, "zbi.CreateSnapshotSchedule")
	defer func() { logger.LogServiceTime(log) }()

	var err error
	if instance.Resources == nil || len(instance.Resources) == 0 {
		instance.Resources, err = z.GetInstanceResources(ctx, instance.Project, instance.Name)
		if err != nil {
			log.WithFields(logrus.Fields{"error": err}).Errorf("failed to get instance resources")
			return err
		}
	}

	projMgr := vars.ManagerFactory.GetProjectDataManager(ctx)
	objects, err := projMgr.CreateSnapshotScheduleResource(ctx, project, instance, schedule)

	if err != nil {
		log.Errorf("failed to create snapshot schedule assets - %s", err)
		return err
	}

	resources, err := z.client.ApplyResources(ctx, objects)
	if err != nil {
		log.Errorf("instance snapshot schedule resource creation failed - %s", err)
		return err
	}

	instance.AddResources(resources...)
	log.Infof("created %d snapshot schedule resources for instance %s in project %s", len(resources), instance.Name, instance.Project)

	return nil
}
