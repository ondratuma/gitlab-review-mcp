// Package gitlabmcp provides the GitLab MCP tool functionality.
package gitlabmcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Run starts the GitLab MCP tool with the provided configuration.
func Run(gitlabToken, projectID string) error {
	// Create configuration from provided values
	config := Config{
		GitLabToken: gitlabToken,
		ProjectID:   projectID,
	}

	// Create a new MCP server
	s := server.NewMCPServer(
		"GitLab Merge Request MCP",
		"0.0.1",
		server.WithToolCapabilities(false),
	)

	// Create and register tools with logging middleware
	registerTools(s, config)

	// Start the stdio server
	return server.ServeStdio(s)
}

// registerTools registers all tools with the MCP server.
func registerTools(s *server.MCPServer, config Config) {
	// Get current branch tool
	getCurrentBranchTool := mcp.NewTool("get_current_branch",
		mcp.WithDescription("Get the current Git branch"),
	)
	s.AddTool(getCurrentBranchTool, GetCurrentBranchHandler)

	// Get merge request info tool
	getMergeRequestInfoTool := mcp.NewTool("get_merge_request_info",
		mcp.WithDescription("Get general information for merge requests from the currently checked out branch"),
	)

	// Wrap the info handler to include the config
	wrappedInfoHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return GetMergeRequestInfoHandler(ctx, request, config)
	}
	s.AddTool(getMergeRequestInfoTool, wrappedInfoHandler)

	// Get merge request comments tool
	getMergeRequestCommentsTool := mcp.NewTool("get_merge_request_comments",
		mcp.WithDescription("Get comments for merge requests from the currently checked out branch"),
		mcp.WithNumber(
			"mergeRequestIID",
			mcp.Required(),
			mcp.Description("IID Of the Merge Request"),
		),
	)

	// Wrap the comments handler to include the config
	wrappedCommentsHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return GetMergeRequestCommentsHandler(ctx, request, config)
	}
	s.AddTool(getMergeRequestCommentsTool, wrappedCommentsHandler)
}
