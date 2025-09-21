package gitlab

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CreateMergeRequestInput struct {
	ProjectID    string `json:"project_id" jsonschema:"Project ID (numeric) or URL-encoded path"`
	SourceBranch string `json:"source_branch" jsonschema:"Source branch name"`
	TargetBranch string `json:"target_branch" jsonschema:"Target branch name"`
	Title        string `json:"title" jsonschema:"Merge request title"`
	Description  string `json:"description,omitempty" jsonschema:"Merge request description (optional)"`
	Draft        bool   `json:"draft,omitempty" jsonschema:"Create as draft (optional)"`
	RemoveSource bool   `json:"remove_source_branch,omitempty" jsonschema:"Remove source branch on merge (optional)"`
}

type CreateMergeRequestOutput struct {
	IID          int    `json:"iid" jsonschema:"Merge request IID"`
	ID           int    `json:"id" jsonschema:"Merge request ID"`
	Title        string `json:"title" jsonschema:"Title"`
	State        string `json:"state" jsonschema:"State"`
	WebURL       string `json:"web_url" jsonschema:"Browser URL"`
	SourceBranch string `json:"source_branch" jsonschema:"Source branch"`
	TargetBranch string `json:"target_branch" jsonschema:"Target branch"`
}

func createMergeRequestTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "gitlab_create_merge_request",
		Description: "Create a GitLab merge request",
	}
}

func (g *GitLab) createMergeRequest(ctx context.Context, _ *mcp.CallToolRequest, in CreateMergeRequestInput) (*mcp.CallToolResult, CreateMergeRequestOutput, error) {
	if in.ProjectID == "" {
		return nil, CreateMergeRequestOutput{}, errors.New("project_id is required")
	}
	if in.SourceBranch == "" {
		return nil, CreateMergeRequestOutput{}, errors.New("source_branch is required")
	}
	if in.TargetBranch == "" {
		return nil, CreateMergeRequestOutput{}, errors.New("target_branch is required")
	}
	if in.Title == "" {
		return nil, CreateMergeRequestOutput{}, errors.New("title is required")
	}

	escapedID := url.PathEscape(in.ProjectID)
	path := fmt.Sprintf("/api/v4/projects/%s/merge_requests", escapedID)

	body := map[string]any{
		"source_branch":        in.SourceBranch,
		"target_branch":        in.TargetBranch,
		"title":                in.Title,
		"remove_source_branch": in.RemoveSource,
	}
	if in.Description != "" {
		body["description"] = in.Description
	}
	if in.Draft {
		body["title"] = "Draft: " + in.Title
	}

	var out CreateMergeRequestOutput
	if err := g.client.Do(ctx, http.MethodPost, path, body, &out); err != nil {
		return nil, CreateMergeRequestOutput{}, err
	}

	respBytes, _ := json.Marshal(out)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(respBytes)}}}, out, nil
}
