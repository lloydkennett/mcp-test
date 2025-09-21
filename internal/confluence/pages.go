package confluence

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type PagesSearchInput struct {
	CQL string `json:"cql" jsonschema:"CQL query string"`
}

type Page struct {
	ID    string `json:"id" jsonschema:"Page ID"`
	Type  string `json:"type" jsonschema:"Type"`
	Title string `json:"title" jsonschema:"Title"`
	Links struct {
		WebUI string `json:"webui" jsonschema:"Web UI path"`
	} `json:"_links" jsonschema:"Links"`
}

type PagesSearchOutput struct {
	Size  int    `json:"size" jsonschema:"Returned result count"`
	Pages []Page `json:"results" jsonschema:"Pages"`
}

func pagesSearchTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "confluence_search_pages",
		Description: "Search for Confluence pages using CQL (first page)",
	}
}

func (c *Confluence) pagesSearch(ctx context.Context, _ *mcp.CallToolRequest, in PagesSearchInput) (*mcp.CallToolResult, PagesSearchOutput, error) {
	q := url.Values{}
	if in.CQL != "" {
		q.Set("cql", in.CQL)
	}
	q.Set("limit", "50")
	path := "/rest/api/content/search"
	if qs := q.Encode(); qs != "" {
		path += "?" + qs
	}

	var out PagesSearchOutput
	if err := c.client.Do(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, PagesSearchOutput{}, err
	}
	b, _ := json.Marshal(out)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, out, nil
}
