// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package loader

import "errors"

var (
	// ErrInvalidRoot is returned when the root path is invalid.
	// This can occur if the path doesn't exist or is not accessible.
	ErrInvalidRoot = errors.New("invalid root path")

	// ErrInvalidPath is returned when a template path cannot be properly resolved
	// or processed. This can occur in several situations:
	//   - Path contains invalid characters
	//   - Path cannot be converted to a relative path
	//   - Path normalization fails
	//
	// This error helps distinguish path processing issues from other types of
	// template loading errors, making debugging easier for library users.
	ErrInvalidPath = errors.New("invalid template path")

	// ErrPathTraversal is returned when detecting an attempt to access files
	// outside of the root directory through path manipulation.
	ErrPathTraversal = errors.New("path traversal attempt detected")
)
