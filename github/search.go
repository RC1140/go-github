// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/url"
	"strconv"
)

// SearchService provides access to the search related functions
// in the GitHub API.
//
// GitHub API docs: http://developer.github.com/v3/search/
type SearchService struct {
	client *Client
}

// SearchOptions specifies optional parameters to the SearchService methods.
type SearchOptions struct {
	// How to sort the search results.  Possible values are:
	//   - for repositories: stars, fork, updated
	//   - for code: indexed
	//   - for issues: comments, created, updated
	//   - for users: followers, repositories, joined
	//
	// Default is to sort by best match.
	Sort string

	// Sort order if sort parameter is provided. Possible values are: asc,
	// desc. Default is desc.
	Order string

	// Page of results to retrieve.
	Page int

	// Number of results to show per page.  This can be up to 100
	// according to GitHub.
	PerPage int
}

// RepositoriesSearchResult represents the result of a repositories search.
type RepositoriesSearchResult struct {
	Total        *int         `json:"total_count,omitempty"`
	Repositories []Repository `json:"items,omitempty"`
}

// Repositories searches repositories via various criteria.
//
// GitHub API docs: http://developer.github.com/v3/search/#search-repositories
func (s *SearchService) Repositories(query string, opt *SearchOptions) (*RepositoriesSearchResult, *Response, error) {
	result := new(RepositoriesSearchResult)
	resp, err := s.search("repositories", query, opt, result)
	return result, resp, err
}

// IssuesSearchResult represents the result of an issues search.
type IssuesSearchResult struct {
	Total  *int    `json:"total_count,omitempty"`
	Issues []Issue `json:"items,omitempty"`
}

// Issues searches issues via various criteria.
//
// GitHub API docs: http://developer.github.com/v3/search/#search-issues
func (s *SearchService) Issues(query string, opt *SearchOptions) (*IssuesSearchResult, *Response, error) {
	result := new(IssuesSearchResult)
	resp, err := s.search("issues", query, opt, result)
	return result, resp, err
}

// UsersSearchResult represents the result of an issues search.
type UsersSearchResult struct {
	Total *int   `json:"total_count,omitempty"`
	Users []User `json:"items,omitempty"`
}

// Users searches users via various criteria.
//
// GitHub API docs: http://developer.github.com/v3/search/#search-users
func (s *SearchService) Users(query string, opt *SearchOptions) (*UsersSearchResult, *Response, error) {
	result := new(UsersSearchResult)
	resp, err := s.search("users", query, opt, result)
	return result, resp, err
}

// CodeSearchResult represents the result of an code search.
type CodeSearchResult struct {
	Total       *int         `json:"total_count,omitempty"`
	CodeResults []CodeResult `json:"items,omitempty"`
}

// CodeResult represents a single search result.
type CodeResult struct {
	Name       *string     `json:"name,omitempty"`
	Path       *string     `json:"path,omitempty"`
	SHA        *string     `json:"sha,omitempty"`
	HTMLURL    *string     `json:"html_url,omitempty"`
	Repository *Repository `json:"repository,omitempty"`
}

func (c CodeResult) String() string {
	return Stringify(c)
}

// Code searches code via various criteria.
//
// GitHub API docs: http://developer.github.com/v3/search/#search-code
func (s *SearchService) Code(query string, opt *SearchOptions) (*CodeSearchResult, *Response, error) {
	result := new(CodeSearchResult)
	resp, err := s.search("code", query, opt, result)
	return result, resp, err
}

// Helper function that executes search queries against different
// GitHub search types (repositories, code, issues, users)
func (s *SearchService) search(searchType string, query string, opt *SearchOptions, result interface{}) (*Response, error) {
	params := url.Values{"q": []string{query}}
	if opt != nil {
		params.Add("sort", opt.Sort)
		params.Add("order", opt.Order)
		params.Add("page", strconv.Itoa(opt.Page))
		params.Add("per_page", strconv.Itoa(opt.PerPage))
	}
	u := fmt.Sprintf("search/%s?%s", searchType, params.Encode())

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", mimePreview)

	resp, err := s.client.Do(req, result)
	return resp, err
}
