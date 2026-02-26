package client

import (
	"context"
	"io"
	"net/http"
	"time"
)

type FetchClient interface {
	Get(ctx context.Context, url string) (status int, contentType string, body []byte, err error)
}

type fetchClient struct {
	http *http.Client
}

func NewFetchClient(timeoutSeconds int) FetchClient {
	t := time.Duration(timeoutSeconds) * time.Second
	if t <= 0 {
		t = 10 * time.Second
	}
	return &fetchClient{http: &http.Client{Timeout: t}}
}

func (c *fetchClient) Get(ctx context.Context, url string) (int, string, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, "", nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, "", nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", nil, err
	}

	return resp.StatusCode, resp.Header.Get("Content-Type"), body, nil
}
