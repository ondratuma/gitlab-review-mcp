// Package main provides the entry point for the GitLab MCP tool.
package main

import (
	"fmt"
	"os"

	"github.com/ondratuma/gitlab-review-mcp/pkg/gitlabmcp"
)

func main() {
	// Load configuration from environment variables
	gitlabToken := os.Getenv("GITLAB_TOKEN")
	if gitlabToken == "" {
		fmt.Fprintf(os.Stderr, "Error: GITLAB_TOKEN environment variable is not set\n")
		os.Exit(1)
	}

	projectID := os.Getenv("GITLAB_PROJECT_ID")
	if projectID == "" {
		fmt.Fprintf(os.Stderr, "Error: GITLAB_PROJECT_ID environment variable is not set\n")
		os.Exit(1)
	}

	// This file serves as a simple entry point that delegates to the actual implementation
	// in the pkg/gitlabmcp package.
	if err := gitlabmcp.Run(gitlabToken, projectID); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
