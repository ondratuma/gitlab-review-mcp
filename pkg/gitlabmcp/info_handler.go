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

// GetMergeRequestInfoHandler handles the getMergeRequestInfo tool request.
func GetMergeRequestInfoHandler(ctx context.Context, request mcp.CallToolRequest, config Config) (*mcp.CallToolResult, error) {
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

	// Create a content slice to hold each MR as a separate content item
	contents := make([]mcp.Content, 0, len(mrs))

	// Add a header content with the count of MRs found
	headerContent := mcp.TextContent{
		Type: "text",
		Text: fmt.Sprintf("Found %d merge request(s) for branch %s", len(mrs), branch),
	}
	contents = append(contents, headerContent)

	// Add each MR as a separate content item
	for _, mr := range mrs {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("!%d: %s\n", mr.IID, mr.Title))
		if mr.Description != "" {
			sb.WriteString(fmt.Sprintf("Description: %s\n", mr.Description))
		}
		if mr.Author != nil && mr.Author.UserName != nil {
			sb.WriteString(fmt.Sprintf("Author: %s\n", *mr.Author.UserName))
		}
		sb.WriteString(fmt.Sprintf("Source Branch: %s\n", mr.SourceBranch))
		sb.WriteString(fmt.Sprintf("Target Branch: %s\n", mr.TargetBranch))
		sb.WriteString(fmt.Sprintf("State: %s\n", mr.State))
		sb.WriteString(fmt.Sprintf("URL: %s", mr.WebURL))

		mrContent := mcp.TextContent{
			Type: "text",
			Text: sb.String(),
		}
		contents = append(contents, mrContent)
	}

	return &mcp.CallToolResult{
		Content: contents,
	}, nil
}
