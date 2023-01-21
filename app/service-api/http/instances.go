package http

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/zbitech/controller/app/service-api/request"
	"github.com/zbitech/controller/app/service-api/response"
	"github.com/zbitech/controller/internal/utils"
	"github.com/zbitech/controller/internal/vars"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
	"net/http"
)

// CreateInstance creates the resources associated with the instance.
// input - the instance to be created and a list of existing instances to be
// peered with the new instance.
// response - the instance and the list of resources created.
func CreateInstance(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := logger.GetLogger(ctx)

	projectName := request.GetParameterValue(r, request.PATH_PARAM, "project")

	var input struct {
		Project  *model.Project         `json:"project"`
		Instance *model.Instance        `json:"instance"`
		Request  *model.ResourceRequest `json:"request"`
	}

	if err := request.ReadJSON(w, r, &input); err != nil {
		log.WithFields(logrus.Fields{"error": err, "project": projectName}).Errorf("failed to read input")
		response.BadRequestResponse(w, r, err)
		return
	}
	log.WithFields(logrus.Fields{"input": input}).Infof("creating new instance")

	//	user := ctx.Value(rctx.USERID).(string)
	//	input.Instance.Owner = user
	input.Instance.Project = projectName
	//	input.Instance.Namespace = projectName

	log.WithFields(logrus.Fields{"instance": input.Instance}).Infof("instance details")

	zclient := vars.KlientFactory.GetZBIClient()
	err := zclient.CreateInstance(ctx, input.Project, input.Instance, input.Request)
	if err != nil {
		//		service.HandleError(ctx, w, r, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	envelope := response.Envelope{"instance": input.Instance}
	if err = response.JSON(w, http.StatusCreated, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func DeleteInstance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	projectName := request.GetParameterValue(r, request.PATH_PARAM, "project")
	instanceName := request.GetParameterValue(r, request.PATH_PARAM, "instance")

	log.WithFields(logrus.Fields{"project": projectName, "instance": instanceName}).Infof("deleting instance")

	var err error
	var input struct {
		Project  *model.Project  `json:"project"`
		Instance *model.Instance `json:"instance"`
	}

	if err = request.ReadJSON(w, r, &input); err != nil {
		log.WithFields(logrus.Fields{"error": err, "project": projectName, "instance": instanceName}).Errorf("failed to read input")
		response.BadRequestResponse(w, r, err)
		return
	}

	//	log.WithFields(logrus.Fields{"input": input}).Infof("delete %d resources from %s.%s", len(input.Instance.Resources), projectName, instanceName)
	zclient := vars.KlientFactory.GetZBIClient()

	err = zclient.DeleteInstance(ctx, input.Project, input.Instance)
	if err != nil {
		log.Errorf("failed to delete instance resources - %s", err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	envelope := response.Envelope{"instance": input.Instance}
	if err = response.JSON(w, http.StatusOK, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func UpdateInstance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	projectName := request.GetParameterValue(r, request.PATH_PARAM, "project")
	instanceName := request.GetParameterValue(r, request.PATH_PARAM, "instance")

	log.WithFields(logrus.Fields{"project": projectName, "instance": instanceName}).Infof("updating instance")

	var input struct {
		Project  *model.Project         `json:"project"`
		Instance *model.Instance        `json:"instance"`
		Request  *model.ResourceRequest `json:"request"`
	}

	if err := request.ReadJSON(w, r, &input); err != nil {
		log.WithFields(logrus.Fields{"error": err, "project": projectName, "instance": instanceName}).Errorf("failed to read input")
		response.BadRequestResponse(w, r, err)
		return
	}

	log.WithFields(logrus.Fields{"project": input.Project, "instance": input.Instance, "request": input.Request}).Infof("updating instance")
	zclient := vars.KlientFactory.GetZBIClient()
	err := zclient.UpdateInstance(ctx, input.Project, input.Instance, input.Request)
	if err != nil {
		//		service.HandleError(ctx, w, r, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	// envelope := response.Envelope{"instance": input.Instance}
	if err = response.JSON(w, http.StatusNoContent, response.Envelope{}); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func RepairInstance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	projectName := request.GetParameterValue(r, request.PATH_PARAM, "project")
	instanceName := request.GetParameterValue(r, request.PATH_PARAM, "instance")

	var input struct {
		Project  *model.Project  `json:"project"`
		Instance *model.Instance `json:"instance"`
	}

	if err := request.ReadJSON(w, r, &input); err != nil {
		log.WithFields(logrus.Fields{"error": err, "project": projectName, "instance": instanceName}).Errorf("failed to read input")
		response.BadRequestResponse(w, r, err)
		return
	}
	log.WithFields(logrus.Fields{"project": input.Project, "instance": input.Instance}).Infof("repairing instance")
	zclient := vars.KlientFactory.GetZBIClient()
	err := zclient.RepairInstance(ctx, input.Project, input.Instance)
	if err != nil {
		//		service.HandleError(ctx, w, r, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	envelope := response.Envelope{"instance": input.Instance}
	if err = response.JSON(w, http.StatusOK, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}

}

func PatchInstance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	projectName := request.GetParameterValue(r, request.PATH_PARAM, "project")
	instanceName := request.GetParameterValue(r, request.PATH_PARAM, "instance")
	action := request.GetParameterValue(r, request.PATH_PARAM, "action")

	log.Infof("patching (%s) instance %s.%s", action, projectName, instanceName)

	var input struct {
		Project  *model.Project             `json:"project"`
		Instance *model.Instance            `json:"instance"`
		Schedule model.SnapshotScheduleType `json:"schedule"`
	}

	if err := request.ReadJSON(w, r, &input); err != nil {
		log.WithFields(logrus.Fields{"error": err, "project": projectName, "instance": instanceName}).Errorf("failed to read input")
		response.BadRequestResponse(w, r, err)
		return
	}

	if input.Project == nil {
		log.WithFields(logrus.Fields{"input": input}).Errorf("project is required")
		response.BadRequestResponse(w, r, errors.New("project is required"))
		return
	}

	if input.Instance == nil {
		log.WithFields(logrus.Fields{"input": input}).Errorf("instance is required")
		response.BadRequestResponse(w, r, errors.New("instance is required"))
		return
	}

	var err error

	zclient := vars.KlientFactory.GetZBIClient()

	switch action {
	case "snapshot":
		err = zclient.CreateSnapshot(ctx, input.Project, input.Instance)
	case "schedule":
		err = zclient.CreateSnapshotSchedule(ctx, input.Project, input.Instance, input.Schedule)
	case "start":
		err = zclient.StartInstance(ctx, input.Project, input.Instance)
	case "stop":
		err = zclient.StopInstance(ctx, input.Project, input.Instance)
	case "rotate":
		err = zclient.RotateInstanceCredentials(ctx, input.Project, input.Instance)
	}

	if err != nil {
		log.Errorf("failed to %s instance %s.%s - %s", action, projectName, instanceName, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	envelope := response.Envelope{"project": input.Project, "instance": input.Instance}
	if err = response.JSON(w, http.StatusOK, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func GetInstances(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	projectName := request.GetParameterValue(r, request.PATH_PARAM, "project")
	log.Infof("getting instances for project %s", projectName)
	zclient := vars.KlientFactory.GetZBIClient()

	project, err := zclient.GetProject(ctx, projectName)
	if err != nil {
		log.Errorf("failed to retrieve project %s", projectName)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	instances, err := zclient.GetAllInstances(ctx, project)
	if err != nil {
		log.Errorf("failed to retrieve intsances for project %s", projectName)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	envelope := response.Envelope{"instances": instances}
	if err = response.JSON(w, http.StatusOK, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func GetInstance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	projectName := request.GetParameterValue(r, request.PATH_PARAM, "project")
	instanceName := request.GetParameterValue(r, request.PATH_PARAM, "instance")

	log.Infof("getting instance %s for project %s", instanceName, projectName)
	zclient := vars.KlientFactory.GetZBIClient()

	project, err := zclient.GetProject(ctx, projectName)
	if err != nil {
		log.Errorf("failed to retrieve project %s", projectName)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	instance, err := zclient.GetInstance(ctx, project, instanceName)
	if err != nil {
		log.Errorf("failed to retrieve intsance %s in project %s", instanceName, projectName)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	envelope := response.Envelope{"instance": instance}
	if err = response.JSON(w, http.StatusOK, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func GetInstanceResources(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	projectName := request.GetParameterValue(r, request.PATH_PARAM, "project")
	instanceName := request.GetParameterValue(r, request.PATH_PARAM, "instance")

	log.Infof("getting resources for instance %s.%s", projectName, instanceName)
	zclient := vars.KlientFactory.GetZBIClient()

	project, err := zclient.GetProject(ctx, projectName)
	if err != nil {
		log.Errorf("failed to retrieve project %s", projectName)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	resources, err := zclient.GetInstanceResources(ctx, project, instanceName)
	if err != nil {
		log.Errorf("failed to retrieve resources for instance %s.%s", projectName, instanceName)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	envelope := response.Envelope{"resources": resources}
	if err = response.JSON(w, http.StatusOK, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func GetInstanceResource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	projectName := request.GetParameterValue(r, request.PATH_PARAM, "project")
	instanceName := request.GetParameterValue(r, request.PATH_PARAM, "instance")
	resourceName := request.GetParameterValue(r, request.PATH_PARAM, "resource")
	resourceType := request.GetParameterValue(r, request.PATH_PARAM, "type")

	log.Infof("getting resource %s (%s) instance %s.%s", resourceName, resourceType, projectName, instanceName)
	zclient := vars.KlientFactory.GetZBIClient()
	project, err := zclient.GetProject(ctx, projectName)
	if err != nil {
		log.Errorf("failed to retrieve project %s", projectName)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	resource, err := zclient.GetInstanceResource(ctx, project, instanceName, resourceName, model.ResourceObjectType(resourceType))
	if err != nil {
		log.Errorf("failed to retrieve resource %s (%s) for instance %s.%s", resourceName, resourceType, projectName, instanceName)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	envelope := response.Envelope{"resource": resource}
	if err = response.JSON(w, http.StatusOK, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func DeleteInstanceResource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	projectName := request.GetParameterValue(r, request.PATH_PARAM, "project")
	instanceName := request.GetParameterValue(r, request.PATH_PARAM, "instance")
	resourceName := request.GetParameterValue(r, request.PATH_PARAM, "resource")
	resourceType := request.GetParameterValue(r, request.PATH_PARAM, "type")

	var input struct {
		Project  *model.Project
		Instance *model.Instance
	}

	if err := request.ReadJSON(w, r, &input); err != nil {
		log.WithFields(logrus.Fields{"error": err, "project": projectName, "instance": instanceName}).Errorf("failed to read input")
		response.BadRequestResponse(w, r, err)
		return
	}
	resource := input.Instance.GetResource(resourceName, utils.ResourceObjectType(resourceType))

	log.Infof("deleting resource %s (%s) instance %s.%s", resource.Name, resource.Type, projectName, instanceName)
	zclient := vars.KlientFactory.GetZBIClient()
	err := zclient.DeleteInstanceResource(ctx, input.Project, input.Instance, resource.Name, resource.Type)
	if err != nil {
		log.Errorf("failed to delete resource %s (%s) for instance %s.%s", resource.Name, resource.Type, projectName, instanceName)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	envelope := response.Envelope{}
	if err = response.JSON(w, http.StatusNoContent, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}
