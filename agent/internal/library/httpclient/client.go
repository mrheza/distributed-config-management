package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Client struct {
	http *http.Client
}

type ResponseMeta struct {
	StatusCode int
	Header     http.Header
}

func New(timeoutSeconds int) *Client {
	return &Client{
		http: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
	}
}

func (c *Client) DoJSON(ctx context.Context, method, url string, headers map[string]string, body interface{}, out interface{}) (*ResponseMeta, error) {
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewBuffer(raw)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		if v != "" {
			req.Header.Set(k, v)
		}
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if out != nil {
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(out); err != nil && err != io.EOF {
			return &ResponseMeta{
				StatusCode: resp.StatusCode,
				Header:     cloneHeader(resp.Header),
			}, err
		}
	}

	return &ResponseMeta{
		StatusCode: resp.StatusCode,
		Header:     cloneHeader(resp.Header),
	}, nil
}

func cloneHeader(h http.Header) http.Header {
	out := make(http.Header, len(h))
	for k, v := range h {
		cp := make([]string, len(v))
		copy(cp, v)
		out[k] = cp
	}
	return out
}
