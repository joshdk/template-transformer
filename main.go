// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"os"
)

func main() {
	if err := mainCmd(); err != nil {
		fmt.Fprintln(os.Stderr, "github.com/joshdk/template-transformer:", err)
		os.Exit(1)
	}
}

func mainCmd() error {
	return nil
}
