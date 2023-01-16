package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/zbitech/controller/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpServer struct {
	server
	router *mux.Router
}

func NewHttpServer(port int) *HttpServer {
	return &HttpServer{
		server: server{port: port},
		router: mux.NewRouter(),
	}
}

func (s *HttpServer) Run(ctx context.Context) {
	log := logger.GetLogger(ctx)
	if err := s.serve(ctx); err != nil {
		log.Errorf("An error occurred while running the http server - %s", err)
	}
}

func (s *HttpServer) GetRouter() *mux.Router {
	return s.router
}

func (s *HttpServer) serve(ctx context.Context) error {

	log := logger.GetLogger(ctx)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"OPTIONS", "GET", "POST", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.server.port),
		Handler:      c.Handler(s.router), //s.recoverPanic(s.enableCORS(s.router)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sign := <-quit
		log.Infof("Shutting down server. signal: %s", sign.String())
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	log.Infof("Starting the server on port %d", s.port)
	err := srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	log.Infof("Stopped the server")
	return nil
}
