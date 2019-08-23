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
	Name         string `json:"name,omitempty"`
	Integrations struct {
		Slack SlackSettings `json:"slack,omitempty"`
	} `json:"integrations,omitempty"`
}

// SlackSettings represents a project's Slack settings.
type SlackSettings struct {
	Active     bool                           `json:"active"`
	Events     []string                       `json:"events"`
	Properties *gitlab.SlackServiceProperties `json:"properties,omitempty"`
}

// SlackSettings will return the Slack settings for a project by looking up the
// it's name space in the config.
func (c *Config) SlackSettings(ns string) SlackSettings {
	for {
		// Loop through groups and return configured Slack Settings.
		for _, g := range c.Groups {
			if strings.EqualFold(g.Name, ns) {
				return g.Integrations.Slack
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
	return SlackSettings{}
}

// newClient returns a configured GitLab client.
func newClient(apiToken, apiURL string) (*gitlab.Client, error) {
	// Test API token and URL params are configured
	if apiToken == "" && apiURL == "" {
		return nil, fmt.Errorf("Missing required API token and/or URL params")
	}

	client := gitlab.NewClient(nil, apiToken)
	if apiURL != "" {
		client.SetBaseURL(apiURL)
	}

	return client, nil
}

// compareObjects compares two objects and returns true if they match or false
// if they don't.
func compareObjects(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

// contains checks a slice for a string and returns true if found.
func contains(s []string, str string) bool {
	for _, n := range s {
		if str == n {
			return true
		}
	}
	return false
}
