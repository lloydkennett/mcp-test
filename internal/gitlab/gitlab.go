package gitlab

import (
	"net/http"

	"github.com/lloydkennett/mcp-test/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GitLab struct {
	client *client.Client
}

func New(baseURL, token string, hc *http.Client) *GitLab {
	g := &GitLab{}
	if baseURL == "" || token == "" {
		return g
	}
	g.client = client.New(baseURL, client.PrivateToken(token), hc)
	return g
}

func (g *GitLab) Name() string {
	return "gitlab"
}

func (g *GitLab) Enabled() bool {
	return g.client != nil
}
func (g *GitLab) Register(server *mcp.Server) {
	mcp.AddTool(server, mergeRequestsTool(), g.mergeRequests)
	mcp.AddTool(server, approvalRulesTool(), g.approvalRules)
	mcp.AddTool(server, createMergeRequestTool(), g.createMergeRequest)
	mcp.AddTool(server, projectDetailsTool(), g.projectDetails)
}
