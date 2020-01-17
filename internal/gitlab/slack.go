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

func updateSlackService(client *gitlab.Client, p *gitlab.Project, cfg *Config, dryRun bool) error {
	// Fetch Slack Service settings from config file and return if a nil object is returned.
	cfgSettings := cfg.SlackSettings(p.Namespace.FullPath)
	// We don't error if no Slack settings are found for the namespace or a parent namespace, we just
	// print a message and move on...
	if cfgSettings == nil {
		fmt.Printf("Cannot find Slack settings for namespace %s or any of it's parents, skipping\n", p.Namespace.FullPath)
		return nil
	}
	if compareObjects(cfgSettings, &SlackSettings{}) {
		return nil
	}

	// Fetch current Slack Service settings and properties.
	projectSettings, _, err := client.Services.GetSlackService(p.ID)
	if err != nil {
		return nil
	}

	// Create a gitlab.SlackService object with our desired service and properties by
	// taking what is current and applying our changes on top, thus we do not override
	// any settings that we don't define in our config.
	newSettings := &gitlab.SlackService{}
	if err := copier.Copy(&newSettings, &cfgSettings); err != nil {
		return err
	}

	if err := mergo.Merge(newSettings.Properties, cfgSettings.Properties, mergo.WithOverride); err != nil {
		return err
	}
	newSettings.Active = cfgSettings.Active
	// This is the service default so we set it to true here; consumers can set it to
	// false via the config file if it should be disabled.
	newSettings.Properties.NotifyOnlyBrokenPipelines = true
	for _, e := range cfgSettings.Events {
		switch e {
		case "issues":
			newSettings.IssuesEvents = true
		case "merge_request":
			newSettings.MergeRequestsEvents = true
		case "pipeline":
			newSettings.PipelineEvents = true
		case "push":
			newSettings.PushEvents = true
		case "tags":
			newSettings.TagPushEvents = true
		default:
			fmt.Printf("Unsupported event type: %s\n", e)
		}
	}

	// We need to set these even though they aren't use otherwise compareObjects
	// will never return true.
	newSettings.ID = projectSettings.ID
	newSettings.Title = projectSettings.Title
	newSettings.CreatedAt = projectSettings.CreatedAt
	newSettings.UpdatedAt = projectSettings.UpdatedAt

	// Return if our proposed config matches the actual config
	if compareObjects(projectSettings, newSettings) {
		fmt.Printf("Project %s's Slack settings don't need updating\n", p.Name)

		return nil
	}

	fmt.Printf("Project %s's Slack settings need updating ... ", p.Name)

	if dryRun {
		fmt.Printf("skipping because this is a dry run\n")
		return nil
	}

	opts := &gitlab.SetSlackServiceOptions{}

	svcData, _ := json.Marshal(newSettings.Service)
	if err := json.Unmarshal(svcData, &opts); err != nil {
		return err
	}

	propData, _ := json.Marshal(newSettings.Properties)
	if err := json.Unmarshal(propData, &opts); err != nil {
		return err
	}

	fmt.Printf("Updating project ... ")

	_, err = client.Services.SetSlackService(p.ID, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Success!\n")

	return nil
}
