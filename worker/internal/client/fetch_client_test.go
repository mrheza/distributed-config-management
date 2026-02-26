package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFetchClient_DefaultTimeout(t *testing.T) {
	c := NewFetchClient(0)
	fc, ok := c.(*fetchClient)
	if assert.True(t, ok) {
		assert.Equal(t, 10*time.Second, fc.http.Timeout)
	}
}

func TestNewFetchClient_CustomTimeout(t *testing.T) {
	c := NewFetchClient(3)
	fc, ok := c.(*fetchClient)
	if assert.True(t, ok) {
		assert.Equal(t, 3*time.Second, fc.http.Timeout)
	}
}

func TestFetchClient_Get_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	c := NewFetchClient(5)
	status, contentType, body, err := c.Get(context.Background(), srv.URL)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)
	assert.Contains(t, contentType, "text/plain")
	assert.Equal(t, []byte("ok"), body)
}

func TestFetchClient_Get_InvalidURL(t *testing.T) {
	c := NewFetchClient(5)
	_, _, _, err := c.Get(context.Background(), "://bad-url")
	assert.Error(t, err)
}

func TestFetchClient_Get_ContextCanceled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	c := NewFetchClient(5)
	_, _, _, err := c.Get(ctx, srv.URL)
	assert.Error(t, err)
}
