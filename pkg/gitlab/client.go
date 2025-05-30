// Package gitlab provides utilities for interacting with the GitLab API.
package gitlab

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client represents a GitLab API client.
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient HTTPClient
}

// NewClient creates a new GitLab API client with the given token.
func NewClient(token string) *Client {
	return &Client{
		BaseURL:    "https://gitlab.com/api/v4",
		Token:      token,
		HTTPClient: &http.Client{},
	}
}

// GetMergeRequestComments retrieves comments for a specific merge request.
func (c *Client) GetMergeRequestComments(projectID string, mrIID int) ([]MergeRequestNote, error) {
	perPage := 100
	page := 1
	var allNotes []MergeRequestNote

	for {
		endpoint := fmt.Sprintf(
			"%s/projects/%s/merge_requests/%d/notes?per_page=%d&page=%d",
			c.BaseURL, url.PathEscape(projectID), mrIID, perPage, page,
		)

		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("PRIVATE-TOKEN", c.Token)

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("GitLab API error: %s - %s", resp.Status, string(body))
		}

		var notes []MergeRequestNote
		if err := json.NewDecoder(resp.Body).Decode(&notes); err != nil {
			return nil, err
		}

		allNotes = append(allNotes, notes...)

		// GitLab uses `X-Next-Page` header for pagination
		nextPage := resp.Header.Get("X-Next-Page")
		if nextPage == "" {
			break
		}
		page++
	}

	return allNotes, nil
}

// GetMergeRequestsBySourceBranch retrieves merge requests for a specific source branch.
func (c *Client) GetMergeRequestsBySourceBranch(projectID, sourceBranch string) ([]MergeRequest, error) {
	endpoint := fmt.Sprintf("%s/projects/%s/merge_requests?source_branch=%s", 
		c.BaseURL, url.PathEscape(projectID), url.QueryEscape(sourceBranch))

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitLab API error: %s - %s", resp.Status, string(body))
	}

	var mrs []MergeRequest
	err = json.NewDecoder(resp.Body).Decode(&mrs)
	if err != nil {
		return nil, err
	}

	return mrs, nil
}
