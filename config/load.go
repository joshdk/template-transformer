// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Load parses the contents of the given filename as YAML and returns a Config.
func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename) //nolint:gosec
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Validate the config struct for correctness.
	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
