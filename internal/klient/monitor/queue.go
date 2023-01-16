package monitor

import (
	"context"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type InformerWorkQueue struct {
	queue   workqueue.RateLimitingInterface
	pending map[string]string
}

type QueueElement struct {
	Action       ResourceAction
	Type         model.ResourceObjectType
	Key          string
	Object       *ResourceStatus
	RequeueCount int
}

func NewInformerWorkQueue() *InformerWorkQueue {
	return &InformerWorkQueue{
		queue:   workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		pending: make(map[string]string),
	}
}

func (inf *InformerWorkQueue) QueueItem(ctx context.Context, qe QueueElement) {
	go func() {

		log := logger.GetLogger(ctx)
		requeueDelay := time.Duration(60) * time.Second
		log.Debugf("Waiting for %s to requeue %s %s", requeueDelay.String(), qe.Type, qe.Object.Resource.Name)
		time.Sleep(requeueDelay)
		log.Infof("Adding %s back to the work queue %d times %s", qe.Key, qe.RequeueCount, qe.Object.Resource.Status)
		if qe.RequeueCount == 0 {
			inf.queue.Add(qe)
		} else {
			inf.queue.AddRateLimited(qe)
		}
		//		inf.pending[qe.Key+string(qe.Object.Resource.Type)] = "pending"
	}()
}

//func (inf *InformerWorkQueue) InQueue(key string, rType model.ResourceObjectType) bool {
//	_, in := inf.pending[key+string(rType)]
//	return in
//}

func (inf *InformerWorkQueue) Forget(qe QueueElement) {
	//delete(inf.pending, qe.Key+string(qe.Object.Resource.Type))
	inf.queue.Forget(qe)
}

func (inf *InformerWorkQueue) Get() (interface{}, bool) {
	return inf.queue.Get()
}

func (inf *InformerWorkQueue) Done(item interface{}) {
	inf.queue.Done(item)
}
