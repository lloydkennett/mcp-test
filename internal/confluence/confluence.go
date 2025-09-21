package confluence

import (
	"net/http"

	"github.com/lloydkennett/mcp-test/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Confluence struct {
	client *client.Client
}

func New(baseURL, token string, hc *http.Client) *Confluence {
	c := &Confluence{}
	if baseURL == "" || token == "" {
		return c
	}
	c.client = client.New(baseURL, client.BearerToken(token), hc)
	return c
}

func (c *Confluence) Name() string {
	return "confluence"
}

func (c *Confluence) Enabled() bool {
	return c.client != nil
}
func (c *Confluence) Register(server *mcp.Server) {
	mcp.AddTool(server, pagesSearchTool(), c.pagesSearch)
}
