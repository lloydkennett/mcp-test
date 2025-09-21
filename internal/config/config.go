package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	GitLabURL       string
	GitLabToken     string
	SonarQubeURL    string
	SonarQubeToken  string
	JiraURL         string
	JiraToken       string
	ConfluenceURL   string
	ConfluenceToken string
}

func envOrDefault(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}

func Load() (Config, error) {
	cfg := Config{
		GitLabURL:       envOrDefault("GITLAB_URL", "https://gitlab.com"),
		GitLabToken:     strings.TrimSpace(os.Getenv("GITLAB_TOKEN")),
		SonarQubeURL:    envOrDefault("SONARQUBE_URL", "https://sonarqube.com"),
		SonarQubeToken:  strings.TrimSpace(os.Getenv("SONARQUBE_TOKEN")),
		JiraURL:         envOrDefault("JIRA_URL", "https://jira.com"),
		JiraToken:       strings.TrimSpace(os.Getenv("JIRA_TOKEN")),
		ConfluenceURL:   envOrDefault("CONFLUENCE_URL", "https://confluence.com"),
		ConfluenceToken: strings.TrimSpace(os.Getenv("CONFLUENCE_TOKEN")),
	}
	if cfg.GitLabToken == "" && cfg.SonarQubeToken == "" && cfg.JiraToken == "" && cfg.ConfluenceToken == "" {
		return cfg, fmt.Errorf("no service tokens supplied (need at least one of GITLAB_TOKEN, SONARQUBE_TOKEN, JIRA_TOKEN, CONFLUENCE_TOKEN)")
	}
	return cfg, nil
}
