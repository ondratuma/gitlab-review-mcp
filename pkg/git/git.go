// Package git provides utilities for interacting with Git repositories.
package git

import (
	"bytes"
	"os/exec"
	"strings"
)

// execCommand is a variable that holds the exec.Command function.
// It can be replaced in tests to mock the command execution.
var execCommand = exec.Command

// GetCurrentBranchFunc is the function type for GetCurrentBranch
type GetCurrentBranchFunc func() (string, error)

// getCurrentBranchImpl is the actual implementation of GetCurrentBranch
func getCurrentBranchImpl() (string, error) {
	cmd := execCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

// GetCurrentBranch is a variable that holds the getCurrentBranchImpl function.
// It can be replaced in tests to mock the function.
var GetCurrentBranch GetCurrentBranchFunc = getCurrentBranchImpl
