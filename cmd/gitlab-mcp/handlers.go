// Package gitlabmcp provides the GitLab MCP tool functionality.
package gitlabmcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/ondratuma/gitlab-review-mcp/pkg/git"
	"github.com/ondratuma/gitlab-review-mcp/pkg/gitlab"
)

// Config holds the application configuration.
type Config struct {
	GitLabToken string
	ProjectID   string
}

// GetCurrentBranchHandler handles the getCurrentBranch tool request.
func GetCurrentBranchHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	branch, err := git.GetCurrentBranch()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Current branch is: %s", branch)), nil
}

// GetMergeRequestCommentsHandler handles the getMergeRequestComments tool request.
func GetMergeRequestCommentsHandler(ctx context.Context, request mcp.CallToolRequest, config *Config) (*mcp.CallToolResult, error) {
	branch, err := git.GetCurrentBranch()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client := gitlab.NewClient(config.GitLabToken)
	sourceBranch := branch

	mrs, err := client.GetMergeRequestsBySourceBranch(config.ProjectID, sourceBranch)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if len(mrs) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No merge requests found for the source branch (%s)", branch)), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d merge request(s):\n", len(mrs)))
	for _, mr := range mrs {
		comments, err := client.GetMergeRequestComments(config.ProjectID, mr.IID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		sb.WriteString(fmt.Sprintf("- !%d: %s\n", mr.IID, mr.Title))
		if len(comments) == 0 {
			sb.WriteString("  (No comments)\n")
			continue
		}

		for _, comment := range comments {
			if comment.System {
				continue
			}

			// Skip resolved comments (currently commented out)
			// if comment.Resolved {
			// 	continue;
			// }

			line := "Unknown"
			path := "Unknown"
			if comment.Position != nil {
				if comment.Position.NewLine != nil {
					line = fmt.Sprintf("%d", *comment.Position.NewLine)
				}

				if comment.Position.NewPath != nil {
					path = *comment.Position.NewPath
				}
			} else {
				continue
			}

			// Format comment information
			sb.WriteString(fmt.Sprintf("  File: %s\n", path))
			sb.WriteString(fmt.Sprintf("  Line: %s\n", line))
			sb.WriteString(fmt.Sprintf("  Comment by %s\n", comment.Author.Username))
			sb.WriteString(fmt.Sprintf("  Resolved %t\n", comment.Resolved)) // Fixed %b to %t for boolean formatting
			sb.WriteString(fmt.Sprintf("  %s\n\n", comment.Body))
		}
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: sb.String(),
			},
		},
	}, nil
}
