package httpclient

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type badBody struct{}

func (b badBody) MarshalJSON() ([]byte, error) {
	return nil, errors.New("marshal failed")
}

func TestDoJSON_GetAndDecode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "token", r.Header.Get("X-Token"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := New(5)
	var out map[string]bool
	resp, err := c.DoJSON(context.Background(), http.MethodGet, srv.URL, map[string]string{"X-Token": "token"}, nil, &out)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, true, out["ok"])
}

func TestDoJSON_PostBodyAndContentType(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":1}`))
	}))
	defer srv.Close()

	c := New(5)
	var out map[string]int
	resp, err := c.DoJSON(context.Background(), http.MethodPost, srv.URL, nil, map[string]string{"name": "x"}, &out)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, 1, out["id"])
}

func TestDoJSON_EmptyBodyNoDecodeError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := New(5)
	resp, err := c.DoJSON(context.Background(), http.MethodGet, srv.URL, nil, nil, &map[string]interface{}{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestDoJSON_InvalidURL(t *testing.T) {
	c := New(5)
	resp, err := c.DoJSON(context.Background(), http.MethodGet, "://bad-url", nil, nil, nil)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestDoJSON_MarshalError(t *testing.T) {
	c := New(5)
	resp, err := c.DoJSON(context.Background(), http.MethodPost, "http://example.com", nil, badBody{}, nil)
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "marshal failed")
}
