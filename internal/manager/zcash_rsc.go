package manager

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/zbitech/controller/internal/helper"
	"github.com/zbitech/controller/internal/utils"
	"github.com/zbitech/controller/internal/vars"
	"github.com/zbitech/controller/pkg/interfaces"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"strconv"
	"text/template"
)

var (
	FUNCTIONS = template.FuncMap{
		"base64Encode": func(value string) string {
			return utils.Base64EncodeString(value)
		},
		"basicCredentials": func(username, password string) string {
			creds := fmt.Sprintf("%s:%s", username, password)
			return utils.Base64EncodeString(creds)
		},
	}
)

type ZcashInstanceResourceManager struct {
}

func NewZcashInstanceResourceManager() interfaces.InstanceResourceManagerIF {
	return &ZcashInstanceResourceManager{}
}

func (z *ZcashInstanceResourceManager) CreateInstanceResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, peers ...model.Instance) ([][]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "zcash.CreateInstanceResource")
	defer func() { logger.LogServiceTime(log) }()

	var specArr []string
	var err error

	policy := helper.Config.GetPolicyConfig()
	ic, err := helper.Config.GetInstanceConfig(instance.InstanceType)
	if err != nil {
		return nil, err
	}

	fileTemplate, err := helper.Config.GetInstanceTemplate(instance.InstanceType)
	if err != nil {
		return nil, err
	}

	var request = instance.Request
	var miner = false
	//	txIndex := instance.Properties["transactionIndex"].(bool)
	if _, ok := request.Properties[MINER_ZCASH_PROPERTY]; ok {
		miner = request.Properties[MINER_ZCASH_PROPERTY].(bool)
	}
	//	peers := instance.Properties["peers"].([]interface{})

	rpcport := strconv.FormatInt(int64(ic.GetPort(SERVICE_PORT)), 10)
	conf := createZcashConf(ic, miner, project.Network, rpcport)
	conf = getZcashPeers(conf, rpcport, project.GetNamespace(), peers...)

	dataVolumeName := fmt.Sprintf("%s-%s", instance.Name, utils.GenerateRandomString(5, true))
	dataVolumeSize := request.VolumeSize

	username := utils.GenerateRandomString(6, true)
	password := utils.GenerateSecurePassword()

	instanceSpec := model.InstanceSpec{
		Name:               instance.Name,
		Project:            instance.Project,
		ServiceAccountName: policy.ServiceAccount,
		Namespace:          project.GetNamespace(),
		Labels:             helper.CreateInstanceLabels(instance),
		DomainName:         policy.DomainName,
		DomainSecret:       policy.CertificateName,
		Envoy:              helper.CreateEnvoySpec(ic.GetPort(ENVOY_PORT)),
		DataVolumeName:     dataVolumeName,
		Images: map[string]string{
			ZCASH:   ic.GetImageRepository(NODE_IMAGE),
			METRICS: ic.GetImageRepository(METRICS_IMAGE),
		},
		Ports: map[string]int32{
			ZCASH:   ic.GetPort(SERVICE_PORT),
			METRICS: ic.GetPort(METRICS_PORT),
		},
		Properties: map[string]interface{}{
			USERNAME:                  username,
			PASSWORD:                  password,
			ZcashConf:                 conf,
			RESOURCE_REQUEST_PROPERTY: utils.MarshalObject(request),
		},
	}

	specArr, err = fileTemplate.ExecuteTemplates([]string{ZCASH_CONF, ENVOY_CONF, CREDENTIALS, DEPLOYMENT, SERVICE}, instanceSpec)

	if err != nil {
		log.Errorf("zcash templates failed - %s", err)
		//		return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	objects, err := helper.CreateYAMLObjects(specArr /*, instance.Project, instance.Name*/)
	if err != nil {
		log.Errorf("zcash templates failed - %s", err)
		//		return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	storageClass := policy.StorageClass

	var volumeSpecs = []model.VolumeSpec{
		{VolumeName: dataVolumeName, StorageClass: storageClass, Namespace: project.GetNamespace(),
			VolumeDataType: string(request.VolumeType), DataSourceType: request.VolumeSourceType,
			SourceName: request.VolumeSourceName,
			Size:       dataVolumeSize, Labels: instanceSpec.Labels},
	}

	appRsc := vars.ManagerFactory.GetAppResourceManager(ctx)
	volumes, err := appRsc.CreateVolumeResource(ctx, instance.Project, instance.Name, volumeSpecs...)
	if err != nil {
		log.Errorf("Zcash volume templates failed - %s", err)
		// return nil, errs.NewApplicationError(errs.ResourceGenerationError, err)
		return nil, err
	}

	objects = append(objects, volumes...)

	obj, err := z.CreateIngressResource(ctx, projIngress, project, instance, model.EventActionCreate)
	if err != nil {
		log.Errorf("Instance ingress object creation failed - %s", err)
		return nil, err
	}

	objects = append(objects, *obj)

	var resources = make([][]unstructured.Unstructured, 0)
	resources = append(resources, objects)

	for index, peer := range peers {
		instances := make([]model.Instance, 0)
		instances = append(instances, *instance)
		if index > 0 {
			instances = append(instances, peers[:index]...)
		}
		if index < len(peers) {
			instances = append(instances, peers[index+1:]...)
		}

		object, err := z.CreateUpdatePeersResource(ctx, project, &peer, instances...)
		if err != nil {
			log.WithFields(logrus.Fields{"error": err, "instance": peer}).Error("failed to update ")
		} else {
			//			objects = append(objects, object...)
			resources = append(resources, object)
		}
	}

	return resources, nil
}

func (z *ZcashInstanceResourceManager) CreateUpdatePeersResource(ctx context.Context, project *model.Project, instance *model.Instance, peers ...model.Instance) ([]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "zcash.CreateUpdatePeersResource")
	defer func() { logger.LogServiceTime(log) }()

	var err error
	fileTemplate, err := helper.Config.GetInstanceTemplate(instance.InstanceType)
	if err != nil {
		return nil, err
	}

	ic, err := helper.Config.GetInstanceConfig(instance.InstanceType)
	if err != nil {
		return nil, err
	}

	rpcport := strconv.FormatInt(int64(ic.GetPort("service")), 10)

	var request = instance.Request

	//request, ok := helper.GetResourceRequest(ctx, instance)
	//if !ok {
	//	return nil, errors.New("unable to retrieve resource request for instance")
	//}

	miner := request.Properties["miner"].(bool)
	conf := createZcashConf(ic, miner, project.Network, rpcport)
	conf = getZcashPeers(conf, rpcport, project.GetNamespace(), peers...)

	instanceSpec := model.InstanceSpec{
		Name:    instance.Name,
		Project: instance.Project,
		//		Version:   instance.Version,
		Namespace: project.GetNamespace(),
		Labels:    helper.CreateInstanceLabels(instance),
		Properties: map[string]interface{}{
			ZcashConf: conf,
		},
	}

	specArr, err := fileTemplate.ExecuteTemplates([]string{ZCASH_CONF}, instanceSpec)
	if err != nil {
		log.Errorf("zcash templates failed - %s", err)
		//		return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	return helper.CreateYAMLObjects(specArr)
}

func (z *ZcashInstanceResourceManager) CreateUpdateResource(ctx context.Context, project *model.Project, instance *model.Instance, peers ...model.Instance) ([][]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "zcash.CreateUpdateResource")
	defer func() { logger.LogServiceTime(log) }()

	var err error
	fileTemplate, err := helper.Config.GetInstanceTemplate(instance.InstanceType)
	if err != nil {
		return nil, err
	}

	ic, err := helper.Config.GetInstanceConfig(instance.InstanceType)
	if err != nil {
		return nil, err
	}

	rpcport := strconv.FormatInt(int64(ic.GetPort("service")), 10)

	var request = instance.Request

	miner := request.Properties["miner"].(bool)
	conf := createZcashConf(ic, miner, project.Network, rpcport)
	conf = getZcashPeers(conf, rpcport, project.GetNamespace(), peers...)

	//connect := make([]string, 0)
	//if peers != nil {
	//	for _, peer := range peers {
	//		connect = append(connect, peer.Name)
	//		conf = append(conf, object.KVPair{Key: "connect", Value: getZcashInstanceHost(peer) + ":" + rpcport})
	//		//conf = updateZcashPeers(conf, ic.GetPort("service"), peers...)
	//	}
	//	instance.Properties["peers"] = connect
	//}

	instanceSpec := model.InstanceSpec{
		Name:    instance.Name,
		Project: instance.Project,
		//		Version:   instance.Version,
		Namespace: project.GetNamespace(),
		Labels:    helper.CreateInstanceLabels(instance),
		Properties: map[string]interface{}{
			ZcashConf: conf,
		},
	}

	specArr, err := fileTemplate.ExecuteTemplates([]string{ZCASH_CONF}, instanceSpec)
	if err != nil {
		log.Errorf("zcash templates failed - %s", err)
		//		return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	objects, err := helper.CreateYAMLObjects(specArr /*, instance.Project, instance.Name*/)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err, "instance": instance}).Errorf("failed to create kubernetes resources")
		return nil, err
	}

	var resources = make([][]unstructured.Unstructured, 0)
	resources = append(resources, objects)

	for index, peer := range peers {
		instances := make([]model.Instance, 0)
		instances = append(instances, *instance)
		if index > 0 {
			instances = append(instances, peers[:index]...)
		}
		if index < len(peers) {
			instances = append(instances, peers[index+1:]...)
		}

		object, err := z.CreateUpdatePeersResource(ctx, project, &peer, instances...)
		if err != nil {
			log.WithFields(logrus.Fields{"error": err, "instance": peer}).Error("failed to update ")
		} else {
			//			objects = append(objects, object...)
			resources = append(resources, object)
		}
	}

	return resources, nil
}

func (z *ZcashInstanceResourceManager) CreateIngressResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, action model.EventAction) (*unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "zcash.CreateIngressResource")
	defer func() { logger.LogServiceTime(log) }()

	var err error
	fileTemplate, err := helper.Config.GetInstanceTemplate(instance.InstanceType)
	if err != nil {
		return nil, err
	}

	policy := helper.Config.GetPolicyConfig()
	ic, err := helper.Config.GetInstanceConfig(instance.InstanceType)
	if err != nil {
		return nil, err
	}

	zcashSpec := model.InstanceSpec{
		Name:               instance.Name,
		ServiceAccountName: policy.ServiceAccount,
		Namespace:          project.GetNamespace(),
		Labels:             helper.CreateInstanceLabels(instance),
		DomainName:         policy.DomainName,
		DomainSecret:       policy.CertificateName,
		Envoy:              helper.CreateEnvoySpec(ic.GetPort(ENVOY_PORT)),
	}
	var specObj string

	if action == model.EventActionStopInstance {
		specObj, err = fileTemplate.ExecuteTemplate(INGRESS_STOPPED, zcashSpec)
	} else {
		specObj, err = fileTemplate.ExecuteTemplate(INGRESS, zcashSpec)
	}

	if err != nil {
		log.Errorf("Zcash templates failed - %s", err)
		//		return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	var route *model.IngressRoute

	route, err = helper.CreateIngressRoute(ctx, specObj)
	if err != nil {
		log.Errorf("zcash route marshal failed - %s", err)
		// return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	if err = helper.UpdateIngressRoute(ctx, projIngress, route, action == model.EventActionDelete); err != nil {
		log.Errorf("error updating ingress route for zcash instance - %s", err)
		// return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	return projIngress, nil
}

func (z *ZcashInstanceResourceManager) CreateStartResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "zcash.CreateStartResource")
	defer func() { logger.LogServiceTime(log) }()

	pvc := instance.GetResourceByType(model.ResourcePersistentVolumeClaim)

	var err error
	fileTemplate, err := helper.Config.GetInstanceTemplate(instance.InstanceType)
	if err != nil {
		return nil, err
	}

	policy := helper.Config.GetPolicyConfig()
	ic, err := helper.Config.GetInstanceConfig(instance.InstanceType)
	if err != nil {
		return nil, err
	}

	zcashSpec := model.InstanceSpec{
		Name:               instance.Name,
		ServiceAccountName: policy.ServiceAccount,
		Namespace:          project.GetNamespace(),
		Labels:             helper.CreateInstanceLabels(instance),
		DomainName:         policy.DomainName,
		DomainSecret:       policy.CertificateName,
		Envoy:              helper.CreateEnvoySpec(ic.GetPort(ENVOY_PORT)),
		DataVolumeName:     pvc.Name,
		Images: map[string]string{
			ZCASH:   ic.GetImageRepository(NODE_IMAGE),
			METRICS: ic.GetImageRepository(METRICS_IMAGE),
		},
		Ports: map[string]int32{
			ZCASH:   ic.GetPort(SERVICE_PORT),
			METRICS: ic.GetPort(METRICS_PORT),
		},
	}
	specArr, err := fileTemplate.ExecuteTemplates([]string{DEPLOYMENT, SERVICE}, zcashSpec)
	if err != nil {
		log.Errorf("zcash templates failed - %s", err)
		// return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	objects, err := helper.CreateYAMLObjects(specArr /*, instance.Project, instance.Name*/)

	obj, err := z.CreateIngressResource(ctx, projIngress, project, instance, model.EventActionStartInstance)
	if err != nil {
		log.Errorf("instance ingress object creation failed - %s", err)
		return nil, err
	}

	objects = append(objects, *obj)

	return objects, nil
}

func (z *ZcashInstanceResourceManager) CreateStopResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance) ([]model.KubernetesResource, []unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "zcash.CreateStopResource")
	defer func() { logger.LogServiceTime(log) }()

	var resources = make([]model.KubernetesResource, 0)
	var objects = make([]unstructured.Unstructured, 0)

	deployment := instance.GetResourceByType(model.ResourceDeployment)
	if deployment != nil && deployment.Status == "active" {
		resources = append(resources, *deployment)
	}

	if len(resources) == 0 {
		return nil, nil, errors.New("instance is not active")
	}

	obj, err := z.CreateIngressResource(ctx, projIngress, project, instance, model.EventActionStopInstance)
	if err != nil {
		return nil, nil, err
	}
	objects = append(objects, *obj)

	return resources, objects, nil
}

func (z *ZcashInstanceResourceManager) CreateRepairResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, peers ...model.Instance) ([]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "zcash.CreateRepairResource")
	defer func() { logger.LogServiceTime(log) }()

	var specArr []string
	var err error

	policy := helper.Config.GetPolicyConfig()
	ic, err := helper.Config.GetInstanceConfig(instance.InstanceType)
	if err != nil {
		return nil, err
	}

	fileTemplate, err := helper.Config.GetInstanceTemplate(instance.InstanceType)
	if err != nil {
		return nil, err
	}

	var username, password, dataVolumeName, dataVolumeSize string

	pvc := instance.GetResourceByType(model.ResourcePersistentVolumeClaim)
	secret := instance.GetResourceByType(model.ResourceSecret)

	//request, ok := helper.GetResourceRequest(ctx, instance)
	//if !ok {
	//	return nil, errors.New("unable to retrieve resource request for instance")
	//}
	var request = instance.Request

	miner := request.Properties[MINER_ZCASH_PROPERTY].(bool)
	rpcport := strconv.FormatInt(int64(ic.GetPort(SERVICE_PORT)), 10)
	conf := createZcashConf(ic, miner, project.Network, rpcport)
	conf = getZcashPeers(conf, rpcport, project.GetNamespace(), peers...)

	if pvc != nil && pvc.Status == "active" {
		dataVolumeName = pvc.Name
		dataVolumeSize = request.VolumeSize
	} else {
		dataVolumeName = fmt.Sprintf("%s-%s", instance.Name, utils.GenerateRandomString(5, true))
		dataVolumeSize = request.VolumeSize
		pvc.Name = dataVolumeName // re-create in-active volume
	}

	if secret != nil && secret.Status == "active" {
		username = secret.Properties["username"].(string)
		password = secret.Properties["password"].(string)
	} else {
		username = utils.GenerateRandomString(6, true)
		password = utils.GenerateSecurePassword()
	}

	instanceSpec := model.InstanceSpec{
		Name:               instance.Name,
		Project:            instance.Project,
		ServiceAccountName: policy.ServiceAccount,
		Namespace:          project.GetNamespace(),
		Labels:             helper.CreateInstanceLabels(instance),
		DomainName:         policy.DomainName,
		DomainSecret:       policy.CertificateName,
		Envoy:              helper.CreateEnvoySpec(ic.GetPort(ENVOY_PORT)),
		DataVolumeName:     dataVolumeName,
		Images: map[string]string{
			ZCASH:   ic.GetImageRepository(NODE_IMAGE),
			METRICS: ic.GetImageRepository(METRICS_IMAGE),
		},
		Ports: map[string]int32{
			ZCASH:   ic.GetPort(SERVICE_PORT),
			METRICS: ic.GetPort(METRICS_PORT),
		},
		Properties: map[string]interface{}{
			USERNAME:  username,
			PASSWORD:  password,
			ZcashConf: conf,
		},
	}

	specArr, err = fileTemplate.ExecuteTemplates([]string{ZCASH_CONF, ENVOY_CONF, CREDENTIALS, DEPLOYMENT, SERVICE}, instanceSpec)

	if err != nil {
		log.Errorf("zcash templates failed - %s", err)
		//		return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	objects, err := helper.CreateYAMLObjects(specArr /*, instance.Project, instance.Name*/)
	if err != nil {
		log.Errorf("zcash templates failed - %s", err)
		//		return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	if pvc == nil || pvc.Status != "active" {
		storageClass := policy.StorageClass
		var volumeSpecs = []model.VolumeSpec{
			{VolumeName: dataVolumeName, StorageClass: storageClass, Namespace: project.GetNamespace(),
				VolumeDataType: string(request.VolumeType), DataSourceType: request.VolumeSourceType,
				SourceName: request.VolumeSourceName,
				Size:       dataVolumeSize, Labels: instanceSpec.Labels},
		}

		appRsc := vars.ManagerFactory.GetAppResourceManager(ctx)
		volumes, err := appRsc.CreateVolumeResource(ctx, instance.Project, instance.Name, volumeSpecs...)
		if err != nil {
			log.Errorf("Zcash volume templates failed - %s", err)
			// return nil, errs.NewApplicationError(errs.ResourceGenerationError, err)
			return nil, err
		}

		objects = append(objects, volumes...)
	}

	obj, err := z.CreateIngressResource(ctx, projIngress, project, instance, model.EventActionCreate)
	if err != nil {
		log.Errorf("Instance ingress object creation failed - %s", err)
		return nil, err
	}

	objects = append(objects, *obj)

	return objects, nil
}

func (z *ZcashInstanceResourceManager) CreateSnapshotResource(ctx context.Context, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "zcash.CreateSnapshotResource")
	defer func() { logger.LogServiceTime(log) }()

	var req model.SnapshotRequest

	policy := helper.Config.GetPolicyConfig()

	resource := instance.GetResourceByType(model.ResourcePersistentVolumeClaim)

	appRsc := vars.ManagerFactory.GetAppResourceManager(ctx)
	req.Namespace = project.GetNamespace()
	req.VolumeName = resource.Name
	req.SnapshotClass = policy.SnapshotClass

	req.Labels = helper.CreateInstanceLabels(instance)

	return appRsc.CreateSnapshotResource(ctx, instance.Project, instance.Name, &req)
}

func (z *ZcashInstanceResourceManager) CreateSnapshotScheduleResource(ctx context.Context, project *model.Project, instance *model.Instance, scheduleType model.SnapshotScheduleType) ([]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "zcash.CreateSnapshotScheduleResource")
	defer func() { logger.LogServiceTime(log) }()

	policy := helper.Config.GetPolicyConfig()

	var req model.SnapshotScheduleRequest

	resource := instance.GetResourceByType(model.ResourcePersistentVolumeClaim)

	appRsc := vars.ManagerFactory.GetAppResourceManager(ctx)
	req.Namespace = project.GetNamespace()
	req.Schedule = scheduleType
	req.VolumeName = resource.Name
	req.SnapshotClass = policy.SnapshotClass

	req.Labels = helper.CreateInstanceLabels(instance)

	return appRsc.CreateSnapshotScheduleResource(ctx, instance.Project, instance.Name, &req)
}

func (z *ZcashInstanceResourceManager) CreateRotationResource(ctx context.Context, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "zcash.CreateRotationResource")
	defer func() { logger.LogServiceTime(log) }()

	var err error
	fileTemplate, err := helper.Config.GetInstanceTemplate(instance.InstanceType)
	if err != nil {
		return nil, err
	}

	policy := helper.Config.GetPolicyConfig()
	ic, err := helper.Config.GetInstanceConfig(instance.InstanceType)
	if err != nil {
		return nil, err
	}

	username := utils.GenerateRandomString(6, true)
	password := utils.GenerateSecurePassword()

	zcashSpec := model.InstanceSpec{
		Name:               instance.Name,
		ServiceAccountName: policy.ServiceAccount,
		Namespace:          project.GetNamespace(),
		Labels:             helper.CreateInstanceLabels(instance),
		Envoy:              helper.CreateEnvoySpec(ic.GetPort(ENVOY_PORT)),
		Properties: map[string]interface{}{
			USERNAME: username,
			PASSWORD: password,
		},
	}

	var specArr []string

	specArr, err = fileTemplate.ExecuteTemplates([]string{"ENVOY_CONF", "CREDENTIALS"}, zcashSpec)
	if err != nil {
		log.Errorf("Zcash templates failed - %s", err)
		//		return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	return helper.CreateYAMLObjects(specArr /*, instance.Project, instance.Name*/)
}

func (z *ZcashInstanceResourceManager) CreateDeleteResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, resources []model.KubernetesResource) ([]model.KubernetesResource, []unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "zcash.CreateDeleteResource")
	defer func() { logger.LogServiceTime(log) }()

	var err error
	ingressResource, err := z.CreateIngressResource(ctx, projIngress, project, instance, model.EventActionDelete)
	if err != nil {
		return nil, nil, err
	}

	for index := 0; index < len(resources); index++ {
		if resources[index].Type == model.ResourceHTTPProxy {
			resources = append(resources[:index], resources[index+1:]...)
		}
	}

	return resources, []unstructured.Unstructured{*ingressResource}, nil
}
