package gitlab

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

// MockHTTPClient is a mock implementation of the http.Client
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do implements the http.Client interface
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// TestGetMergeRequestsBySourceBranch tests the GetMergeRequestsBySourceBranch method
func TestGetMergeRequestsBySourceBranch(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		projectID      string
		sourceBranch   string
		responseStatus int
		responseBody   string
		expectedMRs    []MergeRequest
		expectError    bool
	}{
		{
			name:           "successful retrieval",
			projectID:      "12345",
			sourceBranch:   "feature-branch",
			responseStatus: http.StatusOK,
			responseBody:   `[{"id":1,"iid":1,"title":"Test MR","description":"Test description","web_url":"https://gitlab.com/test/project/-/merge_requests/1","source_branch":"feature-branch","target_branch":"main","state":"opened","author":{"username":"testuser"}}]`,
			expectedMRs: []MergeRequest{
				{
					ID:           1,
					IID:          1,
					Title:        "Test MR",
					Description:  "Test description",
					WebURL:       "https://gitlab.com/test/project/-/merge_requests/1",
					SourceBranch: "feature-branch",
					TargetBranch: "main",
					State:        "opened",
					Author: &struct {
						UserName *string "json:\"username\""
					}{
						UserName: strPtr("testuser"),
					},
				},
			},
			expectError: false,
		},
		{
			name:           "API error",
			projectID:      "12345",
			sourceBranch:   "feature-branch",
			responseStatus: http.StatusInternalServerError,
			responseBody:   `{"error":"Internal server error"}`,
			expectedMRs:    nil,
			expectError:    true,
		},
		{
			name:           "no merge requests",
			projectID:      "12345",
			sourceBranch:   "feature-branch",
			responseStatus: http.StatusOK,
			responseBody:   `[]`,
			expectedMRs:    []MergeRequest{},
			expectError:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock HTTP client
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					// Check that the request is properly formed
					if req.Method != "GET" {
						t.Errorf("expected GET request, got %s", req.Method)
					}

					// Check that the token is set
					if req.Header.Get("PRIVATE-TOKEN") != "test-token" {
						t.Errorf("expected PRIVATE-TOKEN header to be 'test-token', got %s", req.Header.Get("PRIVATE-TOKEN"))
					}

					// Return the mock response
					return &http.Response{
						StatusCode: tc.responseStatus,
						Body:       io.NopCloser(bytes.NewBufferString(tc.responseBody)),
					}, nil
				},
			}

			// Create a client with the mock HTTP client
			client := NewClient("test-token")
			client.HTTPClient = mockClient

			// Call the method
			mrs, err := client.GetMergeRequestsBySourceBranch(tc.projectID, tc.sourceBranch)

			// Check the results
			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(mrs) != len(tc.expectedMRs) {
					t.Errorf("expected %d merge requests, got %d", len(tc.expectedMRs), len(mrs))
				}
				// Check the first MR if there are any
				if len(mrs) > 0 && len(tc.expectedMRs) > 0 {
					if mrs[0].ID != tc.expectedMRs[0].ID {
						t.Errorf("expected MR ID %d, got %d", tc.expectedMRs[0].ID, mrs[0].ID)
					}
					if mrs[0].Title != tc.expectedMRs[0].Title {
						t.Errorf("expected MR title %q, got %q", tc.expectedMRs[0].Title, mrs[0].Title)
					}
				}
			}
		})
	}
}

// TestGetMergeRequestComments tests the GetMergeRequestComments method
func TestGetMergeRequestComments(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		projectID      string
		mrIID          int
		responseStatus int
		responseBody   string
		responseHeader http.Header
		expectedNotes  []MergeRequestNote
		expectError    bool
	}{
		{
			name:           "successful retrieval",
			projectID:      "12345",
			mrIID:          1,
			responseStatus: http.StatusOK,
			responseBody:   `[{"id":1,"body":"Test comment","author":{"username":"testuser"},"system":false,"created_at":"2023-01-01T00:00:00Z","resolved":false,"position":{"new_path":"test.go","new_line":10}}]`,
			responseHeader: http.Header{},
			expectedNotes: []MergeRequestNote{
				{
					ID:   1,
					Body: "Test comment",
					Author: struct {
						Username string "json:\"username\""
					}{
						Username: "testuser",
					},
					System:    false,
					CreatedAt: "2023-01-01T00:00:00Z",
					Resolved:  false,
					Position: &struct {
						NewPath *string "json:\"new_path\""
						NewLine *int    "json:\"new_line\""
					}{
						NewPath: strPtr("test.go"),
						NewLine: intPtr(10),
					},
				},
			},
			expectError: false,
		},
		{
			name:           "API error",
			projectID:      "12345",
			mrIID:          1,
			responseStatus: http.StatusInternalServerError,
			responseBody:   `{"error":"Internal server error"}`,
			responseHeader: http.Header{},
			expectedNotes:  nil,
			expectError:    true,
		},
		{
			name:           "no comments",
			projectID:      "12345",
			mrIID:          1,
			responseStatus: http.StatusOK,
			responseBody:   `[]`,
			responseHeader: http.Header{},
			expectedNotes:  []MergeRequestNote{},
			expectError:    false,
		},
		{
			name:           "pagination",
			projectID:      "12345",
			mrIID:          1,
			responseStatus: http.StatusOK,
			responseBody:   `[{"id":1,"body":"Test comment 1","author":{"username":"testuser"},"system":false,"created_at":"2023-01-01T00:00:00Z","resolved":false,"position":{"new_path":"test.go","new_line":10}}]`,
			responseHeader: http.Header{
				"X-Next-Page": []string{"2"},
			},
			expectedNotes: []MergeRequestNote{
				{
					ID:   1,
					Body: "Test comment 1",
					Author: struct {
						Username string "json:\"username\""
					}{
						Username: "testuser",
					},
					System:    false,
					CreatedAt: "2023-01-01T00:00:00Z",
					Resolved:  false,
					Position: &struct {
						NewPath *string "json:\"new_path\""
						NewLine *int    "json:\"new_line\""
					}{
						NewPath: strPtr("test.go"),
						NewLine: intPtr(10),
					},
				},
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock HTTP client
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					// Check that the request is properly formed
					if req.Method != "GET" {
						t.Errorf("expected GET request, got %s", req.Method)
					}

					// Check that the token is set
					if req.Header.Get("PRIVATE-TOKEN") != "test-token" {
						t.Errorf("expected PRIVATE-TOKEN header to be 'test-token', got %s", req.Header.Get("PRIVATE-TOKEN"))
					}

					// Return the mock response
					resp := &http.Response{
						StatusCode: tc.responseStatus,
						Body:       io.NopCloser(bytes.NewBufferString(tc.responseBody)),
						Header:     tc.responseHeader,
					}
					return resp, nil
				},
			}

			// Create a client with the mock HTTP client
			client := NewClient("test-token")
			client.HTTPClient = mockClient

			// Call the method
			notes, err := client.GetMergeRequestComments(tc.projectID, tc.mrIID)

			// Check the results
			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(notes) != len(tc.expectedNotes) {
					t.Errorf("expected %d notes, got %d", len(tc.expectedNotes), len(notes))
				}
				// Check the first note if there are any
				if len(notes) > 0 && len(tc.expectedNotes) > 0 {
					if notes[0].ID != tc.expectedNotes[0].ID {
						t.Errorf("expected note ID %d, got %d", tc.expectedNotes[0].ID, notes[0].ID)
					}
					if notes[0].Body != tc.expectedNotes[0].Body {
						t.Errorf("expected note body %q, got %q", tc.expectedNotes[0].Body, notes[0].Body)
					}
				}
			}
		})
	}
}

// Helper functions for creating pointers to string and int values
func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}