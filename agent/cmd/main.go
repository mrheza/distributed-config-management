package main

import (
	_ "agent/docs"
	"agent/internal/client"
	"agent/internal/config"
	"agent/internal/handler"
	"agent/internal/library/httpclient"
	"agent/internal/middleware"
	"agent/internal/repository"
	"agent/internal/service"
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Agent API
// @version 1.0
// @description Agent service for controller polling and worker sync
// @BasePath /
func main() {
	cfg := config.Load()
	gin.SetMode(cfg.GinMode)

	httpClient := httpclient.New(cfg.RequestTimeoutSeconds)
	controllerClient := client.NewControllerClient(cfg.ControllerBaseURL, cfg.ControllerAPIKey, httpClient)
	workerClient := client.NewWorkerClient(cfg.WorkerBaseURL, cfg.WorkerAPIKey, httpClient)
	stateRepo := repository.NewFileStateRepository(cfg.StatePath)

	agentSvc := service.NewAgentService(
		controllerClient,
		workerClient,
		stateRepo,
		cfg.PollURL,
		cfg.PollIntervalSeconds,
		cfg.MaxBackoffSeconds,
		cfg.BackoffJitterPercent,
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go agentSvc.Run(ctx)

	h := handler.New(agentSvc)
	r := gin.New()
	r.Use(middleware.RequestLogger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/state", h.GetState)

	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
