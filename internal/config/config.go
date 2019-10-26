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

package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/shoekstra/repo-settings/internal/gitlab"
	"github.com/spf13/viper"
)

// Config represents the app config structure.
type Config struct {
	GitLab *gitlab.Config `json:"gitlab,omitempty"`
}

// Load reads a config file and returns an initialised Config object.
func Load(path string) (*Config, error) {
	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	if ok := contains([]string{"json", "yaml", "yml"}, ext); !ok {
		return nil, fmt.Errorf("Unsupported config type \"%s\"", ext)
	}

	viper.SetConfigType(ext)
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %s", err)
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %s", err)
	}

	return cfg, err
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
