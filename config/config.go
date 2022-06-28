// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package config

import "regexp"

// Config represents a kustomize plugin configuration.
type Config struct {
	// Properties is a list of named properties used for templating.
	Properties []Property `yaml:"properties"`
}

// Property represents a single named value that can be obtained from the
// current working environment.
type Property struct {
	// Name is the name used while templating.
	Name string `yaml:"name"`

	// Description is a textual description of what this property is and what
	// it is intended to be used for.
	Description string `yaml:"description"`

	// Source is an ordered list of environment variable names used for
	// obtaining a value.
	Source []string `yaml:"source"`

	// Default is an optional fallback value used when none of the source
	// environment variables exist or contain a value.
	Default string `yaml:"default"`

	// Mutate is an optional operation to transform the final value by using a
	// regex capture.
	Mutate Mutate `yaml:"mutate"`
}

// Mutate represents a transform that can be applied to a value by using a
// regex capture.
type Mutate struct {
	// Pattern is a regex used to match against a value.
	Pattern string `yaml:"pattern"`

	// Capture is the regex capture group used to transform the value.
	Capture int `yaml:"capture"`

	// Regex is a compiled regex that can be reused.
	Regex *regexp.Regexp `yaml:"-"`
}
