package client

import (
	"agent/internal/library/httpclient"
	"agent/internal/model"
	"context"
	"fmt"
	"net/http"
)

type ControllerClient interface {
	Register(ctx context.Context, existingAgentID string) (*model.RegisterResponse, error)
	GetConfig(ctx context.Context, agentID, etag, pollURL string) (*model.Config, string, int, error)
}

type controllerClient struct {
	baseURL string
	apiKey  string
	http    *httpclient.Client
}

func NewControllerClient(baseURL, apiKey string, httpClient *httpclient.Client) ControllerClient {
	return &controllerClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		http:    httpClient,
	}
}

func (c *controllerClient) Register(ctx context.Context, existingAgentID string) (*model.RegisterResponse, error) {
	var out model.RegisterResponse
	resp, err := c.http.DoJSON(ctx, http.MethodPost, c.baseURL+"/register", map[string]string{
		"X-API-Key":  c.apiKey,
		"X-Agent-ID": existingAgentID,
	}, nil, &out)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("register failed with status %d", resp.StatusCode)
	}
	return &out, nil
}

func (c *controllerClient) GetConfig(ctx context.Context, agentID, etag, pollURL string) (*model.Config, string, int, error) {
	var out model.Config
	resp, err := c.http.DoJSON(ctx, http.MethodGet, c.baseURL+pollURL, map[string]string{
		"X-API-Key":     c.apiKey,
		"X-Agent-ID":    agentID,
		"If-None-Match": etag,
	}, nil, &out)
	if err != nil {
		return nil, "", 0, err
	}

	newETag := resp.Header.Get("ETag")
	switch resp.StatusCode {
	case http.StatusOK:
		return &out, newETag, resp.StatusCode, nil
	case http.StatusNotModified:
		return nil, newETag, resp.StatusCode, nil
	default:
		return nil, newETag, resp.StatusCode, fmt.Errorf("get config failed with status %d", resp.StatusCode)
	}
}
