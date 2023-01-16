package klient

import (
	"context"
	klient "github.com/zbitech/controller/fake-zbi/klient/k8s-client"
	zklient "github.com/zbitech/controller/fake-zbi/klient/zbi-klient"
	"github.com/zbitech/controller/pkg/interfaces"
)

type KlientFactory struct {
	zbiKlient interfaces.ZBIClientIF
	klient    interfaces.KlientIF
}

type FakeKlientFactory struct {
	client interfaces.ZBIClientIF
	rscMon interfaces.KlientMonitorIF
}

func NewFakeKlientFactory() interfaces.KlientFactoryIF {
	return &FakeKlientFactory{}
}

func (k *FakeKlientFactory) Init(ctx context.Context, repoSvc interfaces.RepositoryServiceIF) error {

	klient, err := klient.NewFakeKlient(ctx)
	if err != nil {
		return err
	}

	k.client = zklient.NewFakeZBIClient(klient)
	return nil
}

func (k *FakeKlientFactory) GetZBIClient() interfaces.ZBIClientIF {
	return k.client
}

func (k *FakeKlientFactory) StartMonitor() {
	k.rscMon.Start()
}

func (k *FakeKlientFactory) StopMonitor() {
	k.rscMon.Stop()
}
