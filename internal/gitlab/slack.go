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
	// Fetch current Slack Service settings and properties.
	projectSlack, _, err := client.Services.GetSlackService(p.ID)
	if err != nil {
		return nil
	}

	// Create a gitlab.SlackService object with our desired service and properties by
	// taking what is current and applying our changes on top, thus we do not override
	// any settings that we don't define in our config.
	cfgSlack := cfg.SlackSettings(p.Namespace.FullPath)
	newSlack := &gitlab.SlackService{}
	if err := copier.Copy(&newSlack, &cfgSlack); err != nil {
		return err
	}

	if err := mergo.Merge(newSlack.Properties, cfgSlack.Properties, mergo.WithOverride); err != nil {
		return err
	}
	newSlack.Active = cfgSlack.Active
	// This is the service default so we set it to true here; consumers can set it to
	// false via the config file if it should be disabled.
	newSlack.Properties.NotifyOnlyBrokenPipelines = true
	for _, e := range cfgSlack.Events {
		switch e {
		case "issues":
			newSlack.IssuesEvents = true
		case "merge_request":
			newSlack.MergeRequestsEvents = true
		case "pipeline":
			newSlack.PipelineEvents = true
		case "push":
			newSlack.PushEvents = true
		case "tags":
			newSlack.TagPushEvents = true
		default:
			fmt.Printf("Unsupported event type: %s\n", e)
		}
	}

	// We need to set these even though they aren't use otherwise compareObjects
	// will never return true.
	newSlack.ID = projectSlack.ID
	newSlack.Title = projectSlack.Title
	newSlack.CreatedAt = projectSlack.CreatedAt
	newSlack.UpdatedAt = projectSlack.UpdatedAt

	// Return if our proposed config matches the actual config
	if compareObjects(projectSlack, newSlack) {
		fmt.Printf("Project %s's Slack settings doesn't need updating\n", p.Name)

		return nil
	}

	fmt.Printf("Project %s's Slack settings need updating ... ", p.Name)

	if dryRun {
		fmt.Printf("skipping because this is a dry run\n")
		return nil
	}

	opts := &gitlab.SetSlackServiceOptions{}

	svcData, _ := json.Marshal(newSlack.Service)
	if err := json.Unmarshal(svcData, &opts); err != nil {
		return err
	}

	propData, _ := json.Marshal(newSlack.Properties)
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
