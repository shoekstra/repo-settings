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
	"strings"

	"github.com/xanzy/go-gitlab"
)

type ProtectedBranchSetting struct {
	Name           string `json:"name,omitempty"`
	AllowedToMerge string `json:"allowed_to_merge,omitempty"`
	AllowedToPush  string `json:"allowed_to_push,omitempty"`
}

func updateProtectedBranchesSettings(client *gitlab.Client, p *gitlab.Project, cfg *Config) error {
	// Fetch Protected Branches settings from config file and return if a nil object is returned.
	cfgSettings := cfg.ProtectedBranchesSettings(p.Namespace.FullPath)
	if cfgSettings == nil {
		return nil
	}

	// Fetch current Protected Branches settings
	listOpts := &gitlab.ListProtectedBranchesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 10,
			Page:    1,
		},
	}
	projectSettings, _, err := client.ProtectedBranches.ListProtectedBranches(p.ID, listOpts)
	if err != nil {
		return nil
	}

	// Loop through protected branches mentioned in config file
	for _, cfgSetting := range cfgSettings {
		// Look for an existing protected branch setting; it's ok if nothing is found, we'll
		// just add another
		projectSetting := &gitlab.ProtectedBranch{}
		for _, ps := range projectSettings {
			if strings.EqualFold(ps.Name, cfgSetting.Name) {
				projectSetting = ps
			}
		}

		ma, err := setBranchAccess(cfgSetting.AllowedToMerge)
		if err != nil {
			return err
		}
		pa, err := setBranchAccess(cfgSetting.AllowedToPush)
		if err != nil {
			return err
		}
		newSetting := &gitlab.ProtectedBranch{
			Name:              cfgSetting.Name,
			MergeAccessLevels: []*gitlab.BranchAccessDescription{&ma},
			PushAccessLevels:  []*gitlab.BranchAccessDescription{&pa},
		}

		if compareObjects(projectSetting, newSetting) {
			fmt.Printf("Project %s's %s Branch Protection settings don't need updating\n", p.PathWithNamespace, projectSetting.Name)

			return nil
		}

		fmt.Printf("Project %s's %s Branch Protection settings need updating ... ", p.PathWithNamespace, projectSetting.Name)

		if cfg.DryRun {
			fmt.Printf("skipping because this is a dry run\n")
			return nil
		}

		setOpts := &gitlab.ProtectRepositoryBranchesOptions{
			Name:             &newSetting.Name,
			MergeAccessLevel: &newSetting.MergeAccessLevels[0].AccessLevel,
			PushAccessLevel:  &newSetting.PushAccessLevels[0].AccessLevel,
		}

		fmt.Printf("Updating project ... ")

		_, err = client.ProtectedBranches.UnprotectRepositoryBranches(p.ID, *setOpts.Name)
		if err != nil {
			return nil
		}
		_, _, err = client.ProtectedBranches.ProtectRepositoryBranches(p.ID, setOpts)
		if err != nil {
			return nil
		}
		fmt.Printf("Success!\n")

	}

	return nil
}

func setBranchAccess(s string) (gitlab.BranchAccessDescription, error) {
	if strings.EqualFold(s, "developers") {
		return gitlab.BranchAccessDescription{AccessLevel: 30, AccessLevelDescription: "Developers + Maintainers"}, nil
	}

	if strings.EqualFold(s, "maintainers") {
		return gitlab.BranchAccessDescription{AccessLevel: 40, AccessLevelDescription: "Maintainers"}, nil
	}

	if strings.EqualFold(s, "no one") {
		return gitlab.BranchAccessDescription{AccessLevel: 0, AccessLevelDescription: "No one"}, nil
	}

	return gitlab.BranchAccessDescription{}, fmt.Errorf("Invalid access type: %s", s)
}
