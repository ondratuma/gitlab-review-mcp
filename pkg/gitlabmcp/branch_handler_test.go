package gitlabmcp

import (
	"context"
	"errors"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/ondratuma/gitlab-review-mcp/pkg/git"
)

// TestGetCurrentBranchHandler tests the GetCurrentBranchHandler function
func TestGetCurrentBranchHandler(t *testing.T) {
	// Save the original GetCurrentBranch function and restore it after the test
	originalGetCurrentBranch := git.GetCurrentBranch
	defer func() { git.GetCurrentBranch = originalGetCurrentBranch }()

	// Test cases
	tests := []struct {
		name           string
		mockBranch     string
		mockError      error
		expectErrorMsg string
	}{
		{
			name:           "successful branch retrieval",
			mockBranch:     "main",
			mockError:      nil,
			expectErrorMsg: "",
		},
		{
			name:           "git command error",
			mockBranch:     "",
			mockError:      errors.New("git command failed"),
			expectErrorMsg: "git command failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Mock the git.GetCurrentBranch function
			git.GetCurrentBranch = func() (string, error) {
				return tc.mockBranch, tc.mockError
			}

			// Call the handler
			result, err := GetCurrentBranchHandler(context.Background(), mcp.CallToolRequest{})

			// Check for unexpected errors
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check the result
			if tc.expectErrorMsg != "" {
				// We expect an error result
				if len(result.Content) != 1 {
					t.Errorf("expected 1 content item, got %d", len(result.Content))
					return
				}

				// Check if the content contains the expected error message
				content := result.Content[0]
				textContent, ok := content.(mcp.TextContent)
				if !ok {
					t.Errorf("expected TextContent, got %T", content)
					return
				}

				if textContent.Text != tc.expectErrorMsg {
					t.Errorf("expected error message %q, got %q", tc.expectErrorMsg, textContent.Text)
				}
			} else {
				// We expect a success result
				if len(result.Content) != 1 {
					t.Errorf("expected 1 content item, got %d", len(result.Content))
					return
				}

				// Check if the content contains the expected branch message
				content := result.Content[0]
				textContent, ok := content.(mcp.TextContent)
				if !ok {
					t.Errorf("expected TextContent, got %T", content)
					return
				}

				expectedText := "Current branch is: " + tc.mockBranch
				if textContent.Text != expectedText {
					t.Errorf("expected text %q, got %q", expectedText, textContent.Text)
				}
			}
		})
	}
}