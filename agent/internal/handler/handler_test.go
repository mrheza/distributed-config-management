package handler

import (
	service_mocks "agent/internal/mocks/service"
	"agent/internal/model"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/state", h.GetState)
	return r
}

func TestGetState(t *testing.T) {
	mockSvc := service_mocks.NewAgentService(t)
	mockSvc.EXPECT().GetState().Return(&model.State{AgentID: "agent-1"})

	h := New(mockSvc)
	r := setupRouter(h)
	req := httptest.NewRequest(http.MethodGet, "/state", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var out model.State
	err := json.Unmarshal(resp.Body.Bytes(), &out)
	assert.NoError(t, err)
	assert.Equal(t, "agent-1", out.AgentID)
}

func TestGetState_FullStatePayload(t *testing.T) {
	expected := &model.State{
		AgentID:             "agent-99",
		ETag:                "\"7\"",
		ConfigURL:           "https://example.com/config",
		PollURL:             "/config",
		PollIntervalSeconds: 45,
		LastConfigVersion:   7,
	}

	mockSvc := service_mocks.NewAgentService(t)
	mockSvc.EXPECT().GetState().Return(expected)

	h := New(mockSvc)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/state", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var out model.State
	err := json.Unmarshal(resp.Body.Bytes(), &out)
	assert.NoError(t, err)
	assert.Equal(t, *expected, out)
}
