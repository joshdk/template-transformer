// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package config

import (
	"fmt"
	"path/filepath"
)

// PathError represents an error that is associated with some hierarchical
// resource.
type PathError struct {
	// Wrapped is a wrapped error.
	Wrapped error

	// Message is a plane error message.
	Message string

	// Paths is a list of hierarchical path segments.
	Paths []string
}

func (err PathError) Error() string {
	path := "/" + filepath.Join(err.Paths...)
	if err.Wrapped != nil {
		return fmt.Sprintf("%s: %v", path, err.Wrapped)
	}

	return fmt.Sprintf("%s: %s", path, err.Message)
}
