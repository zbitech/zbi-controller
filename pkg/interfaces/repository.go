package interfaces

import (
	"context"
	"github.com/zbitech/controller/pkg/model"
)

type RepositoryServiceIF interface {
	UpdateProjectResource(ctx context.Context, project string, resource *model.KubernetesResource) error
	UpdateInstanceResource(ctx context.Context, instance string, resource *model.KubernetesResource) error
}

type RepositoryServiceFactoryIF interface {
	Init(ctx context.Context)
	GetRepositoryService() RepositoryServiceIF
}
