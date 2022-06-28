// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joshdk/template-transformer/config"
)

func main() {
	if err := mainCmd(); err != nil {
		fmt.Fprintln(os.Stderr, "github.com/joshdk/template-transformer:", err)
		os.Exit(1)
	}
}

func mainCmd() error {
	// When invoked via "kustomize build" the first (os.Args[1]) is a temporary
	// filename containing the plugin configuration. Validate that we are being
	// properly run as a plugin.
	if len(os.Args) < 2 {
		return errors.New("not invoked as a kustomize plugin")
	}

	// Parse the plugin configuration file.
	cfg, err := config.Load(os.Args[1])
	if err != nil {
		return err
	}

	// For now, just log the plugin configuration.
	log.Printf("plugin config: %+v", cfg)

	// Resolve all property values and build a complete map of property names
	// to their respective values.
	properties := make(map[string]string)
	for _, property := range cfg.Properties {
		value, err := resolve(property)
		if err != nil {
			return config.PathError{
				Wrapped: err,
				Paths:   []string{"properties", property.Name},
			}
		}

		properties[property.Name] = value
	}

	// For now, just log the property values.
	log.Printf("properties: %+v", properties)

	return nil
}

// resolve obtains a final value for the given property.
func resolve(property config.Property) (string, error) {
	// Check every environment variable, in the order they were specified, and
	// return the first value that is found.
	for _, source := range property.Source {
		// If this environment variable contains a value then return it.
		if value := os.Getenv(source); value != "" {
			// If a regex was configured then mutate the value.
			if property.Mutate.Regex != nil {
				// List the regex matches for the original value.
				matches := property.Mutate.Regex.FindStringSubmatch(value)

				// Check that the original value actually matches the regex.
				if len(matches) == 0 {
					return "", fmt.Errorf("value %q did not match mutate regex", value)
				}

				// Return the configured capture group.
				value = matches[property.Mutate.Capture]
			}

			return value, nil
		}
	}

	// No environment variables contained a value. If a default is configured
	// then return it.
	if property.Default != "" {
		return property.Default, nil
	}

	// No value could be obtained.
	return "", errors.New("could not resolve value for property")
}
