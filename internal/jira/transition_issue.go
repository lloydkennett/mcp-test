package jira

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TransitionIssueInput struct {
	IssueKey     string `json:"issue_key" jsonschema:"Issue key, e.g. PROJ-123"`
	TransitionID string `json:"transition_id" jsonschema:"Transition ID to apply"`
	Comment      string `json:"comment,omitempty" jsonschema:"Optional comment to add during transition"`
}

func (in TransitionIssueInput) Validate() error {
	switch {
	case in.IssueKey == "":
		return errors.New("issue_key is required")
	case in.TransitionID == "":
		return errors.New("transition_id is required")
	}
	return nil
}

type TransitionIssueOutput struct {
	Key          string `json:"key" jsonschema:"Issue key"`
	TransitionID string `json:"transition_id" jsonschema:"Applied transition ID"`
}

func transitionIssueTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "jira_transition_issue",
		Description: "Transition a Jira issue to a new status via transition ID",
	}
}

func (j *Jira) transitionIssue(ctx context.Context, _ *mcp.CallToolRequest, in TransitionIssueInput) (*mcp.CallToolResult, TransitionIssueOutput, error) {
	if err := in.Validate(); err != nil {
		return nil, TransitionIssueOutput{}, err
	}

	path := fmt.Sprintf("/rest/api/2/issue/%s/transitions", in.IssueKey)

	body := map[string]any{
		"transition": map[string]string{"id": in.TransitionID},
	}
	if in.Comment != "" {
		body["update"] = map[string]any{
			"comment": []any{map[string]any{"add": map[string]string{"body": in.Comment}}},
		}
	}

	if err := j.client.Do(ctx, http.MethodPost, path, body, nil); err != nil {
		return nil, TransitionIssueOutput{}, err
	}

	out := TransitionIssueOutput{Key: in.IssueKey, TransitionID: in.TransitionID}
	b, _ := json.Marshal(out)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, out, nil
}
