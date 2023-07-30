package klient

import (
	"context"
	"github.com/zbitech/controller/internal/helper"
	"github.com/zbitech/controller/internal/klient/client"
	"github.com/zbitech/controller/internal/klient/monitor"
	"github.com/zbitech/controller/internal/klient/zbi"
	"github.com/zbitech/controller/pkg/interfaces"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
)

type KlientFactory struct {
	//klient    interfaces.KlientIF
	client interfaces.ZBIClientIF
	rscMon interfaces.KlientMonitorIF
}

func NewKlientFactory() interfaces.KlientFactoryIF {
	return &KlientFactory{}
}

func (k *KlientFactory) Init(ctx context.Context, repoSvc interfaces.RepositoryServiceIF) error {

	log := logger.GetLogger(ctx)

	log.Infof("creating kubernetes client")
	clientSvc, err := client.NewKlient(ctx)
	if err != nil {
		return err
	}

	log.Infof("creating zbi client")
	k.client = zbi.NewZBIClient(clientSvc)

	if helper.Config.GetSettings().EnableMonitor {
		k.rscMon = monitor.NewKlientMonitor(ctx, clientSvc, repoSvc)
		rtypes := []model.ResourceObjectType{model.ResourceNamespace, model.ResourceConfigMap, model.ResourceSecret, model.ResourceDeployment,
			model.ResourceService, model.ResourcePersistentVolumeClaim, model.ResourceVolumeSnapshot, model.ResourceSnapshotSchedule,
			model.ResourceHTTPProxy}

		for _, rtype := range rtypes {
			log.Infof("Adding %s informer", rtype)
			k.rscMon.AddInformer(rtype)
		}
	}

	return nil
}

func (k *KlientFactory) GetZBIClient() interfaces.ZBIClientIF {
	return k.client
}

func (k *KlientFactory) StartMonitor() {
	if helper.Config.GetSettings().EnableMonitor {
		k.rscMon.Start()
	}
}

func (k *KlientFactory) StopMonitor() {
	if helper.Config.GetSettings().EnableMonitor {
		k.rscMon.Stop()
	}
}
