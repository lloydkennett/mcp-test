package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type IssuesSearchInput struct {
	JQL string `json:"jql" jsonschema:"JQL query string"`
}

type Issue struct {
	ID     string `json:"id" jsonschema:"Issue ID"`
	Key    string `json:"key" jsonschema:"Issue key"`
	Fields struct {
		Summary   string `json:"summary" jsonschema:"Summary"`
		IssueType struct {
			Name string `json:"name" jsonschema:"Issue type"`
		} `json:"issuetype" jsonschema:"Issue type"`
		Status struct {
			Name string `json:"name" jsonschema:"Status"`
		} `json:"status" jsonschema:"Status"`
		Assignee *struct {
			DisplayName string `json:"displayName" jsonschema:"Assignee"`
		} `json:"assignee" jsonschema:"Assignee"`
		Project struct {
			Key string `json:"key" jsonschema:"Project key"`
		} `json:"project" jsonschema:"Project"`
		Priority *struct {
			Name string `json:"name" jsonschema:"Priority"`
		} `json:"priority" jsonschema:"Priority"`
	} `json:"fields" jsonschema:"Fields"`
}

type IssuesSearchOutput struct {
	Total  int     `json:"total" jsonschema:"Total issues"`
	Issues []Issue `json:"issues" jsonschema:"Issues"`
}

func issuesSearchTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "jira_search_issues",
		Description: "Search for Jira issues using JQL (first page)",
	}
}

func (j *Jira) issuesSearch(ctx context.Context, _ *mcp.CallToolRequest, in IssuesSearchInput) (*mcp.CallToolResult, IssuesSearchOutput, error) {
	q := url.Values{}
	if in.JQL != "" {
		q.Set("jql", in.JQL)
	}
	q.Set("maxResults", "50")
	path := "/rest/api/2/search"
	if qs := q.Encode(); qs != "" {
		path += "?" + qs
	}

	var out IssuesSearchOutput
	if err := j.client.Do(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, IssuesSearchOutput{}, err
	}
	b, _ := json.Marshal(out)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, out, nil
}
