// Package gitlab provides utilities for interacting with the GitLab API.
package gitlab

// MergeRequest represents a GitLab merge request.
type MergeRequest struct {
	ID           int    `json:"id"`
	IID          int    `json:"iid"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	WebURL       string `json:"web_url"`
	SourceBranch string `json:"source_branch"`
	TargetBranch string `json:"target_branch"`
	State        string `json:"state"`
	Author       *struct {
		UserName *string `json:"username"` // Use pointer in case it's null
	} `json:"author,omitempty"`
}

// MergeRequestNote represents a comment on a GitLab merge request.
type MergeRequestNote struct {
	ID     int    `json:"id"`
	Body   string `json:"body"`
	Author struct {
		Username string `json:"username"`
	} `json:"author"`
	System    bool   `json:"system"`
	CreatedAt string `json:"created_at"`
	Resolved  bool   `json:"resolved"`
	Position  *struct {
		NewPath   *string `json:"new_path"` // Use pointer in case it's null
		NewLine   *int    `json:"new_line"` // Use pointer in case it's null
		LineRange *struct {
			Start *struct {
				LineCode string `json:"line_code"`
			} `json:"start,omitempty"`
		} `json:"line_range,omitempty"`
	} `json:"position,omitempty"`
}
