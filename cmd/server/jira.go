package main

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

type JiraClient struct {
	BaseURL    string
	User       string
	Token      string
	httpClient *http.Client
}

type Transition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TransitionsResponse struct {
	Transitions []Transition `json:"transitions"`
}

func NewJiraClient(baseURL, user, token string) *JiraClient {
	return &JiraClient{
		BaseURL:    baseURL,
		User:       user,
		Token:      token,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *JiraClient) doRequest(ctx context.Context, method, url string, body []byte) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewBuffer(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.User, c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

func (c *JiraClient) getTransitions(ctx context.Context, ticketID string) ([]Transition, error) {
	url := fmt.Sprintf("%s/rest/api/3/issue/%s/transitions", c.BaseURL, ticketID)

	resp, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get transitions: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var transitionsResp TransitionsResponse
	if err := json.Unmarshal(body, &transitionsResp); err != nil {
		return nil, fmt.Errorf("failed to parse transitions: %w", err)
	}

	return transitionsResp.Transitions, nil
}

func (c *JiraClient) findTransitionByStatus(transitions []Transition, status string) (string, error) {
	for _, t := range transitions {
		if strings.EqualFold(t.Name, status) {
			return t.ID, nil
		}
	}
	return "", fmt.Errorf("no transition found for status '%s'", status)
}

func (c *JiraClient) transitionTicket(ctx context.Context, ticketID, transitionID string) error {
	url := fmt.Sprintf("%s/rest/api/3/issue/%s/transitions", c.BaseURL, ticketID)

	body := map[string]any{
		"transition": map[string]string{
			"id": transitionID,
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, "POST", url, jsonBody)
	if err != nil {
		return fmt.Errorf("failed to transition ticket: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func (c *JiraClient) SetTicketStatus(ctx context.Context, ticketID, status string) error {
	transitions, err := c.getTransitions(ctx, ticketID)
	if err != nil {
		return err
	}

	transitionID, err := c.findTransitionByStatus(transitions, status)
	if err != nil {
		return err
	}

	return c.transitionTicket(ctx, ticketID, transitionID)
}

func (c *JiraClient) AddComment(ctx context.Context, ticketID, comment string) error {
	url := fmt.Sprintf("%s/rest/api/3/issue/%s/comment", c.BaseURL, ticketID)

	body := map[string]any{
		"body": map[string]any{
			"type":    "doc",
			"version": 1,
			"content": []map[string]any{
				{
					"type": "paragraph",
					"content": []map[string]any{
						{
							"text": comment,
							"type": "text",
						},
					},
				},
			},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal comment request: %w", err)
	}

	resp, err := c.doRequest(ctx, "POST", url, jsonBody)
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
