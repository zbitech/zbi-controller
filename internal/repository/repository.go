package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/zbitech/controller/pkg/interfaces"
	"github.com/zbitech/controller/pkg/model"
	"net/http"
)

var client = &http.Client{}

type RepositoryService struct {
}

func NewRepositoryService() interfaces.RepositoryServiceIF {
	return &RepositoryService{}
}

func (repo *RepositoryService) UpdateProjectResource(ctx context.Context, project string, resource *model.KubernetesResource) error {

	jsonReq, _ := json.Marshal(resource)
	req, err := http.NewRequest(http.MethodPut, "https://zbi-cp.zbi:8080/project/resources", bytes.NewBuffer(jsonReq))
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {

	} else {

	}

	defer resp.Body.Close()
	//
	return nil
}

func (repo *RepositoryService) UpdateInstanceResource(ctx context.Context, project, instance string, resource *model.KubernetesResource) error {

	return nil
}
