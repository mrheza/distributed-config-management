package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	_ "worker/docs"
	"worker/internal/client"
	"worker/internal/config"
	"worker/internal/handler"
	"worker/internal/middleware"
	"worker/internal/repository"
	"worker/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Worker API
// @version 1.0
// @description Worker service for config apply and hit execution
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
func main() {
	cfg := config.Load()
	gin.SetMode(cfg.GinMode)

	repo := repository.NewMemoryConfigRepository()
	fetch := client.NewFetchClient(cfg.RequestTimeoutSeconds)
	workerSvc := service.NewWorkerService(repo, fetch)
	h := handler.New(workerSvc)

	r := gin.New()
	r.Use(middleware.RequestLogger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	agent := r.Group("/", middleware.APIKeyAuth(cfg.AgentAPIKey))
	agent.POST("/config", h.SetConfig)
	r.GET("/hit", h.Hit)
	r.GET("/state", h.GetState)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	addr := ":" + cfg.Port
	srv := &http.Server{Addr: addr, Handler: r}

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
