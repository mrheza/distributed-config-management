package client

import (
	"agent/internal/library/httpclient"
	"agent/internal/model"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkerClient_ApplyConfig_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/config", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "worker-secret", r.Header.Get("X-API-Key"))

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		var payload map[string]interface{}
		err = json.Unmarshal(body, &payload)
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", payload["url"])
		assert.Equal(t, float64(1), payload["version"])
		assert.Equal(t, float64(30), payload["poll_interval_seconds"])

		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewWorkerClient(srv.URL, "worker-secret", httpclient.New(3))
	err := c.ApplyConfig(context.Background(), &model.Config{URL: "https://example.com", Version: 1, PollIntervalSeconds: 30})
	assert.NoError(t, err)
}

func TestWorkerClient_ApplyConfig_StatusError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	c := NewWorkerClient(srv.URL, "worker-secret", httpclient.New(3))
	err := c.ApplyConfig(context.Background(), &model.Config{URL: "https://example.com"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "worker apply failed with status 400")
}

func TestWorkerClient_ApplyConfig_HTTPError(t *testing.T) {
	c := NewWorkerClient("://bad", "worker-secret", httpclient.New(1))
	err := c.ApplyConfig(context.Background(), &model.Config{URL: "https://example.com"})
	assert.Error(t, err)
}
