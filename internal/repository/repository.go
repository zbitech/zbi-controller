package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/zbitech/controller/internal/helper"
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

func (repo *RepositoryService) UpdateProjectResource(ctx context.Context, projectId string, resource *model.KubernetesResource) error {

	var repository = helper.Config.GetSettings().Repository + "/projects/" + projectId + "/resources"
	jsonReq, _ := json.Marshal(resource)
	req, err := http.NewRequest(http.MethodPut, repository, bytes.NewBuffer(jsonReq))
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
		return nil
	} else {
		return errors.New("failed to update project resource")
	}

	resp.Body.Close()
	return nil
}

func (repo *RepositoryService) UpdateInstanceResource(ctx context.Context, instanceId string, resource *model.KubernetesResource) error {

	var repository = helper.Config.GetSettings().Repository + "/instances/" + instanceId + "/resources"
	jsonReq, _ := json.Marshal(resource)
	req, err := http.NewRequest(http.MethodPut, repository, bytes.NewBuffer(jsonReq))
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
		return nil
	} else {
		return errors.New("failed to update instance resource")
	}

	resp.Body.Close()
	return nil
}
