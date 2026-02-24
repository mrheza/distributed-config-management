package httpresponse

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	URL string `json:"url" validate:"required"`
}

func setupContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestError(t *testing.T) {

	c, w := setupContext()

	Error(c, http.StatusBadRequest, "BAD_REQUEST", "bad request")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "BAD_REQUEST", resp.Error.Code)
	assert.Equal(t, "bad request", resp.Error.Message)
}

func TestInternalServerError(t *testing.T) {

	c, w := setupContext()

	InternalServerError(c, errors.New("db error"))

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "INTERNAL_ERROR", resp.Error.Code)
	assert.Equal(t, "internal server error", resp.Error.Message)
}

func TestUnauthorized(t *testing.T) {

	c, w := setupContext()

	Unauthorized(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "UNAUTHORIZED", resp.Error.Code)
}

func TestNotFound(t *testing.T) {

	c, w := setupContext()

	NotFound(c, "config not found")

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Equal(t, "NOT_FOUND", resp.Error.Code)
	assert.Equal(t, "config not found", resp.Error.Message)
}

func TestValidationError_WithValidatorError(t *testing.T) {

	validate := validator.New()

	req := TestStruct{}

	err := validate.Struct(req)

	c, w := setupContext()

	ValidationError(c, err, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp ValidationErrorResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "VALIDATION_ERROR", resp.Error.Code)
	assert.Len(t, resp.Error.Fields, 1)
	assert.Equal(t, "url", resp.Error.Fields[0].Field)
	assert.Equal(t, "is required", resp.Error.Fields[0].Message)
}

func TestValidationError_WithNonValidatorError(t *testing.T) {

	c, w := setupContext()

	ValidationError(c, errors.New("random error"), TestStruct{})

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp ValidationErrorResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "VALIDATION_ERROR", resp.Error.Code)
	assert.Nil(t, resp.Error.Fields)
}

func TestToSnakeCase(t *testing.T) {

	result := toSnakeCase("PollIntervalSeconds")

	assert.Equal(t, "poll_interval_seconds", result)
}

func TestGetJSONFieldName_WithJSONTag(t *testing.T) {

	req := TestStruct{}

	result := getJSONFieldName(req, "URL")

	assert.Equal(t, "url", result)
}

func TestGetJSONFieldName_Fallback(t *testing.T) {

	req := struct{}{}

	result := getJSONFieldName(req, "UnknownField")

	assert.Equal(t, "unknown_field", result)
}
