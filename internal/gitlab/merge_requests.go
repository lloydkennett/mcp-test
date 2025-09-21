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

type MergeRequestsInput struct {
	ProjectID string `json:"project_id" jsonschema:"Project ID (numeric) or URL-encoded path"`
	State     string `json:"state,omitempty" jsonschema:"State: opened|closed|locked|merged (default opened)"`
}

type MergeRequestsOutputItem struct {
	IID    int    `json:"iid" jsonschema:"Merge request IID"`
	Title  string `json:"title" jsonschema:"Title"`
	State  string `json:"state" jsonschema:"State"`
	WebURL string `json:"web_url" jsonschema:"Browser URL"`
	Author struct {
		Username string `json:"username" jsonschema:"Author username"`
	} `json:"author" jsonschema:"Author"`
	TargetBranch string `json:"target_branch" jsonschema:"Target branch"`
	SourceBranch string `json:"source_branch" jsonschema:"Source branch"`
}

type MergeRequestsOutput struct {
	Items []MergeRequestsOutputItem `json:"items" jsonschema:"Merge requests"`
}

func mergeRequestsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "gitlab_list_merge_requests",
		Description: "Get first page of GitLab merge requests for a project (state filter optional)",
	}
}

func (g *GitLab) mergeRequests(ctx context.Context, _ *mcp.CallToolRequest, in MergeRequestsInput) (*mcp.CallToolResult, MergeRequestsOutput, error) {
	if in.ProjectID == "" {
		return nil, MergeRequestsOutput{}, errors.New("project_id is required")
	}
	escapedID := url.PathEscape(in.ProjectID)
	state := in.State
	if state == "" {
		state = "opened"
	}
	url := fmt.Sprintf("/api/v4/projects/%s/merge_requests?state=%s", escapedID, state)

	var items []MergeRequestsOutputItem
	if err := g.client.Do(ctx, http.MethodGet, url, nil, &items); err != nil {
		return nil, MergeRequestsOutput{}, err
	}
	out := MergeRequestsOutput{Items: items}
	respBytes, _ := json.Marshal(out)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(respBytes)}}}, out, nil
}
