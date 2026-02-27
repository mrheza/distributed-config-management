package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"worker/internal/middleware"
	serviceMocks "worker/internal/mocks/service"
	"worker/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	agent := r.Group("/", middleware.APIKeyAuth("worker-secret"))
	agent.POST("/config", h.SetConfig)
	r.GET("/hit", h.Hit)
	r.GET("/state", h.GetState)
	return r
}

func TestSetConfig_Success(t *testing.T) {
	mockSvc := new(serviceMocks.WorkerService)
	h := New(mockSvc)
	r := setupRouter(h)

	body := `{"version":1,"url":"https://example.com","poll_interval_seconds":30}`
	mockSvc.On("ApplyConfig", mock.AnythingOfType("*model.Config")).Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/config", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "worker-secret")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var out model.ConfigUpdateResponse
	err := json.Unmarshal(resp.Body.Bytes(), &out)
	assert.NoError(t, err)
	assert.Equal(t, "config updated", out.Message)
}

func TestSetConfig_ValidationError(t *testing.T) {
	mockSvc := new(serviceMocks.WorkerService)
	h := New(mockSvc)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/config", bytes.NewBufferString(`{"url":"invalid"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "worker-secret")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestSetConfig_ServiceError(t *testing.T) {
	mockSvc := new(serviceMocks.WorkerService)
	h := New(mockSvc)
	r := setupRouter(h)

	body := `{"version":1,"url":"https://example.com"}`
	mockSvc.On("ApplyConfig", mock.AnythingOfType("*model.Config")).Return(errors.New("save failed")).Once()

	req := httptest.NewRequest(http.MethodPost, "/config", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "worker-secret")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestSetConfig_Unauthorized(t *testing.T) {
	mockSvc := new(serviceMocks.WorkerService)
	h := New(mockSvc)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/config", bytes.NewBufferString(`{"version":1,"url":"https://example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	mockSvc.AssertNotCalled(t, "ApplyConfig", mock.Anything)
}

func TestHit_Success(t *testing.T) {
	mockSvc := new(serviceMocks.WorkerService)
	h := New(mockSvc)
	r := setupRouter(h)

	mockSvc.On("Hit", mock.Anything).Return(200, "text/plain", []byte("1.2.3.4"), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/hit", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "1.2.3.4", resp.Body.String())
}

func TestHit_Error(t *testing.T) {
	mockSvc := new(serviceMocks.WorkerService)
	h := New(mockSvc)
	r := setupRouter(h)

	mockSvc.On("Hit", mock.Anything).Return(0, "", nil, errors.New("upstream error")).Once()

	req := httptest.NewRequest(http.MethodGet, "/hit", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestHit_DefaultStatusAndContentType(t *testing.T) {
	mockSvc := new(serviceMocks.WorkerService)
	h := New(mockSvc)
	r := setupRouter(h)

	mockSvc.On("Hit", mock.Anything).Return(0, "", []byte("raw"), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/hit", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "raw", resp.Body.String())
	assert.Contains(t, resp.Header().Get("Content-Type"), "text/plain")
}

func TestGetState_Success(t *testing.T) {
	mockSvc := new(serviceMocks.WorkerService)
	h := New(mockSvc)
	r := setupRouter(h)

	expected := &model.Config{Version: 2, URL: "https://example.com", PollIntervalSeconds: 30}
	mockSvc.On("GetCurrentConfig").Return(expected, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/state", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var out model.Config
	err := json.Unmarshal(resp.Body.Bytes(), &out)
	assert.NoError(t, err)
	assert.Equal(t, *expected, out)
}

func TestGetState_NotFound(t *testing.T) {
	mockSvc := new(serviceMocks.WorkerService)
	h := New(mockSvc)
	r := setupRouter(h)

	mockSvc.On("GetCurrentConfig").Return((*model.Config)(nil), sql.ErrNoRows).Once()

	req := httptest.NewRequest(http.MethodGet, "/state", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestGetState_InternalError(t *testing.T) {
	mockSvc := new(serviceMocks.WorkerService)
	h := New(mockSvc)
	r := setupRouter(h)

	mockSvc.On("GetCurrentConfig").Return((*model.Config)(nil), errors.New("db down")).Once()

	req := httptest.NewRequest(http.MethodGet, "/state", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}
