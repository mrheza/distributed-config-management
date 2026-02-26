package handler

import (
	"agent/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	agent service.AgentService
}

func New(agent service.AgentService) *Handler {
	return &Handler{agent: agent}
}

// GetState godoc
// @Summary Agent state
// @Tags system
// @Produce json
// @Success 200 {object} model.State
// @Router /state [get]
func (h *Handler) GetState(c *gin.Context) {
	c.JSON(http.StatusOK, h.agent.GetState())
}
