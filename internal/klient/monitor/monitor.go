package monitor

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zbitech/controller/internal/helper"
	"github.com/zbitech/controller/internal/utils"
	"github.com/zbitech/controller/internal/vars"
	"github.com/zbitech/controller/pkg/interfaces"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"time"
)

type ResourceAction string

const (
	AddResource    ResourceAction = "ADD"
	UpdateResource ResourceAction = "UPDATE"
	DeleteResource ResourceAction = "DELETE"
)

type ResourceStatus struct {
	Resource *model.KubernetesResource
	Id       string
	Level    string
	//	Project  string
	//	Instance string
	Ignore bool
	Reason string
	Ready  bool
}

type KlientInformer struct {
	informer   cache.SharedIndexInformer
	objectType model.ResourceObjectType
}

type KlientMonitor struct {
	typedFactory   informers.SharedInformerFactory
	dynamicFactory dynamicinformer.DynamicSharedInformerFactory
	stopper        chan struct{}
	informers      map[model.ResourceObjectType]KlientInformer
	workQueue      *InformerWorkQueue
	ctx            context.Context
	log            *logrus.Entry
	repoSvc        interfaces.RepositoryServiceIF
	clientSvc      interfaces.KlientIF
}

func (k *KlientInformer) GetIndexer() cache.Indexer {
	return k.informer.GetIndexer()
}

func NewKlientMonitor(ctx context.Context, clientSvc interfaces.KlientIF, repoSvc interfaces.RepositoryServiceIF) interfaces.KlientMonitorIF {

	resync := time.Duration(helper.Config.GetSettings().InformerResync)

	return &KlientMonitor{
		typedFactory:   informers.NewSharedInformerFactoryWithOptions(clientSvc.GetKubernetesClient(), time.Second*resync),
		dynamicFactory: dynamicinformer.NewDynamicSharedInformerFactory(clientSvc.GetDynamicClient(), time.Second*resync),
		stopper:        make(chan struct{}),
		informers:      make(map[model.ResourceObjectType]KlientInformer, 0),
		workQueue:      NewInformerWorkQueue(),
		ctx:            ctx,
		log:            logger.GetLogger(ctx),
		repoSvc:        repoSvc,
		clientSvc:      clientSvc,
	}
}

func (k *KlientMonitor) AddInformer(rType model.ResourceObjectType) {

	var inf cache.SharedIndexInformer

	k.log.WithFields(logrus.Fields{"type": rType}).Infof("Adding informer")
	switch rType {
	case model.ResourceNamespace:
		inf = k.typedFactory.Core().V1().Namespaces().Informer()
	case model.ResourceConfigMap:
		inf = k.typedFactory.Core().V1().ConfigMaps().Informer()
	case model.ResourceSecret:
		inf = k.typedFactory.Core().V1().Secrets().Informer()
	case model.ResourceService:
		inf = k.typedFactory.Core().V1().Services().Informer()
	case model.ResourceDeployment:
		inf = k.typedFactory.Apps().V1().Deployments().Informer()
	case model.ResourcePod:
		inf = k.typedFactory.Core().V1().Pods().Informer()
	case model.ResourcePersistentVolume:
		inf = k.typedFactory.Core().V1().PersistentVolumes().Informer()
	case model.ResourcePersistentVolumeClaim:
		inf = k.typedFactory.Core().V1().PersistentVolumeClaims().Informer()
	case model.ResourceVolumeSnapshot:
		inf = k.dynamicFactory.ForResource(helper.GvrMap[rType]).Informer()
	case model.ResourceSnapshotSchedule:
		inf = k.dynamicFactory.ForResource(helper.GvrMap[rType]).Informer()
	case model.ResourceHTTPProxy:
		inf = k.dynamicFactory.ForResource(helper.GvrMap[rType]).Informer()
	default:
		k.log.WithFields(logrus.Fields{"type": rType}).Warnf("Unable to create informer")
		return
	}

	kinf := KlientInformer{informer: inf, objectType: rType}
	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			k.AddEvent(rType, obj)
		},
		UpdateFunc: func(obj1, obj2 interface{}) {
			k.UpdateEvent(rType, obj1, obj2)
		},
		DeleteFunc: func(obj interface{}) {
			k.DeleteEvent(rType, obj)
		},
	}

	kinf.informer.AddEventHandler(handlers)
	k.informers[rType] = kinf
}

func (k *KlientMonitor) AddEvent(rType model.ResourceObjectType, obj interface{}) {
	k.log.WithFields(logrus.Fields{"type": rType, "object": obj}).Tracef("Processing AddEvent")
	k.processEvent(AddResource, rType, obj, 0)
}

func (k *KlientMonitor) UpdateEvent(rType model.ResourceObjectType, obj1, obj2 interface{}) {
	k.log.WithFields(logrus.Fields{"type": rType, "object1": obj1, "object2": obj2}).Tracef("Processing UpdateEvent")
	if isObjectModified(obj1, obj2) {
		k.processEvent(UpdateResource, rType, obj2, 0)
	}
}

func (k *KlientMonitor) DeleteEvent(rType model.ResourceObjectType, obj interface{}) {
	k.log.WithFields(logrus.Fields{"type": rType, "object": obj}).Tracef("Processing DeleteEvent")
	k.processEvent(DeleteResource, rType, obj, 0)
}

func (k *KlientMonitor) processEvent(action ResourceAction, rType model.ResourceObjectType, obj interface{}, requeueCount int) {

	log := logger.GetServiceLogger(k.ctx, "informer.processEvent")
	//	defer logger.LogServiceTime(log)

	kObj, ok := obj.(runtime.Object)
	if ok {
		var result *ResourceStatus

		switch rType {

		case model.ResourceConfigMap:
			rsc := kObj.(*corev1.ConfigMap)
			if !isZBIObject(rsc.Labels) {
				result = &ResourceStatus{Ignore: true}
			} else {
				result = &ResourceStatus{Resource: helper.CreateCoreResource(k.ctx, rType, rsc, k.clientSvc), Ignore: false, Reason: "",
					Ready: true, Id: rsc.Labels["id"], Level: rsc.Labels["level"]}
			}
		case model.ResourceSecret:
			rsc := kObj.(*corev1.Secret)
			if !isZBIObject(rsc.Labels) {
				result = &ResourceStatus{Ignore: true}
			} else {
				result = &ResourceStatus{Resource: helper.CreateCoreResource(k.ctx, rType, rsc, k.clientSvc), Ignore: false, Reason: "",
					Ready: true, Id: rsc.Labels["id"], Level: rsc.Labels["level"]}
			}

		case model.ResourceDeployment:
			result = DeploymentEvent(k.ctx, action, kObj.(*appsv1.Deployment), k.clientSvc)

		case model.ResourcePod:
			result = PodEvent(k.ctx, action, kObj.(*corev1.Pod), k.clientSvc)

		case model.ResourceService:
			rsc := kObj.(*corev1.Service)
			if !isZBIObject(rsc.Labels) {
				result = &ResourceStatus{Ignore: true}
			} else {
				result = &ResourceStatus{Resource: helper.CreateCoreResource(k.ctx, rType, rsc, k.clientSvc), Ignore: false, Reason: "",
					Ready: true, Id: rsc.Labels["id"], Level: rsc.Labels["level"]}
			}

		case model.ResourcePersistentVolumeClaim:
			result = PersistentVolumeClaimEvent(k.ctx, action, kObj.(*corev1.PersistentVolumeClaim), k.clientSvc)

		case model.ResourceVolumeSnapshot:
			result = VolumeSnapshotEvent(k.ctx, action, kObj.(*unstructured.Unstructured), k.clientSvc)

		case model.ResourceSnapshotSchedule:
			result = SnapshotScheduleEvent(k.ctx, action, kObj.(*unstructured.Unstructured), k.clientSvc)

		case model.ResourceHTTPProxy:
			result = IngressEvent(k.ctx, action, kObj.(*unstructured.Unstructured), k.clientSvc)

		default:
			result = &ResourceStatus{Ignore: true}
			return
		}

		if result != nil && !result.Ignore {

			if action == DeleteResource {
				result.Ready = false
			}

			log.WithFields(logrus.Fields{"result": result, "action": action}).Debugf("processing resource")

			// TODO - just update the resource even if not ready. Need to check if informer will be triggered again
			// as deployment progresses
			// k.UpdateResourceStatus(context.Background(), result.Project, result.Instance, result.Resource)

			if result.Ready {
				k.UpdateResourceStatus(context.Background(), result.Id, result.Level, result.Resource)
			} else {

				var key string
				var err error

				if action == DeleteResource {
					key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
				} else {
					key, err = cache.MetaNamespaceKeyFunc(obj)
				}

				if err != nil {
					log.Errorf("Unable to add item %s to queue - %s", key, err)
				} else {
					fields := logrus.Fields{"key": key, "type": rType, "name": result.Resource.Name, "status": result.Resource.Status, "ready": result.Ready}
					//if !k.workQueue.InQueue(key, rType) {
					//log.Infof("Adding item %s to queue - %s %v", key, result.Resource.Status, result.Ready)
					log.WithFields(fields).Infof("adding item to queue")
					qe := QueueElement{Action: action, Key: key, Type: rType, Object: result, RequeueCount: requeueCount}
					k.workQueue.QueueItem(k.ctx, qe)
					//} else {
					//	log.WithFields(fields).Infof("item already in queue")
					//}
				}
			}
		}
	}
}

func (k *KlientMonitor) GetIndexer(rType model.ResourceObjectType) cache.Indexer {
	inf, ok := k.informers[rType]

	if ok {
		return inf.GetIndexer()
	}

	return nil
}

func (k *KlientMonitor) Stop() {
	k.log.Infof("stopping monitor")

	close(k.stopper)
}

func (k *KlientMonitor) Start() {

	k.log.Infof("starting monitor")

	defer utilruntime.HandleCrash()
	defer k.workQueue.queue.ShutDown()
	//	defer close(k.stopper)

	//	k.Start()
	k.typedFactory.Start(k.stopper)
	k.typedFactory.WaitForCacheSync(k.stopper)

	k.dynamicFactory.Start(k.stopper)
	k.dynamicFactory.WaitForCacheSync(k.stopper)

	k.log.Infof("Starting runWorker ...")
	go wait.Until(k.runWorker, time.Second, k.stopper)
	k.log.Infof("Waiting for informers to complete")

	select {
	case <-k.stopper:
		k.log.Infof("stopping monitor")
		return
	}
}

func (k *KlientMonitor) runWorker() {
	for k.processNextItem(context.Background()) {
	}
}

func (k *KlientMonitor) processNextItem(ctx context.Context) bool {

	log := logger.GetLogger(ctx)
	item, quit := k.workQueue.queue.Get()
	if quit {
		return false
	}

	defer k.workQueue.Done(item)
	qItem := item.(QueueElement)

	log.Infof("action: %s, Key: %s, Kind: %s, Requeue Count: %d", qItem.Action, qItem.Key, qItem.Type, qItem.RequeueCount)

	indexer := k.GetIndexer(qItem.Type)
	if indexer != nil {
		obj, exists, err := indexer.GetByKey(qItem.Key)

		if err != nil {
			log.Errorf("fetching object with key %s from store failed - %s", qItem.Key, err)
			k.workQueue.Forget(qItem)
		} else if !exists {
			log.Warnf("object with key %s no longer exists", qItem.Key)

			k.workQueue.Forget(qItem)
			qItem.Object.Resource.Status = "deleted"
			k.UpdateResourceStatus(ctx, qItem.Object.Id, qItem.Object.Level, qItem.Object.Resource)

		} else {
			k.processEvent(qItem.Action, qItem.Type, obj, qItem.RequeueCount+1)
		}
	} else {
		log.Errorf("Unable to get indexer for %s", qItem.Type)
	}

	return false
}

func (k *KlientMonitor) UpdateResourceStatus(ctx context.Context, id, level string, resource *model.KubernetesResource) {

	if helper.Config.GetSettings().EnableRepository {
		log := logger.GetServiceLogger(ctx, "monitor.UpdateProjectResource")
		repoService := vars.RepositoryFactory.GetRepositoryService()

		if level == "instance" {
			log.Infof("updating instance %s resource %s (%s) status - %s", id, resource.Name, resource.Type, resource.Status)

			if err := repoService.UpdateInstanceResource(ctx, id, resource); err != nil {
				log.Errorf("unable to update instance %s resource %s (%s) - %s", id, resource.Name, resource.Type, err)
			}

		} else if level == "project" {
			log.Infof("updating project %s resource %s (%s) status - %s", id, resource.Name, resource.Type, resource.Status)
			if err := repoService.UpdateProjectResource(ctx, id, resource); err != nil {
				log.Errorf("unable to update project %s resource %s (%s) - %s", id, resource.Name, resource.Type, err)
			}
		}
	}
}

func isObjectModified(old, new interface{}) bool {
	old_m := utils.MarshalObject(old)
	new_m := utils.MarshalObject(new)

	modified := old_m != new_m

	return modified
}

func isZBIObject(labels map[string]string) bool {
	platform, ok := labels["platform"]
	if ok && platform == "zbi" {
		return true
	}

	return false
}

func DeploymentEvent(ctx context.Context, action ResourceAction, obj *appsv1.Deployment, clientSvc interfaces.KlientIF) *ResourceStatus {

	if !isZBIObject(obj.GetLabels()) {
		return &ResourceStatus{Ignore: true}
	}

	resStatus := ResourceStatus{Resource: helper.CreateCoreResource(ctx, model.ResourceDeployment, obj, clientSvc), Id: obj.GetLabels()["id"], Level: obj.GetLabels()["level"]}

	if action != DeleteResource {
		resStatus.Ready = resStatus.Resource.Status == "active"
	}

	return &resStatus
}

func PodEvent(ctx context.Context, action ResourceAction, obj *corev1.Pod, clientSvc interfaces.KlientIF) *ResourceStatus {

	if !isZBIObject(obj.GetLabels()) {
		return &ResourceStatus{Ignore: true}
	}

	resStatus := ResourceStatus{Resource: helper.CreateCoreResource(ctx, model.ResourcePod, obj, clientSvc), Id: obj.GetLabels()["id"], Level: obj.GetLabels()["level"]}

	if action != DeleteResource {
		resStatus.Ready = resStatus.Resource.Status == "running"
	}

	return &resStatus
}

func PersistentVolumeClaimEvent(ctx context.Context, action ResourceAction, obj *corev1.PersistentVolumeClaim, clientSvc interfaces.KlientIF) *ResourceStatus {

	if !isZBIObject(obj.GetLabels()) {
		return &ResourceStatus{Ignore: true}
	}

	resStatus := ResourceStatus{Resource: helper.CreateCoreResource(ctx, model.ResourcePersistentVolumeClaim, obj, clientSvc), Id: obj.GetLabels()["id"], Level: obj.GetLabels()["level"]}

	if action != DeleteResource {
		resStatus.Ready = obj.Status.Phase == corev1.ClaimBound
	}

	return &resStatus
}

func VolumeSnapshotEvent(ctx context.Context, action ResourceAction, obj *unstructured.Unstructured, clientSvc interfaces.KlientIF) *ResourceStatus {

	if !isZBIObject(obj.GetLabels()) {
		return &ResourceStatus{Ignore: true}
	}

	resStatus := ResourceStatus{Resource: helper.CreateCoreResource(ctx, model.ResourceVolumeSnapshot, obj, clientSvc), Id: obj.GetLabels()["id"], Level: obj.GetLabels()["level"]}

	if action != DeleteResource {
		resStatus.Ready = resStatus.Resource.Status == "ready"
	}
	return &resStatus
}

func SnapshotScheduleEvent(ctx context.Context, action ResourceAction, obj *unstructured.Unstructured, clientSvc interfaces.KlientIF) *ResourceStatus {

	if !isZBIObject(obj.GetLabels()) {
		return &ResourceStatus{Ignore: true}
	}

	resStatus := ResourceStatus{Resource: helper.CreateCoreResource(ctx, model.ResourceSnapshotSchedule, obj, clientSvc), Id: obj.GetLabels()["id"], Level: obj.GetLabels()["level"]}

	if action != DeleteResource {
		resStatus.Ready = true
	}

	return &resStatus
}

func IngressEvent(ctx context.Context, action ResourceAction, obj *unstructured.Unstructured, clientSvc interfaces.KlientIF) *ResourceStatus {

	if !isZBIObject(obj.GetLabels()) {
		return &ResourceStatus{Ignore: true}
	}

	resStatus := ResourceStatus{Resource: helper.CreateCoreResource(ctx, model.ResourceHTTPProxy, obj, clientSvc), Id: obj.GetLabels()["id"], Level: obj.GetLabels()["level"]}

	if action != DeleteResource {
		resStatus.Ready = resStatus.Resource.Status == "valid"
	}

	return &resStatus
}
