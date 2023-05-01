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
	"log"
	"regexp"

	"github.com/xanzy/go-gitlab"
)

// Groups exists to provide helper methods to []*gitlab.Group
type Groups []*gitlab.Group

// getID returns a group's ID
func (gs Groups) getID(name string) (int, error) {
	for _, g := range gs {
		matched, _ := regexp.MatchString(fmt.Sprintf("(?i)groups/%s$", name), g.WebURL)
		if matched {
			return g.ID, nil
		}
	}

	return 0, fmt.Errorf("Cannot find group with name \"%s\"; if this is a subgroup include it's parent group(s) in the name", name)
}

// listGroups returns a slice containing all readable GitLab groups.
func listGroups(client *gitlab.Client) (Groups, error) {
	groups := Groups{}

	orderBy := "name"
	sort := "asc"
	opt := &gitlab.ListGroupsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 10,
			Page:    1,
		},
		OrderBy: &orderBy,
		Sort:    &sort,
	}

	for {
		gs, resp, err := client.Groups.ListGroups(opt)
		if err != nil {
			log.Fatal(err)
		}

		groups = append(groups, gs...)

		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		opt.Page = resp.NextPage
	}

	return groups, nil
}

// UpdateProjectsInGroups is the entry point for this package and will update any projects
// found within the groups defined in *Config.Groups.
func UpdateProjectsInGroups(cfg *Config) error {
	client, err := newClient(*cfg.APIToken, *cfg.APIURL)
	if err != nil {
		return err
	}

	groups, err := listGroups(client)
	if err != nil {
		return err
	}

	for _, g := range cfg.Groups {
		fmt.Printf("Looking up group with name \"%s\" ... ", g.Name)
		id, err := groups.getID(g.Name)
		if err != nil {
			return err
		}
		fmt.Printf("matched group to ID %d\n", id)

		includeSubGroups := true
		orderBy := "name"
		sort := "asc"
		opt := &gitlab.ListGroupProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 20,
				Page:    1,
			},
			IncludeSubGroups: &includeSubGroups,
			OrderBy:          &orderBy,
			Sort:             &sort,
		}

		for {
			ps, resp, err := client.Groups.ListGroupProjects(id, opt)
			if err != nil {
				log.Fatal(err)
			}

			for _, p := range ps {
				// General settings
				// - Update Merge Request Approval settings
				if err := updateMergeRequestAppovalsSettings(client, p, cfg); err != nil {
					return err
				}
				// Repository settings
				// - Protected Branch settings
				if err := updateProtectedBranchesSettings(client, p, cfg); err != nil {
					return err
				}
				// Integrations
				// - Update Slack integration
				if err := updateSlackService(client, p, cfg); err != nil {
					return err
				}
			}

			if resp.CurrentPage >= resp.TotalPages {
				break
			}

			opt.Page = resp.NextPage
		}
	}

	return nil
}
