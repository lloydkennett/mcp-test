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

type ApprovalRulesInput struct {
	ProjectID string `json:"project_id" jsonschema:"Project ID (numeric) or URL-encoded path"`
}

func (in ApprovalRulesInput) Validate() error {
	if in.ProjectID == "" {
		return errors.New("project_id is required")
	}
	return nil
}

type ApprovalRule struct {
	ID                int    `json:"id" jsonschema:"Rule ID"`
	Name              string `json:"name" jsonschema:"Name"`
	ApprovalsRequired int    `json:"approvals_required" jsonschema:"Required approvals"`
	Users             []struct {
		Username string `json:"username" jsonschema:"User"`
	} `json:"users" jsonschema:"Users"`
	Groups []struct {
		FullPath string `json:"full_path" jsonschema:"Group path"`
	} `json:"groups" jsonschema:"Groups"`
}

type ApprovalRulesOutput struct {
	Items []ApprovalRule `json:"items" jsonschema:"Approval rules"`
}

func approvalRulesTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "gitlab_list_approval_rules",
		Description: "Get GitLab project approval rules",
	}
}

func (g *GitLab) approvalRules(ctx context.Context, _ *mcp.CallToolRequest, in ApprovalRulesInput) (*mcp.CallToolResult, ApprovalRulesOutput, error) {
	if err := in.Validate(); err != nil {
		return nil, ApprovalRulesOutput{}, err
	}
	escapedID := url.PathEscape(in.ProjectID)
	url := fmt.Sprintf("/api/v4/projects/%s/approval_rules", escapedID)

	var items []ApprovalRule
	if err := g.client.Do(ctx, http.MethodGet, url, nil, &items); err != nil {
		return nil, ApprovalRulesOutput{}, err
	}
	out := ApprovalRulesOutput{Items: items}
	respBytes, _ := json.Marshal(out)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(respBytes)}}}, out, nil
}
