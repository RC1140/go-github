// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// UsersService handles communication with the user related
// methods of the GitHub API.
//
// GitHub API docs: http://developer.github.com/v3/users/
type UsersService struct {
	client *Client
}

// User represents a GitHub user.
type User struct {
	Login       *string    `json:"login,omitempty"`
	ID          *int       `json:"id,omitempty"`
	URL         *string    `json:"url,omitempty"`
	AvatarURL   *string    `json:"avatar_url,omitempty"`
	GravatarID  *string    `json:"gravatar_id,omitempty"`
	Name        *string    `json:"name,omitempty"`
	Company     *string    `json:"company,omitempty"`
	Blog        *string    `json:"blog,omitempty"`
	Location    *string    `json:"location,omitempty"`
	Email       *string    `json:"email,omitempty"`
	Hireable    *bool      `json:"hireable,omitempty"`
	PublicRepos *int       `json:"public_repos,omitempty"`
	Followers   *int       `json:"followers,omitempty"`
	Following   *int       `json:"following,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}

func (u User) String() string {
	return Stringify(u)
}
// FollowingUser represents a subset of fields of a GitHub user.
// but also contains information that can be used to query
// the users details further
type FollowingUser struct {
	User
	Type              string `json:"type,omitempty"`
	FollowersURL      string `json:"followers_url,omitempty"`
	FollowingURL      string `json:"following_url,omitempty"`
	GistsURL          string `json:"gists_url,omitempty"`
	StarredURL        string `json:"starred_url,omitempty"`
	SubscriptionsURL  string `json:"subscriptions_url,omitempty"`
	OrganizationsURL  string `json:"organizations_url,omitempty"`
	ReposURL          string `json:"repos_url,omitempty"`
	EventsURL         string `json:"events_url,omitempty"`
	ReceivedEventsURL string `json:"received_events_url,omitempty"`
}

// UserEmail represents user's email address
type UserEmail string

// Get fetches a user.  Passing the empty string will fetch the authenticated
// user.
//
// GitHub API docs: http://developer.github.com/v3/users/#get-a-single-user
func (s *UsersService) Get(user string) (*User, *Response, error) {
	var u string
	if user != "" {
		u = fmt.Sprintf("users/%v", user)
	} else {
		u = "user"
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	uResp := new(User)
	resp, err := s.client.Do(req, uResp)
	return uResp, resp, err
}

// Edit the authenticated user.
//
// GitHub API docs: http://developer.github.com/v3/users/#update-the-authenticated-user
func (s *UsersService) Edit(user *User) (*User, *Response, error) {
	u := "user"
	req, err := s.client.NewRequest("PATCH", u, user)
	if err != nil {
		return nil, nil, err
	}

	uResp := new(User)
	resp, err := s.client.Do(req, uResp)
	return uResp, resp, err
}

// UserListOptions specifies optional parameters to the UsersService.List
// method.
type UserListOptions struct {
	// ID of the last user seen
	Since int
}

// ListAll lists all GitHub users.
//
// GitHub API docs: http://developer.github.com/v3/users/#get-all-users
func (s *UsersService) ListAll(opt *UserListOptions) ([]User, *Response, error) {
	u := "users"
	if opt != nil {
		params := url.Values{
			"since": []string{strconv.Itoa(opt.Since)},
		}
		u += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	users := new([]User)
	resp, err := s.client.Do(req, users)
	return *users, resp, err
}

// ListFollowers lists an authenticated users followers or a specified users followers
//
// http://developer.github.com/v3/users/followers/#list-followers-of-a-user
func (s *UsersService) ListFollowers(user string) ([]FollowingUser, error) {
	u := "/user/followers"
	if user != "" {
		u = fmt.Sprintf("/users/%v/followers", user)
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	users := new([]FollowingUser)
	_, err = s.client.Do(req, users)
	return *users, err
}

// FollowingUser this lets you determine if you are following a user or not
//
// TODO : replace me with correct url
func (s *UsersService) FollowingUser(user string) (bool, error) {
    u := fmt.Sprintf("/user/following/%v", user)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return false, err
	}

    resp , err := s.client.Do(req,nil)
    if resp.StatusCode == 202{
        return true,err
    }else if resp.StatusCode == 404{
        return false,err
    }
    //handle returning a erros
	return false, err
}
