package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/zbitech/controller/internal/helper"
	"github.com/zbitech/controller/pkg/interfaces"
	"github.com/zbitech/controller/pkg/logger"
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

	log := logger.GetServiceLogger(ctx, "repo.UpdateProjectResource")

	var repository = helper.Config.GetSettings().Repository + "/projects/" + project + "/resources"
	jsonReq, _ := json.Marshal(resource)
	req, err := http.NewRequest(http.MethodPut, repository, bytes.NewBuffer(jsonReq))
	if err != nil {
		return err
	}

	log.WithFields(logrus.Fields{"repo": repository}).Infof("updating project resource")
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

func (repo *RepositoryService) UpdateInstanceResource(ctx context.Context, instance string, resource *model.KubernetesResource) error {

	log := logger.GetServiceLogger(ctx, "repo.UpdateInstanceResource")

	//	var repository = helper.Config.GetSettings().Repository + "/projects/" + project + "/instances/" + instance + "/resources"
	var repository = helper.Config.GetSettings().Repository + "/instances/" + instance + "/resources"

	jsonReq, _ := json.Marshal(resource)
	req, err := http.NewRequest(http.MethodPut, repository, bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Errorf("failed to update instance resource: %s", err)
		return err
	}

	log.WithFields(logrus.Fields{"repo": repository}).Infof("updating instance resource")

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	} else {
		log.WithFields(logrus.Fields{"status": resp.StatusCode, "detail": resp.Body}).Errorf("failed to update instance resource")
		return errors.New("failed to update instance resource")
	}

	resp.Body.Close()
	return nil
}
