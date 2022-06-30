// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/joshdk/template-transformer/config"
	"gopkg.in/yaml.v3"
	"jdk.sh/meta"
)

func main() {
	if err := mainCmd(); err != nil {
		fmt.Fprintln(os.Stderr, "github.com/joshdk/template-transformer:", err)
		os.Exit(1)
	}
}

func mainCmd() error {
	if len(os.Args) >= 2 && os.Args[1] == "--version" {
		version()
		return nil
	}
	log.Println("https://github.com/joshdk/template-transformer version", meta.Version())

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

	// Log the resolved values for each property.
	for key, value := range properties {
		log.Printf("property %q resolved to %q", key, value)
	}

	// Read resources from the input stream, transform them, and write them
	// back to the output stream.
	return transform(os.Stdin, os.Stdout, properties)
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
				matches := property.Mutate.Regex.FindStringSubmatchIndex(value)

				// Check that the original value actually matches the regex.
				if len(matches) == 0 {
					return "", fmt.Errorf("value %q did not match mutate regex", value)
				}

				// Expand the template with the matched values.
				replaced := property.Mutate.Regex.ExpandString(nil, property.Mutate.Replace, value, matches)
				if len(replaced) == 0 {
					return "", fmt.Errorf("replaced value was blank")
				}

				// Return the configured capture group.
				value = string(replaced)
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

// transform reads resources from the input stream, modifies them, and writes
// them back to the output stream.
func transform(in io.Reader, out io.Writer, properties map[string]string) error {
	decoder := yaml.NewDecoder(in)
	for {
		// Unmarshal a single yaml document (kubernetes resource) from the
		// input stream.
		var resource interface{}
		if err := decoder.Decode(&resource); err != nil {
			// Quit processing if there are no more documents left in the
			// stream.
			if err == io.EOF {
				return nil
			}
			return err
		}

		// Immediately marshal the resource back into yaml as we will need to
		// template it as an opaque slice of bytes.
		data, err := yaml.Marshal(resource)
		if err != nil {
			return err
		}

		// These are the template delimiters used for referencing named
		// properties like "${{.PROPERTY}}".
		// These delimiters were chosen to:
		// - Survive round-tripping through the yaml parser.
		// - Not conflict with the kustomize vars syntax.
		// - Not conflict with kubernetes var references.
		// - Not conflict with various types of bash expansion.
		const prefix, suffix = "${{", "}}"

		// Parse the entire resource body as a text template.
		tpl, err := template.New("stream").Option("missingkey=error").Delims(prefix, suffix).Parse(string(data))
		if err != nil {
			return err
		}

		// Execute the template and write the result to the output stream.
		if err := tpl.Execute(out, properties); err != nil {
			return err
		}

		// Write a yaml document separator so that multiple documents can be written to the output stream.
		if _, err := out.Write([]byte("---\n")); err != nil {
			return err
		}
	}
}

func version() {
	fmt.Println("homepage: https://github.com/joshdk/template-transformer")
	fmt.Println("author:   Josh Komoroske")
	fmt.Println("license:  MIT")
	if meta.Version() == "" {
		return
	}
	fmt.Println("version: ", meta.Version())
	fmt.Println("sha:     ", meta.ShortSHA())
	fmt.Println("date:    ", meta.DateFormat(time.RFC3339))
}
