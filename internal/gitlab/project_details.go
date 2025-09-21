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

type ProjectDetailsInput struct {
	ProjectID string `json:"project_id" jsonschema:"Project ID (numeric) or URL-encoded path"`
}

type ProjectDetailsOutput struct {
	ID            int    `json:"id" jsonschema:"Project ID"`
	PathWithNS    string `json:"path_with_namespace" jsonschema:"Full path with namespace"`
	DefaultBranch string `json:"default_branch" jsonschema:"Default branch"`
	Visibility    string `json:"visibility" jsonschema:"Visibility"`
	WebURL        string `json:"web_url" jsonschema:"Browser URL"`
}

func projectDetailsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "gitlab_get_project",
		Description: "Get basic GitLab project details (limited fields)",
	}
}

func (g *GitLab) projectDetails(ctx context.Context, _ *mcp.CallToolRequest, in ProjectDetailsInput) (*mcp.CallToolResult, ProjectDetailsOutput, error) {
	if in.ProjectID == "" {
		return nil, ProjectDetailsOutput{}, errors.New("project_id is required")
	}
	escapedID := url.PathEscape(in.ProjectID)
	path := fmt.Sprintf("/api/v4/projects/%s", escapedID)

	out := ProjectDetailsOutput{}
	if err := g.client.Do(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, ProjectDetailsOutput{}, err
	}

	respBytes, _ := json.Marshal(out)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(respBytes)}}}, out, nil
}
