package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type StyleGuideParams struct {
	Language string `json:"language" jsonschema:"the programming language (go, java, javascript, python)"`
	Topic    string `json:"topic" jsonschema:"the coding standard topic to get information about"`
}

var styleGuideURLs = map[string]string{
	"go":         "https://google.github.io/styleguide/go/",
	"java":       "https://google.github.io/styleguide/javaguide.html",
	"javascript": "https://google.github.io/styleguide/jsguide.html",
	"python":     "https://google.github.io/styleguide/pyguide.html",
}

var styleTopics = map[string][]string{
	"naming":     {"Naming", "Names", "Identifiers"},
	"formatting": {"Formatting", "Style", "Layout", "Source file"},
	"comments":   {"Documentation", "Comments", "Javadoc", "Docstrings"},
	"imports":    {"Imports", "Packages", "Package imports", "Module"},
	"practices":  {"Best Practices", "Guidelines", "Programming Practices", "Conventions"},
}

func extractAllSections(content string) string {
	lower := strings.ToLower(content)
	var foundSections []string

	for topic, searchTerms := range styleTopics {
		for _, term := range searchTerms {
			termLower := strings.ToLower(term)
			start := strings.Index(lower, termLower)
			if start == -1 {
				continue
			}

			end := len(content)
			nextSection := strings.Index(lower[start+len(term):], "<h")
			if nextSection != -1 {
				end = start + len(term) + nextSection
			}

			section := content[start:end]
			section = strings.ReplaceAll(section, "<br>", "\n")
			section = strings.ReplaceAll(section, "</p>", "\n")
			section = strings.ReplaceAll(section, "</li>", "\n")

			foundSections = append(foundSections, "=== "+topic+": "+term+" ===\n"+section+"\n\n")
		}
	}

	return strings.Join(foundSections, "---\n")
}

func GetStyleGuide(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[StyleGuideParams]) (*mcp.CallToolResultFor[any], error) {
	language := strings.ToLower(params.Arguments.Language)

	url, ok := styleGuideURLs[language]
	if !ok {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "Available languages: go, java, javascript, python"}},
		}, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "Error fetching style guide: " + err.Error()}},
		}, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "Error reading response: " + err.Error()}},
		}, nil
	}

	//content := extractAllSections(string(body))
	return &mcp.CallToolResultFor[any]{
		// Content: []mcp.Content{&mcp.TextContent{Text: content}},
		Content: []mcp.Content{&mcp.TextContent{Text: string(body)}},
	}, nil
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{Name: "style-guides", Version: "v1.0.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "style-guide",
		Description: "Get Google's Style Guide information for different programming languages",
	}, GetStyleGuide)

	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		log.Fatal(err)
	}
}
