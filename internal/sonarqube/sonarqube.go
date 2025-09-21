package sonarqube

import (
	"net/http"

	"github.com/lloydkennett/mcp-test/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SonarQube struct {
	client *client.Client
}

func New(baseURL, token string, hc *http.Client) *SonarQube {
	s := &SonarQube{}
	if baseURL == "" || token == "" {
		return s
	}
	s.client = client.New(baseURL, client.BearerToken(token), hc)
	return s
}

func (s *SonarQube) Name() string {
	return "sonarqube"
}

func (s *SonarQube) Enabled() bool {
	return s.client != nil
}
func (s *SonarQube) Register(server *mcp.Server) {
	mcp.AddTool(server, issuesTool(), s.issues)
	mcp.AddTool(server, hotspotsTool(), s.hotspots)
}
