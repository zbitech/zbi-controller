package repository

import (
	"context"
	"github.com/zbitech/controller/pkg/interfaces"
	"github.com/zbitech/controller/pkg/model"
)

type RepositoryService struct {
}

func NewRepositoryService() interfaces.RepositoryServiceIF {
	return &RepositoryService{}
}

func (repo *RepositoryService) UpdateProjectResource(ctx context.Context, project string, resource *model.KubernetesResource) error {

	//http.Post("", "application/json", nil)
	return nil
}

func (repo *RepositoryService) UpdateInstanceResource(ctx context.Context, project, instance string, resource *model.KubernetesResource) error {

	return nil
}
