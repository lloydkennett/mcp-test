package sonarqube

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type HotspotsInput struct {
	ProjectKey string `json:"project_key" jsonschema:"Project key"`
}

type Hotspot struct {
	Key                      string `json:"key" jsonschema:"Hotspot key"`
	Component                string `json:"component" jsonschema:"Component (file path)"`
	Project                  string `json:"project" jsonschema:"Project key"`
	SecurityCategory         string `json:"securityCategory" jsonschema:"Category"`
	VulnerabilityProbability string `json:"vulnerabilityProbability" jsonschema:"Probability"`
	Status                   string `json:"status" jsonschema:"Status"`
	Message                  string `json:"message" jsonschema:"Message"`
}

type HotspotsOutput struct {
	Paging struct {
		Total int `json:"total" jsonschema:"Total hotspots"`
	} `json:"paging" jsonschema:"Paging"`
	Hotspots []Hotspot `json:"hotspots" jsonschema:"Hotspots"`
}

func hotspotsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "sonarqube_list_hotspots",
		Description: "Get first page of SonarQube security hotspots for a project",
	}
}

func (s *SonarQube) hotspots(ctx context.Context, _ *mcp.CallToolRequest, in HotspotsInput) (*mcp.CallToolResult, HotspotsOutput, error) {
	if in.ProjectKey == "" {
		return nil, HotspotsOutput{}, errors.New("project_key is required")
	}
	q := url.Values{}
	q.Set("projectKey", in.ProjectKey)
	path := "/api/hotspots/search"
	if qs := q.Encode(); qs != "" {
		path += "?" + qs
	}

	var out HotspotsOutput
	if err := s.client.Do(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, HotspotsOutput{}, err
	}
	b, _ := json.Marshal(out)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, out, nil
}
