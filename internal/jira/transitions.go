package jira

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type IssueTransitionsInput struct {
	IssueKey string `json:"issue_key" jsonschema:"Issue key, e.g. PROJ-123"`
}

type IssueTransition struct {
	ID   string `json:"id" jsonschema:"Transition ID"`
	Name string `json:"name" jsonschema:"Transition name"`
	To   struct {
		Name string `json:"name" jsonschema:"Destination status"`
	} `json:"to" jsonschema:"Destination status"`
}

type IssueTransitionsOutput struct {
	Transitions []IssueTransition `json:"transitions" jsonschema:"Available transitions"`
}

func issueTransitionsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "jira_issue_transitions",
		Description: "List available transitions for a Jira issue",
	}
}

func (j *Jira) issueTransitions(ctx context.Context, _ *mcp.CallToolRequest, in IssueTransitionsInput) (*mcp.CallToolResult, IssueTransitionsOutput, error) {
	if in.IssueKey == "" {
		return nil, IssueTransitionsOutput{}, errors.New("issue_key is required")
	}
	path := fmt.Sprintf("/rest/api/2/issue/%s/transitions", in.IssueKey)
	var out IssueTransitionsOutput
	if err := j.client.Do(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, IssueTransitionsOutput{}, err
	}
	b, _ := json.Marshal(out)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, out, nil
}
