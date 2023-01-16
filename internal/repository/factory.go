package repository

import (
	"context"
	"github.com/zbitech/controller/pkg/interfaces"
)

type RepositoryFactory struct {
	service interfaces.RepositoryServiceIF
}

func NewRepositoryFactory() interfaces.RepositoryServiceFactoryIF {
	return &RepositoryFactory{}
}

func (r *RepositoryFactory) Init(ctx context.Context) {
	r.service = NewRepositoryService()
}

func (r *RepositoryFactory) GetRepositoryService() interfaces.RepositoryServiceIF {
	return r.service
}
