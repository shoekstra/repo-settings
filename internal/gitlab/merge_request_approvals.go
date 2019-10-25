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
	"encoding/json"
	"fmt"

	"github.com/imdario/mergo"
	"github.com/jinzhu/copier"
	"github.com/xanzy/go-gitlab"
)

func updateMergeRequestAppovalsSettings(client *gitlab.Client, p *gitlab.Project, cfg *Config, dryRun bool) error {
	// Fetch Merge Request Approval settings from config file and return if a nil object is returned.
	cfgSettings := cfg.MergeRequestApprovalSettings(p.Namespace.FullPath)
	if compareObjects(cfgSettings, &gitlab.ProjectApprovals{}) {
		return nil
	}

	// Fetch current Merge Request Approval settings.
	projectSettings, _, err := client.Projects.GetApprovalConfiguration(p.ID)
	if err != nil {
		return nil
	}

	// Populate newSettings with cfgSettings values, we do this otherwise mergo.Merge wigs out and raises
	// an exception.
	newSettings := &gitlab.ProjectApprovals{}
	if err := copier.Copy(&newSettings, &cfgSettings); err != nil {
		return err
	}

	// Set defaults, otherwise these will be set as nil values causing the compare to always fail.
	newSettings.Approvers = []*gitlab.MergeRequestApproverUser{}
	newSettings.ApproverGroups = []*gitlab.MergeRequestApproverGroup{}

	// Merge our changes on top of existing settings.
	if err := mergo.Merge(newSettings, projectSettings, mergo.WithOverride); err != nil {
		return err
	}

	// Return if our proposed config matches the actual config
	if compareObjects(projectSettings, newSettings) {
		fmt.Printf("Project %s's Merge Request Approval settings don't need updating\n", p.Name)

		return nil
	}

	fmt.Printf("Project %s's Merge Request Approval settings need updating ... ", p.Name)

	if dryRun {
		fmt.Printf("skipping because this is a dry run\n")
		return nil
	}

	opts := &gitlab.ChangeApprovalConfigurationOptions{}

	settingsData, _ := json.Marshal(newSettings)
	if err := json.Unmarshal(settingsData, &opts); err != nil {
		return err
	}

	fmt.Printf("Updating project ... ")

	_, _, err = client.Projects.ChangeApprovalConfiguration(p.ID, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Success!\n")

	return nil
}
