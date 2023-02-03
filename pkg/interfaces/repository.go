package interfaces

import (
	"context"
	"github.com/zbitech/controller/pkg/model"
)

type RepositoryServiceIF interface {
	UpdateProjectResource(ctx context.Context, projectId string, resource *model.KubernetesResource) error
	UpdateInstanceResource(ctx context.Context, instanceId string, resource *model.KubernetesResource) error
}

type RepositoryServiceFactoryIF interface {
	Init(ctx context.Context)
	GetRepositoryService() RepositoryServiceIF
}
