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

	return nil
}
