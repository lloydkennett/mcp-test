package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/lloydkennett/mcp-test/internal/config"
	"github.com/lloydkennett/mcp-test/internal/confluence"
	"github.com/lloydkennett/mcp-test/internal/gitlab"
	"github.com/lloydkennett/mcp-test/internal/jira"
	"github.com/lloydkennett/mcp-test/internal/prompts"
	"github.com/lloydkennett/mcp-test/internal/registry"
	"github.com/lloydkennett/mcp-test/internal/sonarqube"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var Version = "v0.0.1"

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := mcp.NewServer(&mcp.Implementation{Name: "devtools", Version: Version}, nil)

	prompts.Register(server)

	var reg registry.Registry

	sharedHC := &http.Client{Timeout: 30 * time.Second}

	reg.Add(gitlab.New(cfg.GitLabURL, cfg.GitLabToken, sharedHC))
	reg.Add(sonarqube.New(cfg.SonarQubeURL, cfg.SonarQubeToken, sharedHC))
	reg.Add(jira.New(cfg.JiraURL, cfg.JiraToken, sharedHC))
	reg.Add(confluence.New(cfg.ConfluenceURL, cfg.ConfluenceToken, sharedHC))

	reg.RegisterAll(server)
	log.Printf("services enabled: %v", reg.EnabledServices())
	log.Printf("services disabled: %v", reg.DisabledServices())

	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
