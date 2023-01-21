package http

import (
	"context"
	"github.com/zbitech/controller/app/service-api/http/middleware"
	"github.com/zbitech/controller/app/service-api/response"
	"github.com/zbitech/controller/app/service-api/server"
	"github.com/zbitech/controller/pkg/logger"
	"net/http"
)

func SetupRoutes(ctx context.Context, server *server.HttpServer) {

	log := logger.GetLogger(ctx)

	router := server.GetRouter()

	log.Infof("initializing middlewares")
	router.Use(middleware.InitRequest, middleware.Logging)

	router.NotFoundHandler = http.HandlerFunc(response.NotFoundResponse)
	router.MethodNotAllowedHandler = http.HandlerFunc(response.MethodNotAllowedResponse)

	cfg := router.PathPrefix("/config").Subrouter()
	cfg.Handle("", middleware.Chain(nil)).Methods(http.MethodPost)
	cfg.Handle("", middleware.Chain(nil)).Methods(http.MethodGet)

	log.Infof("setting project routers")
	project := router.PathPrefix("/projects").Subrouter()
	project.Handle("", middleware.Chain(GetProjects)).Methods(http.MethodGet)
	project.Handle("", middleware.Chain(CreateProject)).Methods(http.MethodPost)
	project.Handle("/{project}", middleware.Chain(GetProject)).Methods(http.MethodGet)
	project.Handle("/{project}", middleware.Chain(DeleteProject)).Methods(http.MethodDelete)
	project.Handle("/{project}", middleware.Chain(UpdateProject)).Methods(http.MethodPut)   // update
	project.Handle("/{project}", middleware.Chain(RepairProject)).Methods(http.MethodPatch) // repair

	project.Handle("/{project}/resources", middleware.Chain(GetProjectResources)).Methods(http.MethodGet)
	project.Handle("/{project}/resources/{resource}/{type}", middleware.Chain(GetProjectResource)).Methods(http.MethodGet)
	//project.Handle("/{project}/resources/{resource}/{type}", middleware.Chain(GetProjectResource)).Methods(http.MethodDelete)

	log.Infof("setting instance routers")
	instances := project.PathPrefix("/{project}/instances").Subrouter()
	instances.Handle("", middleware.Chain(CreateInstance)).Methods(http.MethodPost)
	instances.Handle("", middleware.Chain(GetInstances)).Methods(http.MethodGet)
	instances.Handle("/{instance}", middleware.Chain(GetInstance)).Methods(http.MethodGet)
	instances.Handle("/{instance}", middleware.Chain(UpdateInstance)).Methods(http.MethodPut)
	instances.Handle("/{instance}", middleware.Chain(DeleteInstance)).Methods(http.MethodDelete)

	instances.Handle("/{instance}", middleware.Chain(RepairInstance)).Methods(http.MethodPatch)                                           // repair
	instances.Handle("/{instance}/{action:stop|start|snapshot|schedule|rotate}", middleware.Chain(PatchInstance)).Methods(http.MethodPut) // activate, deactivate, snapshot, backup
	instances.Handle("/{instance}/resources", middleware.Chain(GetInstanceResources)).Methods(http.MethodGet)
	instances.Handle("/{instance}/resources/{resource}/{type}", middleware.Chain(GetInstanceResource)).Methods(http.MethodGet)
	instances.Handle("/{instance}/resources/{resource}/{type}", middleware.Chain(DeleteInstanceResource)).Methods(http.MethodDelete)

}
