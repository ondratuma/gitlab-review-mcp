// Package gitlabmcp provides the GitLab MCP tool functionality.
package gitlabmcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/ondratuma/gitlab-review-mcp/pkg/gitlab"
)

func GetGroupedCommentThreads(comments []gitlab.MergeRequestNote) (map[string]map[string][]gitlab.MergeRequestNote, error) {
	// Group comments by file path and line code
	groupedComments := make(map[string]map[string][]gitlab.MergeRequestNote)
	for _, comment := range comments {
		if comment.System {
			continue
		}

		// Skip comments without position
		if comment.Position == nil {
			continue
		}

		// Get the file path
		path := "Unknown"
		if comment.Position.NewPath != nil {
			path = *comment.Position.NewPath
		}

		// Get the line code identifier
		lineCode := "general"
		if comment.Position.LineRange != nil &&
			comment.Position.LineRange.Start != nil &&
			comment.Position.LineRange.Start.LineCode != "" {
			lineCode = comment.Position.LineRange.Start.LineCode
		} else if comment.Position.NewLine != nil {
			// Fallback to line number if line_code is not available
			lineCode = fmt.Sprintf("line_%d", *comment.Position.NewLine)
		}

		// Initialize maps if needed
		if _, ok := groupedComments[path]; !ok {
			groupedComments[path] = make(map[string][]gitlab.MergeRequestNote)
		}

		// Add comment to the appropriate group
		groupedComments[path][lineCode] = append(groupedComments[path][lineCode], comment)
	}

	return groupedComments, nil
}

func GetCommentsForMergeRequest(
	mr int,
	client *gitlab.Client,
	config Config,
) ([]string, error) {
	comments, err := client.GetMergeRequestComments(config.ProjectID, mr)
	if err != nil {
		return []string{}, err
	}

	threadsByFiles, err := GetGroupedCommentThreads(comments)
	if err != nil {
		return []string{}, err
	}

	outputThreads := []string{}

	// Format grouped comments
	for path, threadsByLine := range threadsByFiles {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("\nFile: %s\n", path))

		// Process each line code group
		for _, comments := range threadsByLine {
			if len(comments) > 0 {
				comment := comments[0]
				if comment.Position != nil && comment.Position.NewLine != nil {
					sb.WriteString(fmt.Sprintf("Line: %s\n", comment.Position.NewLine))
				}

				sb.WriteString(fmt.Sprintf("Resolved: %t\n", comment.Resolved))
			}

			// Write line information once per line code group

			// Write all comments for this line code
			for _, comment := range comments {
				sb.WriteString(fmt.Sprintf("[%s] Comment by %s\n", comment.CreatedAt, comment.Author.Username))
				sb.WriteString(fmt.Sprintf("%s\n\n", comment.Body))
			}
		}

		outputThreads = append(outputThreads, sb.String())
	}

	return outputThreads, nil
}

// GetMergeRequestCommentsHandler handles the getMergeRequestComments tool request.
func GetMergeRequestCommentsHandler(ctx context.Context, request mcp.CallToolRequest, config Config) (*mcp.CallToolResult, error) {
	mergeRequestId := request.GetInt("mergeRequestIID", -1)
	if mergeRequestId == -1 {
		return mcp.NewToolResultError("Merge request ID is required"), nil
	}

	client := gitlab.NewClient(config.GitLabToken)
	contents := []mcp.Content{}

	// Add each MR's comments as a separate content item
	threadsForMr, err := GetCommentsForMergeRequest(mergeRequestId, client, config)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	contents = append(contents, mcp.TextContent{
		Type: "text",
		Text: fmt.Sprintf("Comments for Merge Request:"),
	})

	for _, thread := range threadsForMr {
		contents = append(contents, mcp.TextContent{
			Type: "text",
			Text: thread,
		})
	}

	return &mcp.CallToolResult{
		Content: contents,
	}, nil
}
