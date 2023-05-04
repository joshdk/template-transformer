// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package config

import (
	"fmt"
	"regexp"
)

func validate(cfg *Config) error { //nolint:cyclop
	// At least one property needs to be configured.
	if len(cfg.Properties) == 0 {
		return PathError{
			Message: "no properties",
			Paths:   []string{"properties"},
		}
	}

	for propertyIndex, property := range cfg.Properties {
		// A property must have a name.
		if property.Name == "" {
			return PathError{
				Message: "name is empty",
				Paths:   []string{"properties", fmt.Sprint(propertyIndex), "name"},
			}
		}

		// A property must have a description.
		if property.Description == "" {
			return PathError{
				Message: "description is empty",
				Paths:   []string{"properties", property.Name, "description"},
			}
		}

		// At least one source needs to be configured.
		if len(property.Source) == 0 {
			return PathError{
				Message: "no sources",
				Paths:   []string{"properties", property.Name, "source"},
			}
		}

		for sourceIndex, source := range property.Source {
			// Source names cannot be blank.
			if source == "" {
				return PathError{
					Message: "source is empty",
					Paths:   []string{"properties", property.Name, "source", fmt.Sprint(sourceIndex)},
				}
			}
		}

		if property.Mutate.Pattern != "" {
			// The regex needs to actually compile.
			regex, err := regexp.Compile(property.Mutate.Pattern)
			if err != nil {
				return PathError{
					Wrapped: err,
					Paths:   []string{"properties", property.Name, "mutate", "pattern"},
				}
			}

			// Replace cannot be blank.
			if property.Mutate.Replace == "" {
				return PathError{
					Message: "no value",
					Paths:   []string{"properties", property.Name, "mutate", "replace"},
				}
			}

			// Cache the compiled regex for later use.
			property.Mutate.Regex = regex
		}

		// If no regex is configured than replace must also not be configured.
		if property.Mutate.Pattern == "" && property.Mutate.Replace != "" {
			return PathError{
				Message: "no value",
				Paths:   []string{"properties", property.Name, "mutate", "pattern"},
			}
		}

		// Persist the updated property.
		cfg.Properties[propertyIndex] = property
	}

	return nil
}
