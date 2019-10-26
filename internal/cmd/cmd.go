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

package cmd

import (
	"fmt"
	"os"

	"github.com/shoekstra/repo-settings/internal/config"
	"github.com/shoekstra/repo-settings/internal/gitlab"
	"github.com/spf13/cobra"
)

var dryRun bool
var gitlabToken string
var gitlabURL string
var cfgFile string

// NewRepoDefaultsCmd represents the command
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo-settings",
		Short: "CLI to configure repository settings in a GitLab group or project.",
		Long: `
A simple CLI to configure repositories settings across various repository
hosters.

It reads a configuration file containing a GitHub organisation (not yet
supported) and/or GitLab group and will configure all repositories
found within with the defined settings.

To get started, create a configuration file and pass the --config option.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runCmd(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	// Add some flags.
	cmd.Flags().BoolVarP(&dryRun, "dry-drun", "d", false, "perform a dry run")
	cmd.Flags().StringVar(&gitlabToken, "gitlab-token", "", "GitLab API token")
	cmd.Flags().StringVar(&gitlabURL, "gitlab-url", "", "GitLab API URL")
	cmd.Flags().StringVarP(&cfgFile, "config", "c", "", "path to config file")

	return cmd
}

func runCmd() error {
	if err := validateCfgFile(); err != nil {
		return err
	}

	cfg, err := config.Load(cfgFile)
	if err != nil {
		return err
	}

	if cfg.GitLab != nil {
		if err := cfg.GitLab.LoadCreds(gitlabToken, gitlabURL); err != nil {
			return err
		}

		// Update projects in configured groups.
		if cfg.GitLab.Groups != nil {
			if err := gitlab.UpdateProjectsInGroups(cfg.GitLab, dryRun); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateCfgFile() error {
	// Print help if no config file is passed
	if cfgFile == "" {
		cmd := NewCmd()
		cmd.Help()
		os.Exit(0)
	}

	// Test config file exists
	if _, err := os.Stat(cfgFile); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Cannot open config file %s", cfgFile)
		}
	}

	return nil
}
