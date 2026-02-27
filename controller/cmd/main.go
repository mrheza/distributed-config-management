package main

import (
	_ "controller/docs"
	"controller/internal/config"
	"controller/internal/db"
	"controller/internal/handler"
	"controller/internal/middleware"
	postgresRepo "controller/internal/repository/postgres"
	"controller/internal/service"
	"log"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Controller API
// @version 1.0
// @description API for agent registration and configuration polling
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
func main() {

	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}
	log.Printf(
		"event=controller_config_loaded port=%s gin_mode=%s poll_url=%s database_url_set=%t",
		cfg.Port,
		cfg.GinMode,
		cfg.PollURL,
		cfg.DatabaseURL != "",
	)
	gin.SetMode(cfg.GinMode)

	database, err := db.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	configRepo := postgresRepo.NewConfigRepository(database)
	agentRepo := postgresRepo.NewAgentRepository(database)

	configService := service.NewConfigService(configRepo)
	agentService := service.NewAgentService(agentRepo)

	h := handler.New(cfg, configService, agentService)

	r := gin.New()
	if err := r.SetTrustedProxies(nil); err != nil {
		panic(err)
	}

	r.Use(middleware.RequestLogger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	agent := r.Group("/", middleware.APIKeyAuth(cfg.AgentAPIKey))
	agent.POST("/register", h.RegisterAgent)
	agent.GET("/config", h.GetConfig)

	admin := r.Group("/", middleware.APIKeyAuth(cfg.AdminAPIKey))
	admin.POST("/config", h.CreateConfig)

	addr := ":" + cfg.Port
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
