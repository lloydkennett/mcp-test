package jira

import (
	"net/http"

	"github.com/lloydkennett/mcp-test/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Jira struct {
	client *client.Client
}

func New(baseURL, token string, hc *http.Client) *Jira {
	j := &Jira{}
	if baseURL == "" || token == "" {
		return j
	}
	j.client = client.New(baseURL, client.BearerToken(token), hc)
	return j
}

func (j *Jira) Name() string {
	return "jira"
}

func (j *Jira) Enabled() bool {
	return j.client != nil
}

func (j *Jira) Register(server *mcp.Server) {
	mcp.AddTool(server, issuesSearchTool(), j.issuesSearch)
	mcp.AddTool(server, transitionIssueTool(), j.transitionIssue)
	mcp.AddTool(server, issueTransitionsTool(), j.issueTransitions)
}
