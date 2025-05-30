// Package gitlabmcp provides the GitLab MCP tool functionality.
package gitlabmcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerTools registers all tools with the MCP server.
func registerTools(s *server.MCPServer, config *Config) {
	// Get current branch tool
	getCurrentBranchTool := mcp.NewTool("get_current_branch",
		mcp.WithDescription("Get the current Git branch"),
	)
	s.AddTool(getCurrentBranchTool, GetCurrentBranchHandler)

	// Get merge request comments tool
	getMergeRequestCommentsTool := mcp.NewTool("get_merge_request_comments",
		mcp.WithDescription("Get comments for merge requests from the current branch"),
	)

	// Wrap the handler to include the config
	wrappedHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return GetMergeRequestCommentsHandler(ctx, request, config)
	}
	s.AddTool(getMergeRequestCommentsTool, wrappedHandler)
}
