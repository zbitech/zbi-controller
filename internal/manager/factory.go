package manager

import (
	"context"
	"github.com/zbitech/controller/internal/helper"
	"github.com/zbitech/controller/pkg/interfaces"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
)

type ResourceManagerFactory struct {
	appManager         interfaces.AppResourceManagerIF
	projectManager     interfaces.ProjectResourceManagerIF
	zcashManager       interfaces.InstanceResourceManagerIF
	lightWalletManager interfaces.InstanceResourceManagerIF
}

func NewResourceManagerFactory() interfaces.ResourceManagerFactoryIF {
	return &ResourceManagerFactory{}
}

func (m *ResourceManagerFactory) Init(ctx context.Context) error {

	var log = logger.GetLogger(ctx)
	//	var err error

	//zcashCfg, ok := vars.ResourceConfig.GetInstanceResourceConfig(ztypes.InstanceTypeZCASH)
	//if !ok {
	//	return errors.New("unable to retrieve zcash configuration")
	//}

	log.Infof("initializing resource managers")
	helper.Config.LoadConfig(ctx)
	helper.Config.LoadTemplates(ctx)

	log.Infof("creating zcash resource manager")
	m.zcashManager = NewZcashInstanceResourceManager()

	//lwdCfg, ok := vars.ResourceConfig.GetInstanceResourceConfig(ztypes.InstanceTypeLWD)
	//if !ok {
	//	return errors.New("unable to retrieve lightwalletd configuration")
	//}

	log.Infof("creating lightwalletd resource manager")
	m.lightWalletManager = NewLWDInstanceResourceManager()

	log.Infof("creating project resource manager")
	m.projectManager = NewProjectResourceManager(map[model.InstanceType]interfaces.InstanceResourceManagerIF{
		model.InstanceTypeZCASH: m.zcashManager,
		model.InstanceTypeLWD:   m.lightWalletManager,
	})
	//if err != nil {
	//	logger.Errorf(ctx, "Failed to create project resource manager - %s", err)
	//	return errs.NewApplicationError(errs.ResourceRetrievalError, err)
	//}

	log.Infof("creating App Manager")
	m.appManager = NewAppResourceManager()
	//if err != nil {
	//	logger.Errorf(ctx, "Failed to create ingress resource manager - %s", err)
	//	return errs.NewApplicationError(errs.ResourceRetrievalError, err)
	//}

	return nil
}

func (m *ResourceManagerFactory) GetAppResourceManager(ctx context.Context) interfaces.AppResourceManagerIF {
	return m.appManager
}

func (m *ResourceManagerFactory) GetProjectDataManager(ctx context.Context) interfaces.ProjectResourceManagerIF {
	return m.projectManager
}
