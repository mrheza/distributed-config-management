package handler

import (
	"controller/internal/config"
	"controller/internal/httpresponse"
	"controller/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	config        *config.Config
	configService service.ConfigService
	agentService  service.AgentService
}

type RegisterAgentResponse struct {
	AgentID             string `json:"agent_id"`
	PollURL             string `json:"poll_url"`
	PollIntervalSeconds int    `json:"poll_interval_seconds"`
}

type CreateConfigRequest struct {
	URL                 string `json:"url" binding:"required,url"`
	PollIntervalSeconds int    `json:"poll_interval_seconds" binding:"required,gte=1"`
}

func New(cf *config.Config, cs service.ConfigService, as service.AgentService) *Handler {
	return &Handler{cf, cs, as}
}

// RegisterAgent godoc
// @Summary Register agent
// @Description Register a new agent and return polling info
// @Tags agent
// @Produce json
// @Param X-API-Key header string true "API key"
// @Success 200 {object} RegisterAgentResponse
// @Failure 401 {object} httpresponse.ErrorResponse
// @Failure 500 {object} httpresponse.ErrorResponse
// @Security ApiKeyAuth
// @Router /register [post]
func (h *Handler) RegisterAgent(c *gin.Context) {
	id, err := h.agentService.Register()
	if err != nil {
		httpresponse.FromError(c, err)
		return
	}

	// get global config
	cfg, err := h.configService.GetLatest()

	// fallback if config not set
	pollInterval := 30
	if err == nil && cfg != nil && cfg.PollIntervalSeconds > 0 {
		pollInterval = cfg.PollIntervalSeconds
	}

	c.JSON(http.StatusOK, RegisterAgentResponse{
		AgentID:             id,
		PollURL:             h.config.PollURL,
		PollIntervalSeconds: pollInterval,
	})
}

// GetConfig godoc
// @Summary Get latest config
// @Description Returns latest configuration with ETag support
// @Tags config
// @Produce json
// @Param X-API-Key header string true "API key"
// @Success 200 {object} model.Config
// @Success 304 "Not Modified"
// @Failure 401 {object} httpresponse.ErrorResponse
// @Failure 404 {object} httpresponse.ErrorResponse
// @Failure 500 {object} httpresponse.ErrorResponse
// @Param If-None-Match header string false "ETag value"
// @Security ApiKeyAuth
// @Router /config [get]
func (h *Handler) GetConfig(c *gin.Context) {
	cfg, err := h.configService.GetLatest()
	if err != nil {
		httpresponse.FromError(c, err)
		return
	}

	etag := fmt.Sprintf(`"%d"`, cfg.Version)
	c.Header("ETag", etag)
	if c.GetHeader("If-None-Match") == etag {
		c.Status(http.StatusNotModified)
		return
	}

	c.JSON(http.StatusOK, cfg)
}

// CreateConfig godoc
// @Summary Create config
// @Description Create a new configuration version
// @Tags config
// @Accept json
// @Produce json
// @Param X-API-Key header string true "API key"
// @Param request body CreateConfigRequest true "config payload"
// @Success 201 {object} model.Config
// @Failure 401 {object} httpresponse.ErrorResponse
// @Failure 400 {object} httpresponse.ValidationErrorResponse
// @Failure 500 {object} httpresponse.ErrorResponse
// @Security ApiKeyAuth
// @Router /config [post]
func (h *Handler) CreateConfig(c *gin.Context) {
	var req CreateConfigRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.ValidationError(c, err, req)
		return
	}

	err := h.configService.Create(req.URL, req.PollIntervalSeconds)
	if err != nil {
		httpresponse.FromError(c, err)
		return
	}

	cfg, err := h.configService.GetLatest()
	if err != nil {
		httpresponse.FromError(c, err)
		return
	}

	c.JSON(http.StatusCreated, cfg)
}
