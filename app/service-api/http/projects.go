package http

import (
	"github.com/pkg/errors"
	"github.com/zbitech/controller/app/service-api/request"
	"github.com/zbitech/controller/app/service-api/response"
	"github.com/zbitech/controller/internal/utils"
	"github.com/zbitech/controller/internal/vars"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
	"github.com/zbitech/controller/pkg/rctx"
	"net/http"
)

func CreateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	var project model.Project
	if err := request.ReadJSON(w, r, &project); err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	user := ctx.Value(rctx.USERID).(string)
	project.Owner = user

	zclient := vars.KlientFactory.GetZBIClient()
	err := zclient.CreateProject(ctx, &project)
	if err != nil {
		log.Errorf("Failed to create project %s - %s", project.Name, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	envelope := response.Envelope{"project": project}
	if err = response.JSON(w, http.StatusCreated, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//	log := logger.GetLogger(ctx)

	projectName := request.GetParameterValue(r, request.PATH_PARAM, "project")
	if len(projectName) == 0 {
		response.BadRequestResponse(w, r, errors.New("project is required"))
		return
	}

	var input struct {
		Project   *model.Project   `json:"project"`
		Instances []model.Instance `json:"instances"`
	}

	if err := request.ReadJSON(w, r, &input); err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	zclient := vars.KlientFactory.GetZBIClient()
	if err := zclient.DeleteProject(ctx, input.Project, input.Instances); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	if err := response.JSON(w, http.StatusNoContent, response.Envelope{}); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func UpdateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	//action := request.GetParameterValue(r, request.GET_PARAM, "action")
	//if len(action) == 0 {
	//	response.BadRequestResponse(w, r, errors.New("action is required"))
	//	return
	//}

	var project model.Project
	if err := request.ReadJSON(w, r, &project); err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	zclient := vars.KlientFactory.GetZBIClient()
	err := zclient.RepairProject(ctx, &project)
	if err != nil {
		log.Errorf("Failed to update project %s - %s", project.Name, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	envelope := response.Envelope{"project": project}
	if err = response.JSON(w, http.StatusOK, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func RepairProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	var project model.Project
	if err := request.ReadJSON(w, r, &project); err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	zclient := vars.KlientFactory.GetZBIClient()
	err := zclient.RepairProject(ctx, &project)
	if err != nil {
		log.Errorf("Failed to update project %s - %s", project.Name, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	envelope := response.Envelope{"project": project}
	if err = response.JSON(w, http.StatusOK, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func GetProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	log.Infof("getting projects")
	zclient := vars.KlientFactory.GetZBIClient()
	projects, err := zclient.GetProjects(ctx)
	if err != nil {
		log.Errorf("failed to retrieve projects")
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	envelope := response.Envelope{"projects": projects}
	if err = response.JSON(w, http.StatusOK, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func GetProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	projectName := request.GetParameterValue(r, request.PATH_PARAM, "project")

	log.Infof("getting projects %s", projectName)
	zclient := vars.KlientFactory.GetZBIClient()
	project, err := zclient.GetProject(ctx, projectName)
	if err != nil {
		log.Errorf("failed to retrieve projects")
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	envelope := response.Envelope{"project": project}
	if err = response.JSON(w, http.StatusOK, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func GetProjectResources(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	projectName := request.GetParameterValue(r, request.PATH_PARAM, "project")

	log.Infof("getting resources for project %s", projectName)
	zclient := vars.KlientFactory.GetZBIClient()
	resources, err := zclient.GetProjectResources(ctx, projectName)
	if err != nil {
		log.Errorf("failed to retrieve resources for project %s", projectName)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	envelope := response.Envelope{"resources": resources}
	if err = response.JSON(w, http.StatusOK, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func GetProjectResource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	project := request.GetParameterValue(r, request.PATH_PARAM, "project")
	resourceName := request.GetParameterValue(r, request.PATH_PARAM, "resource")
	resourceType := request.GetParameterValue(r, request.PATH_PARAM, "type")
	if len(project) == 0 || len(resourceName) == 0 || len(resourceType) == 0 {
		response.BadRequestResponse(w, r, errors.New("project is required"))
		return
	}

	zclient := vars.KlientFactory.GetZBIClient()
	resource, err := zclient.GetProjectResource(ctx, project, resourceName, utils.ResourceObjectType(resourceType))

	if err != nil {
		log.Errorf("Failed to retrieve resource %s (%s) for project %s - %s", resourceName, resourceType, project, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	envelope := response.Envelope{"resource": resource}
	if err = response.JSON(w, http.StatusOK, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}
