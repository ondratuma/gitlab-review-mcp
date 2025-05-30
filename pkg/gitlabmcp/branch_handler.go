// Package gitlabmcp provides the GitLab MCP tool functionality.
package gitlabmcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/ondratuma/gitlab-review-mcp/pkg/git"
)

// GetCurrentBranchHandler handles the getCurrentBranch tool request.
func GetCurrentBranchHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	branch, err := git.GetCurrentBranch()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Current branch is: %s", branch)), nil
}