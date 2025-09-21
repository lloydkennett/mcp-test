package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Auth interface {
	Apply(r *http.Request)
}

type BearerToken string

func (b BearerToken) Apply(r *http.Request) {
	if b != "" {
		r.Header.Set("Authorization", "Bearer "+string(b))
	}
}

type PrivateToken string

func (p PrivateToken) Apply(r *http.Request) {
	if p != "" {
		r.Header.Set("X-PRIVATE-TOKEN", string(p))
	}
}

type Client struct {
	hc      *http.Client
	auth    Auth
	baseURL string
}

func New(baseURL string, auth Auth, hc *http.Client) *Client {
	if hc == nil {
		hc = &http.Client{Timeout: 30 * time.Second}
	}
	baseURL = strings.TrimSpace(baseURL)
	baseURL = strings.TrimRight(baseURL, "/")
	return &Client{hc: hc, auth: auth, baseURL: baseURL}
}

func (c *Client) Do(ctx context.Context, method, path string, body any, out any) error {
	url := c.baseURL + path

	var rdr io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		rdr = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, rdr)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	if c.auth != nil {
		c.auth.Apply(req)
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("http %s %s: status %d: %s", method, path, resp.StatusCode, string(respBody))
	}

	if out != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}
	return nil
}
