package git

import (
	"os"
	"os/exec"
	"testing"
)

// TestGetCurrentBranch tests the GetCurrentBranch function
func TestGetCurrentBranch(t *testing.T) {
	// Save the original exec.Command function and restore it after the test
	originalExecCommand := execCommand
	defer func() { execCommand = originalExecCommand }()

	// Test cases
	tests := []struct {
		name           string
		mockOutput     string
		mockError      error
		expectedBranch string
		expectError    bool
	}{
		{
			name:           "successful branch retrieval",
			mockOutput:     "main\n",
			mockError:      nil,
			expectedBranch: "main",
			expectError:    false,
		},
		{
			name:           "git command error",
			mockOutput:     "",
			mockError:      os.ErrNotExist,
			expectedBranch: "",
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Mock the exec.Command function
			execCommand = func(command string, args ...string) *exec.Cmd {
				cmd := &exec.Cmd{}
				if tc.mockError != nil {
					// Return a command that will fail
					cmd = exec.Command("non-existent-command")
				} else {
					// Create a fake command that outputs the mock output
					cmd = exec.Command("echo", tc.mockOutput)
				}
				return cmd
			}

			// Call the function
			branch, err := GetCurrentBranch()

			// Check the results
			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if branch != tc.expectedBranch {
					t.Errorf("expected branch %q, got %q", tc.expectedBranch, branch)
				}
			}
		})
	}
}