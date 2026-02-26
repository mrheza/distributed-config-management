package client

import (
	"agent/internal/library/httpclient"
	"agent/internal/model"
	"context"
	"fmt"
	"net/http"
)

type WorkerClient interface {
	ApplyConfig(ctx context.Context, cfg *model.Config) error
}

type workerClient struct {
	baseURL string
	apiKey  string
	http    *httpclient.Client
}

func NewWorkerClient(baseURL, apiKey string, httpClient *httpclient.Client) WorkerClient {
	return &workerClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		http:    httpClient,
	}
}

func (w *workerClient) ApplyConfig(ctx context.Context, cfg *model.Config) error {
	resp, err := w.http.DoJSON(ctx, http.MethodPost, w.baseURL+"/config", map[string]string{
		"X-API-Key": w.apiKey,
	}, cfg, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("worker apply failed with status %d", resp.StatusCode)
	}
	return nil
}
