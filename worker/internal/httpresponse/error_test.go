package httpresponse

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestError(t *testing.T) {
	c, w := newTestContext()
	Error(c, http.StatusBadRequest, "BAD_REQUEST", "bad request")
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "BAD_REQUEST")
}

func TestFromError_NotFound(t *testing.T) {
	c, w := newTestContext()
	FromError(c, sql.ErrNoRows)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "NOT_FOUND")
}

func TestFromError_Internal(t *testing.T) {
	c, w := newTestContext()
	FromError(c, errors.New("boom"))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
}

func TestBadRequest_DefaultMessage(t *testing.T) {
	c, w := newTestContext()
	BadRequest(c, "")
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "bad request")
}

func TestNotFound_DefaultMessage(t *testing.T) {
	c, w := newTestContext()
	NotFound(c, "")
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "resource not found")
}

func TestUnauthorized(t *testing.T) {
	c, w := newTestContext()
	Unauthorized(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "UNAUTHORIZED")
}
