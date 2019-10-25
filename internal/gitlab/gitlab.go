//
// Copyright © 2019 Stephen Hoekstra
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package gitlab

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/xanzy/go-gitlab"
)

// Config represents the GitLab section of the config file.
type Config struct {
	APIToken *string     `json:"apitoken,omitempty"`
	APIURL   *string     `json:"apiurl,omitempty"`
	Defaults *Settings   `json:"defaults,omitempty"`
	Groups   []*Settings `json:"groups,omitempty"`
}

// Settings represents a group's settings.
type Settings struct {
	Name    string `json:"name,omitempty"`
	General struct {
		MergeRequestApprovals gitlab.ProjectApprovals `json:"merge_request_approvals,omitempty"`
	} `json:"general,omitempty"`
	Integrations struct {
		Slack SlackSettings `json:"slack,omitempty"`
	} `json:"integrations,omitempty"`
}

// SlackSettings represents a project's Slack settings.
type SlackSettings struct {
	Active     bool                          `json:"active"`
	Events     []string                      `json:"events"`
	Properties gitlab.SlackServiceProperties `json:"properties,omitempty"`
}

// LoadCreds accepts a token and url string; if these are empty it will attempt
// read the GITLAB_TOKEN and GITLAB_URL env vars as a source for credentials. If
// these are also empty it returns an error.
func (c *Config) LoadCreds(token, url string) error {
	if token == "" {
		token = os.Getenv("GITLAB_TOKEN")
	}
	if url == "" {
		url = os.Getenv("GITLAB_URL")
	}
	if token == "" || url == "" {
		return fmt.Errorf("Missing required API token and/or URL params")
	}

	c.APIToken = &token
	c.APIURL = &url

	return nil
}

// MergeRequestApprovalSettings will return the Merge Request Approval settings
// for a project by looking up it's namespace in the config.
func (c *Config) MergeRequestApprovalSettings(ns string) *gitlab.ProjectApprovals {
	for {
		// Loop through groups and return configured Slack Settings.
		for _, g := range c.Groups {
			if strings.EqualFold(g.Name, ns) {
				return &g.General.MergeRequestApprovals
			}
		}

		// Pop last name in namespace before trying again.
		s := strings.Split(ns, "/")
		if len(s) == 1 {
			break
		}
		ns = strings.Join(s[:len(s)-1], "/")
	}

	// This should never happen as we only look up settings for projects
	// we've already found.
	return nil
}

// SlackSettings will return the Slack settings for a project by looking up
// it's namespace in the config.
func (c *Config) SlackSettings(ns string) *SlackSettings {
	for {
		// Loop through groups and return configured Slack Settings.
		for _, g := range c.Groups {
			if strings.EqualFold(g.Name, ns) {
				return &g.Integrations.Slack
			}
		}

		// Pop last name in namespace before trying again.
		s := strings.Split(ns, "/")
		if len(s) == 1 {
			break
		}
		ns = strings.Join(s[:len(s)-1], "/")
	}

	// This should never happen as we only look up settings for projects
	// we've already found.
	return nil
}

// newClient returns a configured GitLab client.
func newClient(token, url string) (*gitlab.Client, error) {
	client := gitlab.NewClient(nil, token)
	if url != "" {
		client.SetBaseURL(url)
	}

	return client, nil
}

// compareObjects compares two objects and returns true if they match or false
// if they don't.
func compareObjects(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}
