package handler

import (
	"log"
	"net/http"
	"worker/internal/httpresponse"
	"worker/internal/model"
	"worker/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	workerService service.WorkerService
}

func New(ws service.WorkerService) *Handler {
	return &Handler{workerService: ws}
}

// GetState godoc
// @Summary Worker state
// @Description Returns current configuration used by worker
// @Tags worker
// @Produce json
// @Success 200 {object} model.Config
// @Failure 404 {object} httpresponse.ErrorResponse
// @Failure 500 {object} httpresponse.ErrorResponse
// @Router /state [get]
func (h *Handler) GetState(c *gin.Context) {
	cfg, err := h.workerService.GetCurrentConfig()
	if err != nil {
		httpresponse.FromError(c, err)
		return
	}

	c.JSON(http.StatusOK, cfg)
}

// SetConfig godoc
// @Summary Apply worker config
// @Description Called by Agent to update worker configuration
// @Tags worker
// @Accept json
// @Produce json
// @Param X-API-Key header string true "API key"
// @Param request body model.Config true "Worker config"
// @Success 200 {object} model.ConfigUpdateResponse
// @Failure 400 {object} httpresponse.ErrorResponse
// @Failure 401 {object} httpresponse.ErrorResponse
// @Failure 500 {object} httpresponse.ErrorResponse
// @Security ApiKeyAuth
// @Router /config [post]
func (h *Handler) SetConfig(c *gin.Context) {
	var req model.Config
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.ValidationError(c, err, req)
		return
	}

	if err := h.workerService.ApplyConfig(&req); err != nil {
		httpresponse.FromError(c, err)
		return
	}

	log.Printf("event=worker_config_updated version=%d url=%s poll_interval_secs=%d", req.Version, req.URL, req.PollIntervalSeconds)
	c.JSON(http.StatusOK, model.ConfigUpdateResponse{Message: "config updated"})
}

// Hit godoc
// @Summary Execute hit task
// @Description Executes HTTP GET to configured URL and returns raw response body
// @Tags worker
// @Produce plain
// @Success 200 {string} string
// @Failure 404 {object} httpresponse.ErrorResponse
// @Failure 500 {object} httpresponse.ErrorResponse
// @Router /hit [get]
func (h *Handler) Hit(c *gin.Context) {
	status, contentType, body, err := h.workerService.Hit(c.Request.Context())
	if err != nil {
		httpresponse.FromError(c, err)
		return
	}

	if status <= 0 {
		status = http.StatusOK
	}
	if contentType == "" {
		contentType = "text/plain; charset=utf-8"
	}

	c.Data(status, contentType, body)
}
