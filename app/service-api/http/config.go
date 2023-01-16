package http

import (
	"github.com/zbitech/controller/app/service-api/response"
	"github.com/zbitech/controller/internal/helper"
	"github.com/zbitech/controller/pkg/model"
	"github.com/zbitech/controller/pkg/object"
	"net/http"
)

func SetConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	envelope := response.Envelope{}
	if err := response.JSON(w, http.StatusCreated, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}

}

func GetConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//	log := logger.GetLogger(ctx)

	instanceConfig := make(map[model.InstanceType]*object.InstanceConfig)
	instanceTemplate := make(map[model.InstanceType]string)

	pc := helper.Config.GetPolicyConfig()

	//	var err error
	for _, instanceType := range pc.InstanceTypes {
		instanceConfig[instanceType], _ = helper.Config.GetInstanceConfig(instanceType)
		ft, _ := helper.Config.GetInstanceTemplate(instanceType)
		if ft != nil {
			instanceTemplate[instanceType] = ft.Content
		}
	}

	envelope := response.Envelope{"policy": pc}
	if err := response.JSON(w, http.StatusOK, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}

}
