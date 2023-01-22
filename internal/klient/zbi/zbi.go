package zbi

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zbitech/controller/internal/helper"
	client "github.com/zbitech/controller/internal/klient/client"
	"github.com/zbitech/controller/internal/utils"
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
			log.WithFields(logrus.Fields{"namespace": namespace}).Debugf("project details")

			project := model.Project{
				Name:      namespace.Labels["project"],
				Network:   model.NetworkType(namespace.Labels["network"]),
				Owner:     namespace.Labels["owner"],
				TeamId:    namespace.Labels["team"],
				Instances: make([]model.Instance, 0),
			}

			cm, err := z.client.GetConfigMapByName(ctx, project.GetNamespace(), "instances")
			if err != nil {
				log.WithFields(logrus.Fields{"error": err}).Errorf("failed to retrieve instances")
			} else {
				if err = utils.UnMarshalObject(cm.Data["instances"], &(project.Instances)); err != nil {
					log.WithFields(logrus.Fields{"error": err}).Errorf("failed to retrieve instances")
				}
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

	log.WithFields(logrus.Fields{"namespace": namespace}).Debugf("project details")

	var proj = &model.Project{
		Name:      namespace.Labels["project"],
		Network:   model.NetworkType(namespace.Labels["network"]),
		Owner:     namespace.Labels["owner"],
		TeamId:    namespace.Labels["team"],
		Instances: make([]model.Instance, 0),
	}

	cm, err := z.client.GetConfigMapByName(ctx, project, "instances")
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Errorf("failed to retrieve instances")
	} else {
		if err = utils.UnMarshalObject(cm.Data["instances"], &(proj.Instances)); err != nil {
			log.WithFields(logrus.Fields{"error": err}).Errorf("failed to retrieve instances")
		}
	}

	return proj, nil
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

func (z *ZBIClient) CreateInstance(ctx context.Context, project *model.Project, instance *model.Instance) error {

	var log = logger.GetServiceLogger(ctx, "zbi.CreateInstance")
	defer func() { logger.LogServiceTime(log) }()

	dataMgr := vars.ManagerFactory.GetProjectDataManager(ctx)

	var request = instance.Request
	peers, err := z.GetInstances(ctx, project, request.Peers)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Errorf("failed to get peer instances")
		return err
	}

	projectIngress, err := z.client.GetIngress(ctx, project.GetNamespace(), "project-ingress")
	if err != nil {
		log.WithFields(logrus.Fields{"error": err, "instance": instance}).Errorf("failed to create instance")
		return err
	}

	presources, objects, err := dataMgr.CreateInstanceResource(ctx, projectIngress, project, instance, peers...)
	if err != nil {
		log.Errorf("instance kubernetes resource generation failed - %s", err)
		return err
	}

	resources, err := z.client.ApplyResources(ctx, objects[0])
	if err != nil {
		log.WithFields(logrus.Fields{"error": err, "instance": instance}).Errorf("instance kubernetes resource creation failed")
		return err
	}

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

	_, err = z.client.ApplyResources(ctx, presources)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Errorf("project kubernetes resource creation failed")
		return err
	}

	return nil
}

func (z *ZBIClient) GetAllInstances(ctx context.Context, project *model.Project) ([]model.Instance, error) {

	log := logger.GetServiceLogger(ctx, "zbi.GetInstances")
	defer logger.LogServiceTime(log)

	instances := make([]model.Instance, 0)
	for _, inst := range project.Instances {
		instance, err := z.GetInstance(ctx, project, inst.Name)
		if err != nil {

		} else {
			instances = append(instances, *instance)
		}

	}

	return instances, nil
}

func (z *ZBIClient) GetInstances(ctx context.Context, project *model.Project, instanceNames []string) ([]model.Instance, error) {

	log := logger.GetServiceLogger(ctx, "zbi.GetInstance")
	defer logger.LogServiceTime(log)

	instances := make([]model.Instance, 0)
	for _, iname := range instanceNames {
		instance, err := z.GetInstance(ctx, project, iname)
		if err != nil {

		} else {
			instances = append(instances, *instance)
		}
	}

	return instances, nil
}

func (z *ZBIClient) GetInstance(ctx context.Context, project *model.Project, name string) (*model.Instance, error) {

	log := logger.GetServiceLogger(ctx, "zbi.GetInstance")
	defer logger.LogServiceTime(log)

	//	labels := map[string]string{"platform": "zbi", "project": project.Name, "instance": name}

	var resources *model.KubernetesResources
	var err error

	resources, err = z.GetInstanceResources(ctx, project, name)
	if err != nil {
		return nil, err
	}

	var instance = model.Instance{
		Name:         name,
		Project:      project.Name,
		InstanceType: project.GetInstanceType(name),
		Resources:    resources,
	}

	//cmaps, _ := z.client.GetConfigMaps(ctx, project.Name, labels)
	//if cmaps != nil {
	//	for _, object := range cmaps {
	//		instance.InstanceType = model.InstanceType(labels["type"])
	//		instance.Resources = append(instance.Resources, *helper.CreateCoreResource(ctx, model.ResourceConfigMap, &object, z.client))
	//	}
	//}
	//
	//secrets, _ := z.client.GetSecrets(ctx, project, labels)
	//if secrets != nil {
	//	for _, object := range secrets {
	//		instance.Resources = append(instance.Resources, *helper.CreateCoreResource(ctx, model.ResourceSecret, &object, z.client))
	//	}
	//}
	//
	//pvcs, _ := z.client.GetPersistentVolumeClaims(ctx, project, labels)
	//if pvcs != nil {
	//	for _, object := range pvcs {
	//		instance.Resources = append(instance.Resources, *helper.CreateCoreResource(ctx, model.ResourcePersistentVolumeClaim, &object, z.client))
	//	}
	//}
	//
	//deps := z.client.GetDeployments(ctx, project, labels)
	//if deps != nil {
	//	for _, object := range deps {
	//		deployment := helper.CreateCoreResource(ctx, model.ResourceDeployment, &object, z.client)
	//		instance.Resources = append(instance.Resources, *deployment)
	//	}
	//}
	//
	//svcs, _ := z.client.GetServices(ctx, project, labels)
	//if svcs != nil {
	//	for _, object := range svcs {
	//		instance.Resources = append(instance.Resources, *helper.CreateCoreResource(ctx, model.ResourceService, &object, z.client))
	//	}
	//}
	//
	//vs := z.client.GetVolumeSnapshots(ctx, project, labels)
	//if vs != nil {
	//	for _, object := range vs {
	//		instance.Resources = append(instance.Resources, *helper.CreateCoreResource(ctx, model.ResourceVolumeSnapshot, &object, z.client))
	//	}
	//}
	//
	//schs := z.client.GetSnapshotSchedules(ctx, project, labels)
	//if schs != nil {
	//	for _, object := range schs {
	//		instance.Resources = append(instance.Resources, *helper.CreateCoreResource(ctx, model.ResourceSnapshotSchedule, &object, z.client))
	//	}
	//}
	//
	//ingresses := z.client.GetIngresses(ctx, project, labels)
	//if ingresses != nil {
	//	for _, object := range ingresses {
	//		instance.Resources = append(instance.Resources, *helper.CreateCoreResource(ctx, model.ResourceHTTPProxy, &object, z.client))
	//	}
	//}

	return &instance, nil
}

func (z *ZBIClient) GetInstanceResources(ctx context.Context, project *model.Project, instance string) (*model.KubernetesResources, error) {

	var log = logger.GetServiceLogger(ctx, "zbi.GetInstanceResources")
	defer func() { logger.LogServiceTime(log) }()

	labels := map[string]string{"project": project.Name, "instance": instance}

	resources := model.KubernetesResources{
		Resources: make([]model.KubernetesResource, 0),
		Snapshots: make([]model.KubernetesResource, 0),
		Schedules: make([]model.KubernetesResource, 0),
	}

	cmaps, _ := z.client.GetConfigMaps(ctx, project.GetNamespace(), labels)
	if cmaps != nil {
		for _, object := range cmaps {
			resources.Resources = append(resources.Resources, *helper.CreateCoreResource(ctx, model.ResourceConfigMap, &object, z.client))
		}
	}

	secrets, _ := z.client.GetSecrets(ctx, project.GetNamespace(), labels)
	if secrets != nil {
		for _, object := range secrets {
			resources.Resources = append(resources.Resources, *helper.CreateCoreResource(ctx, model.ResourceSecret, &object, z.client))
		}
	}

	pvcs, _ := z.client.GetPersistentVolumeClaims(ctx, project.GetNamespace(), labels)
	if pvcs != nil {
		for _, object := range pvcs {
			resources.Resources = append(resources.Resources, *helper.CreateCoreResource(ctx, model.ResourcePersistentVolumeClaim, &object, z.client))
		}
	}

	deps := z.client.GetDeployments(ctx, project.GetNamespace(), labels)
	if deps != nil {
		for _, object := range deps {
			deployment := helper.CreateCoreResource(ctx, model.ResourceDeployment, &object, z.client)
			resources.Resources = append(resources.Resources, *deployment)
		}
	}

	svcs, _ := z.client.GetServices(ctx, project.GetNamespace(), labels)
	if svcs != nil {
		for _, object := range svcs {
			resources.Resources = append(resources.Resources, *helper.CreateCoreResource(ctx, model.ResourceService, &object, z.client))
		}
	}

	vs := z.client.GetVolumeSnapshots(ctx, project.GetNamespace(), labels)
	if vs != nil {
		for _, object := range vs {
			resources.Snapshots = append(resources.Snapshots, *helper.CreateCoreResource(ctx, model.ResourceVolumeSnapshot, &object, z.client))
		}
	}

	schs := z.client.GetSnapshotSchedules(ctx, project.GetNamespace(), labels)
	if schs != nil {
		for _, object := range schs {
			resources.Schedules = append(resources.Schedules, *helper.CreateCoreResource(ctx, model.ResourceSnapshotSchedule, &object, z.client))
		}
	}

	ingresses := z.client.GetIngresses(ctx, project.GetNamespace(), labels)
	if ingresses != nil {
		for _, object := range ingresses {
			resources.Resources = append(resources.Resources, *helper.CreateCoreResource(ctx, model.ResourceHTTPProxy, &object, z.client))
		}
	}

	return &resources, nil
}

func (z *ZBIClient) GetInstanceResource(ctx context.Context, project *model.Project, instance, resourceName string, resourceType model.ResourceObjectType) (*model.KubernetesResource, error) {

	var log = logger.GetServiceLogger(ctx, "zbi.GetInstanceResource")
	defer func() { logger.LogServiceTime(log) }()

	object, err := z.client.GetDynamicResource(ctx, project.GetNamespace(), resourceName, helper.GvrMap[resourceType])
	if err != nil {
		return nil, err
	}

	resource, err := helper.CreateUnstructuredResource(object), nil
	if err != nil {
		return nil, err
	}

	if resourceType == model.ResourceDeployment {
		labels := helper.GetResourceField(object, "spec.template.labels").(map[string]string)
		pods := z.client.GetPods(ctx, project.GetNamespace(), labels)
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

func (z *ZBIClient) UpdateInstance(ctx context.Context, project *model.Project, instance *model.Instance) error {

	var log = logger.GetServiceLogger(ctx, "zbi.UpdateInstance")
	defer func() { logger.LogServiceTime(log) }()

	dataMgr := vars.ManagerFactory.GetProjectDataManager(ctx)

	var request = instance.Request
	peers, err := z.GetInstances(ctx, project, request.Peers)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Errorf("failed to get peer instances")
		return err
	}

	objects, err := dataMgr.CreateUpdateResource(ctx, project, instance, peers...)
	if err != nil {
		log.Errorf("instance kubernetes resource generation failed - %s", err)
		return err
	}

	resources, err := z.client.ApplyResources(ctx, objects[0])
	if err != nil {
		log.WithFields(logrus.Fields{"error": err, "instance": instance}).Errorf("instance kubernetes resource creation failed")
		return err
	}

	//	instance.Resources = make([]model.KubernetesResource, 0)
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

	if !instance.HasResources() {
		instance.Resources, err = z.GetInstanceResources(ctx, project, instance.Name)
		if err != nil {
			log.WithFields(logrus.Fields{"error": err}).Errorf("failed to get instance resources")
			return err
		}
	}

	resources, newResources, err := vars.ManagerFactory.GetProjectDataManager(ctx).CreateDeleteResource(ctx, projIngress, project, instance, instance.GetResources())

	log.Infof("apply new instance resources for deleted instance")

	log.Infof("removing instance resources")
	resources, err = z.client.DeleteResources(ctx, resources)
	instance.AddResources(resources...)

	_, err = z.client.ApplyResources(ctx, newResources)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Errorf("failed to create new resources")
	}

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

	if !instance.HasResources() {
		instance.Resources, err = z.GetInstanceResources(ctx, project, instance.Name)
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

	newResources, err := z.client.ApplyResources(ctx, objects)
	if err != nil {
		log.Errorf("instance kubernetes resource creation failed - %s", err)
		return err
	}

	instance.AddResources(newResources...)

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

	if !instance.HasResources() {
		instance.Resources, err = z.GetInstanceResources(ctx, project, instance.Name)
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

	if !instance.HasResources() {
		instance.Resources, err = z.GetInstanceResources(ctx, project, instance.Name)
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
	if !instance.HasResources() {
		instance.Resources, err = z.GetInstanceResources(ctx, project, instance.Name)
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
	if !instance.HasResources() {
		instance.Resources, err = z.GetInstanceResources(ctx, project, instance.Name)
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
	if !instance.HasResources() {
		instance.Resources, err = z.GetInstanceResources(ctx, project, instance.Name)
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
