package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type StyleClient struct {
	httpClient *http.Client
}

var (
	goTopics = []string{"guide", "best-practices", "decisions"}

	languageURLs = map[string]string{
		"go":         "https://google.github.io/styleguide/go/index.html",
		"java":       "https://google.github.io/styleguide/javaguide.html",
		"javascript": "https://google.github.io/styleguide/jsguide.html",
		"js":         "https://google.github.io/styleguide/jsguide.html",
		"python":     "https://google.github.io/styleguide/pyguide.html",
		"py":         "https://google.github.io/styleguide/pyguide.html",
		"cpp":        "https://google.github.io/styleguide/cppguide.html",
		"c++":        "https://google.github.io/styleguide/cppguide.html",
		"html":       "https://google.github.io/styleguide/htmlcssguide.html",
		"css":        "https://google.github.io/styleguide/htmlcssguide.html",
		"shell":      "https://google.github.io/styleguide/shellguide.html",
		"bash":       "https://google.github.io/styleguide/shellguide.html",
		"xml":        "https://google.github.io/styleguide/xmlstyle.html",
		"angular":    "https://google.github.io/styleguide/angularjs-google-style.html",
		"angularjs":  "https://google.github.io/styleguide/angularjs-google-style.html",
		"lisp":       "https://google.github.io/styleguide/lispguide.xml",
		"vim":        "https://google.github.io/styleguide/vimscriptguide.xml",
	}
)

func NewStyleClient() *StyleClient {
	return &StyleClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *StyleClient) fetchURL(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

func (c *StyleClient) getGoGuides(ctx context.Context) (string, error) {
	var combined strings.Builder

	for _, topic := range goTopics {
		url := strings.Replace(languageURLs["go"], "index", topic, 1)
		content, err := c.fetchURL(ctx, url)
		if err != nil {
			return "", fmt.Errorf("failed to fetch %s: %w", topic, err)
		}
		combined.WriteString(fmt.Sprintf("=== %s ===\n\n%s\n\n", strings.ToUpper(topic), content))
	}

	return combined.String(), nil
}

func (c *StyleClient) getLanguageGuide(ctx context.Context, language string) (string, error) {
	url, exists := languageURLs[language]
	if !exists {
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	return c.fetchURL(ctx, url)
}

func (c *StyleClient) GetStyleGuide(ctx context.Context, language string) (string, error) {
	language = strings.ToLower(language)

	if language == "go" {
		return c.getGoGuides(ctx)
	}

	return c.getLanguageGuide(ctx, language)
}
