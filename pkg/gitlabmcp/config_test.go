package gitlabmcp

import (
	"testing"
)

// TestNewDefaultConfig tests the NewDefaultConfig function
func TestNewDefaultConfig(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		gitlabToken string
		projectID   string
	}{
		{
			name:        "valid config",
			gitlabToken: "test-token",
			projectID:   "12345",
		},
		{
			name:        "empty token",
			gitlabToken: "",
			projectID:   "12345",
		},
		{
			name:        "empty project ID",
			gitlabToken: "test-token",
			projectID:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function
			config := NewDefaultConfig(tc.gitlabToken, tc.projectID)

			// Check the results
			if config.GitLabToken != tc.gitlabToken {
				t.Errorf("expected GitLabToken %q, got %q", tc.gitlabToken, config.GitLabToken)
			}
			if config.ProjectID != tc.projectID {
				t.Errorf("expected ProjectID %q, got %q", tc.projectID, config.ProjectID)
			}
		})
	}
}