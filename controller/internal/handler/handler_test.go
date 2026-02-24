package handler

import (
	"bytes"
	"controller/internal/config"
	serviceMocks "controller/internal/mocks/service"
	"controller/internal/model"
	"controller/internal/repository"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/register", handler.RegisterAgent)
	r.GET("/config", handler.GetConfig)
	r.POST("/config", handler.CreateConfig)

	return r
}

// Success - config exists
func TestRegisterAgent_Success_WithConfig(t *testing.T) {

	mockAgent := new(serviceMocks.AgentService)
	mockConfig := new(serviceMocks.ConfigService)

	cfg := &config.Config{
		PollURL: "/config",
	}

	expectedAgentID := "agent-123"

	mockAgent.
		On("Register").
		Return(expectedAgentID, nil).
		Once()

	mockConfig.
		On("GetLatest").
		Return(&model.Config{
			Version:             1,
			URL:                 "https://example.com",
			PollIntervalSeconds: 60,
		}, nil).
		Once()

	handler := New(cfg, mockConfig, mockAgent)

	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/register", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var body map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &body)

	assert.NoError(t, err)

	assert.Equal(t, expectedAgentID, body["agent_id"])
	assert.Equal(t, "/config", body["poll_url"])
	assert.Equal(t, float64(60), body["poll_interval_seconds"])

	mockAgent.AssertExpectations(t)
	mockConfig.AssertExpectations(t)
}

// Success - configService error → fallback default interval
func TestRegisterAgent_Fallback_DefaultInterval(t *testing.T) {

	mockAgent := new(serviceMocks.AgentService)
	mockConfig := new(serviceMocks.ConfigService)

	cfg := &config.Config{
		PollURL: "/config",
	}

	mockAgent.
		On("Register").
		Return("agent-123", nil).
		Once()

	mockConfig.
		On("GetLatest").
		Return(nil, errors.New("not found")).
		Once()

	handler := New(cfg, mockConfig, mockAgent)

	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/register", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var body map[string]interface{}
	_ = json.Unmarshal(resp.Body.Bytes(), &body)

	assert.Equal(t, float64(30), body["poll_interval_seconds"])
	assert.Equal(t, "/config", body["poll_url"])

	mockAgent.AssertExpectations(t)
	mockConfig.AssertExpectations(t)
}

// Success - config exists but interval = 0 → fallback
func TestRegisterAgent_ConfigIntervalZero_Fallback(t *testing.T) {

	mockAgent := new(serviceMocks.AgentService)
	mockConfig := new(serviceMocks.ConfigService)

	cfg := &config.Config{
		PollURL: "/config",
	}

	mockAgent.
		On("Register").
		Return("agent-123", nil).
		Once()

	mockConfig.
		On("GetLatest").
		Return(&model.Config{
			Version:             1,
			URL:                 "https://example.com",
			PollIntervalSeconds: 0,
		}, nil).
		Once()

	handler := New(cfg, mockConfig, mockAgent)

	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/register", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var body map[string]interface{}
	_ = json.Unmarshal(resp.Body.Bytes(), &body)

	assert.Equal(t, float64(30), body["poll_interval_seconds"])

	mockAgent.AssertExpectations(t)
	mockConfig.AssertExpectations(t)
}

// AgentService error
func TestRegisterAgent_AgentServiceError(t *testing.T) {

	mockAgent := new(serviceMocks.AgentService)
	mockConfig := new(serviceMocks.ConfigService)

	cfg := &config.Config{
		PollURL: "/config",
	}

	mockAgent.
		On("Register").
		Return("", errors.New("register failed")).
		Once()

	handler := New(cfg, mockConfig, mockAgent)

	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/register", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	mockAgent.AssertExpectations(t)
}

//
// GetConfig Tests
//

func TestGetConfig_Success(t *testing.T) {

	mockConfigService := new(serviceMocks.ConfigService)

	expected := &model.Config{
		Version: 1,
		URL:     "https://example.com",
	}

	mockConfigService.
		On("GetLatest").
		Return(expected, nil).
		Once()

	handler := New(nil, mockConfigService, nil)

	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/config", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	assert.Equal(t, `"1"`, resp.Header().Get("ETag"))

	var body model.Config
	err := json.Unmarshal(resp.Body.Bytes(), &body)

	assert.NoError(t, err)
	assert.Equal(t, expected.Version, body.Version)
	assert.Equal(t, expected.URL, body.URL)

	mockConfigService.AssertExpectations(t)
}

func TestGetConfig_NotModified_ETag(t *testing.T) {

	mockConfigService := new(serviceMocks.ConfigService)

	expected := &model.Config{
		Version: 2,
		URL:     "https://example.com",
	}

	mockConfigService.
		On("GetLatest").
		Return(expected, nil).
		Once()

	handler := New(nil, mockConfigService, nil)

	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/config", nil)
	req.Header.Set("If-None-Match", `"2"`)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotModified, resp.Code)

	mockConfigService.AssertExpectations(t)
}

func TestGetConfig_Error(t *testing.T) {

	mockConfigService := new(serviceMocks.ConfigService)

	mockConfigService.
		On("GetLatest").
		Return(nil, errors.New("database error")).
		Once()

	handler := New(nil, mockConfigService, nil)

	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/config", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	mockConfigService.AssertExpectations(t)
}

func TestGetConfig_NotFound(t *testing.T) {

	mockConfigService := new(serviceMocks.ConfigService)

	mockConfigService.
		On("GetLatest").
		Return(nil, repository.ErrConfigNotFound).
		Once()

	handler := New(nil, mockConfigService, nil)

	router := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/config", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	var body map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	assert.NoError(t, err)

	errorObj := body["error"].(map[string]interface{})
	assert.Equal(t, "NOT_FOUND", errorObj["code"])
	assert.Equal(t, "config not found", errorObj["message"])

	mockConfigService.AssertExpectations(t)
}

//
// CreateConfig Tests
//

func TestCreateConfig_Success(t *testing.T) {

	mockConfigService := new(serviceMocks.ConfigService)

	reqBody := `{
		"url": "https://example.com",
		"poll_interval_seconds": 60
	}`

	expectedConfig := &model.Config{
		Version:             3,
		URL:                 "https://example.com",
		PollIntervalSeconds: 60,
	}

	mockConfigService.
		On("Create", "https://example.com", 60).
		Return(nil).
		Once()

	mockConfigService.
		On("GetLatest").
		Return(expectedConfig, nil).
		Once()

	handler := New(nil, mockConfigService, nil)

	router := setupRouter(handler)

	req := httptest.NewRequest(
		http.MethodPost,
		"/config",
		bytes.NewBufferString(reqBody),
	)

	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var body model.Config

	err := json.Unmarshal(resp.Body.Bytes(), &body)

	assert.NoError(t, err)
	assert.Equal(t, expectedConfig.Version, body.Version)
	assert.Equal(t, expectedConfig.URL, body.URL)
	assert.Equal(t, expectedConfig.PollIntervalSeconds, body.PollIntervalSeconds)

	mockConfigService.AssertExpectations(t)
}

func TestCreateConfig_ValidationError_MissingFields(t *testing.T) {

	mockConfigService := new(serviceMocks.ConfigService)

	reqBody := `{}`

	handler := New(nil, mockConfigService, nil)

	router := setupRouter(handler)

	req := httptest.NewRequest(
		http.MethodPost,
		"/config",
		bytes.NewBufferString(reqBody),
	)

	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var body map[string]interface{}

	err := json.Unmarshal(resp.Body.Bytes(), &body)

	assert.NoError(t, err)

	errorObj := body["error"].(map[string]interface{})

	assert.Equal(t, "VALIDATION_ERROR", errorObj["code"])
	assert.Equal(t, "validation failed", errorObj["message"])

	fields := errorObj["fields"].([]interface{})
	assert.Len(t, fields, 2)

	found := map[string]string{}
	for _, f := range fields {
		fieldObj := f.(map[string]interface{})
		found[fieldObj["field"].(string)] = fieldObj["message"].(string)
	}

	assert.Equal(t, "is required", found["url"])
	assert.Equal(t, "is required", found["poll_interval_seconds"])
}

func TestCreateConfig_ValidationError_InvalidURL(t *testing.T) {

	mockConfigService := new(serviceMocks.ConfigService)

	reqBody := `{
		"url": "invalid-url",
		"poll_interval_seconds": 60
	}`

	handler := New(nil, mockConfigService, nil)

	router := setupRouter(handler)

	req := httptest.NewRequest(
		http.MethodPost,
		"/config",
		bytes.NewBufferString(reqBody),
	)

	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var body map[string]interface{}

	err := json.Unmarshal(resp.Body.Bytes(), &body)

	assert.NoError(t, err)

	errorObj := body["error"].(map[string]interface{})

	assert.Equal(t, "VALIDATION_ERROR", errorObj["code"])

	fields := errorObj["fields"].([]interface{})
	assert.Len(t, fields, 1)
	fieldObj := fields[0].(map[string]interface{})
	assert.Equal(t, "url", fieldObj["field"])
	assert.Equal(t, "must be a valid URL", fieldObj["message"])
}

func TestCreateConfig_ValidationError_InvalidType(t *testing.T) {

	mockConfigService := new(serviceMocks.ConfigService)

	reqBody := `{
		"url": 123,
		"poll_interval_seconds": 60
	}`

	handler := New(nil, mockConfigService, nil)
	router := setupRouter(handler)

	req := httptest.NewRequest(
		http.MethodPost,
		"/config",
		bytes.NewBufferString(reqBody),
	)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var body map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	assert.NoError(t, err)

	errorObj := body["error"].(map[string]interface{})
	assert.Equal(t, "VALIDATION_ERROR", errorObj["code"])

	fields := errorObj["fields"].([]interface{})
	assert.Len(t, fields, 1)
	fieldObj := fields[0].(map[string]interface{})
	assert.Equal(t, "url", fieldObj["field"])
	assert.Equal(t, "must be string", fieldObj["message"])
}

func TestCreateConfig_CreateServiceError(t *testing.T) {

	mockConfigService := new(serviceMocks.ConfigService)

	reqBody := `{
		"url": "https://example.com",
		"poll_interval_seconds": 60
	}`

	mockConfigService.
		On("Create", "https://example.com", 60).
		Return(errors.New("db error")).
		Once()

	handler := New(nil, mockConfigService, nil)

	router := setupRouter(handler)

	req := httptest.NewRequest(
		http.MethodPost,
		"/config",
		bytes.NewBufferString(reqBody),
	)

	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	mockConfigService.AssertExpectations(t)
}

func TestCreateConfig_GetLatestError(t *testing.T) {

	mockConfigService := new(serviceMocks.ConfigService)

	reqBody := `{
		"url": "https://example.com",
		"poll_interval_seconds": 60
	}`

	mockConfigService.
		On("Create", "https://example.com", 60).
		Return(nil).
		Once()

	mockConfigService.
		On("GetLatest").
		Return(nil, errors.New("db error")).
		Once()

	handler := New(nil, mockConfigService, nil)

	router := setupRouter(handler)

	req := httptest.NewRequest(
		http.MethodPost,
		"/config",
		bytes.NewBufferString(reqBody),
	)

	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	mockConfigService.AssertExpectations(t)
}
