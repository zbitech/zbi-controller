package main

import (
	"context"
	"flag"
	"github.com/zbitech/controller/app/service-api/http"
	"github.com/zbitech/controller/app/service-api/server"
	"github.com/zbitech/controller/internal/klient"
	"github.com/zbitech/controller/internal/manager"
	"github.com/zbitech/controller/internal/repository"
	"github.com/zbitech/controller/internal/vars"
	"github.com/zbitech/controller/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

var (
	Port = flag.Int("port", 8080, "controller listening port --port")
)

func main() {

	flag.Parse()

	ctx := context.Background()

	logger.Init()
	log := logger.GetLogger(ctx)

	log.Info("starting zbi controller ...")

	// TODO - move to factory?
	vars.ManagerFactory = manager.NewResourceManagerFactory()
	vars.KlientFactory = klient.NewKlientFactory()
	vars.RepositoryFactory = repository.NewRepositoryFactory()

	vars.RepositoryFactory.Init(ctx)
	vars.ManagerFactory.Init(ctx)
	vars.KlientFactory.Init(ctx, vars.RepositoryFactory.GetRepositoryService())

	svr := server.NewHttpServer(*Port)
	http.SetupRoutes(ctx, svr)

	log.Info("starting http server")
	go svr.Run(ctx)

	go vars.KlientFactory.StartMonitor()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sign := <-quit

	vars.KlientFactory.StopMonitor()
	log.Infof("Shutting down server. signal: %s", sign.String())

}
