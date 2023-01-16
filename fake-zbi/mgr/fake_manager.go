package mgr

import (
	"context"
	"github.com/zbitech/controller/fake-zbi/mgr/rsc"
	"github.com/zbitech/controller/pkg/interfaces"
)

type FakeResourceManagerFactory struct {
	ingressManager interfaces.AppResourceManagerIF
	projectManager interfaces.ProjectResourceManagerIF
}

func NewFakeResourceManagerFactory() interfaces.ResourceManagerFactoryIF {
	return &FakeResourceManagerFactory{}
}

func (f *FakeResourceManagerFactory) Init(ctx context.Context) error {
	f.ingressManager = rsc.NewFakeIngressResourceManager()
	f.projectManager = rsc.NewFakeProjectResourceManager()
	return nil
}

func (f *FakeResourceManagerFactory) GetAppResourceManager(ctx context.Context) interfaces.AppResourceManagerIF {
	return f.ingressManager
}

func (f *FakeResourceManagerFactory) GetProjectDataManager(ctx context.Context) interfaces.ProjectResourceManagerIF {
	return f.projectManager
}
