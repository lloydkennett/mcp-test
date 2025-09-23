package sonarqube

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type IssuesInput struct {
	ProjectKey string `json:"project_key" jsonschema:"Project key"`
	Types      string `json:"types,omitempty" jsonschema:"Types CSV: BUG,VULNERABILITY,CODE_SMELL"`
	Severities string `json:"severities,omitempty" jsonschema:"Severities CSV: BLOCKER,CRITICAL,MAJOR,MINOR,INFO"`
}

func (in IssuesInput) Validate() error {
	if in.ProjectKey == "" {
		return errors.New("project_key is required")
	}
	return nil
}

type Issue struct {
	Key       string `json:"key" jsonschema:"Issue key"`
	Rule      string `json:"rule" jsonschema:"Rule key"`
	Severity  string `json:"severity" jsonschema:"Severity"`
	Type      string `json:"type" jsonschema:"Type"`
	Component string `json:"component" jsonschema:"Component (file path)"`
	Project   string `json:"project" jsonschema:"Project key"`
	Message   string `json:"message" jsonschema:"Message"`
	Line      int    `json:"line,omitempty" jsonschema:"Line"`
}

type IssuesOutput struct {
	Total  int     `json:"total" jsonschema:"Total issues"`
	Issues []Issue `json:"issues" jsonschema:"Issues"`
}

func issuesTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "sonarqube_list_issues",
		Description: "Get first page of SonarQube issues for a project (filters optional)",
	}
}

func (s *SonarQube) issues(ctx context.Context, _ *mcp.CallToolRequest, in IssuesInput) (*mcp.CallToolResult, IssuesOutput, error) {
	if err := in.Validate(); err != nil {
		return nil, IssuesOutput{}, err
	}
	q := url.Values{}
	q.Set("componentKeys", in.ProjectKey)
	if in.Types != "" {
		q.Set("types", in.Types)
	}
	if in.Severities != "" {
		q.Set("severities", in.Severities)
	}
	path := "/api/issues/search"
	if qs := q.Encode(); qs != "" {
		path += "?" + qs
	}

	var out IssuesOutput
	if err := s.client.Do(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, IssuesOutput{}, err
	}
	data, _ := json.Marshal(out)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(data)}}}, out, nil
}
