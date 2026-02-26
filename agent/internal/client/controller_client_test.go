package client

import (
	"agent/internal/library/httpclient"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestControllerClient_Register_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/register", r.URL.Path)
		assert.Equal(t, "agent-key", r.Header.Get("X-API-Key"))
		assert.Equal(t, "existing-agent", r.Header.Get("X-Agent-ID"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"agent_id":"agent-1","poll_url":"/config","poll_interval_seconds":30}`))
	}))
	defer srv.Close()

	c := NewControllerClient(srv.URL, "agent-key", httpclient.New(3))
	out, err := c.Register(context.Background(), "existing-agent")

	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, "agent-1", out.AgentID)
	assert.Equal(t, "/config", out.PollURL)
	assert.Equal(t, 30, out.PollIntervalSeconds)
}

func TestControllerClient_Register_StatusError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	c := NewControllerClient(srv.URL, "agent-key", httpclient.New(3))
	out, err := c.Register(context.Background(), "existing-agent")

	assert.Nil(t, out)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "register failed with status 401")
}

func TestControllerClient_GetConfig_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/config", r.URL.Path)
		assert.Equal(t, "agent-key", r.Header.Get("X-API-Key"))
		assert.Equal(t, "agent-1", r.Header.Get("X-Agent-ID"))
		assert.Equal(t, `"1"`, r.Header.Get("If-None-Match"))
		w.Header().Set("ETag", `"2"`)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"version":2,"url":"https://example.com","poll_interval_seconds":45}`))
	}))
	defer srv.Close()

	c := NewControllerClient(srv.URL, "agent-key", httpclient.New(3))
	cfg, etag, status, err := c.GetConfig(context.Background(), "agent-1", `"1"`, "/config")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, `"2"`, etag)
	assert.NotNil(t, cfg)
	assert.Equal(t, 2, cfg.Version)
	assert.Equal(t, "https://example.com", cfg.URL)
	assert.Equal(t, 45, cfg.PollIntervalSeconds)
}

func TestControllerClient_GetConfig_NotModified(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `"9"`)
		w.WriteHeader(http.StatusNotModified)
	}))
	defer srv.Close()

	c := NewControllerClient(srv.URL, "agent-key", httpclient.New(3))
	cfg, etag, status, err := c.GetConfig(context.Background(), "agent-1", `"9"`, "/config")

	assert.NoError(t, err)
	assert.Nil(t, cfg)
	assert.Equal(t, `"9"`, etag)
	assert.Equal(t, http.StatusNotModified, status)
}

func TestControllerClient_GetConfig_StatusError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	c := NewControllerClient(srv.URL, "agent-key", httpclient.New(3))
	cfg, etag, status, err := c.GetConfig(context.Background(), "agent-1", "", "/config")

	assert.Nil(t, cfg)
	assert.Equal(t, "", etag)
	assert.Equal(t, http.StatusBadGateway, status)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "get config failed with status 502")
}

func TestControllerClient_GetConfig_HTTPError(t *testing.T) {
	c := NewControllerClient("://bad", "agent-key", httpclient.New(1))
	cfg, etag, status, err := c.GetConfig(context.Background(), "agent-1", "", "/config")

	assert.Nil(t, cfg)
	assert.Equal(t, "", etag)
	assert.Equal(t, 0, status)
	assert.Error(t, err)
}
