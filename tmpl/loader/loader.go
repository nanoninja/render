// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package loader

import (
	"path/filepath"
	"strings"

	"github.com/nanoninja/render/tmpl"
)

// baseLoader provides common functionality for different template loader implementations.
// It handles basic configuration and shared operations for both filesystem and embedded loaders.
type baseLoader struct {
	// root defines the base directory for template loading
	root string

	// extension specifies the template file extension to filter (e.g., ".html")
	extension string

	// separator is the path separator character used for path manipulation
	// It's "/" for embed.FS and os-specific for filesystem operations
	separator string
}

// newBase creates a new baseLoader instance with appropriate configuration.
// The forEmbed parameter determines which path separator to use:
// - true: uses forward slash for embed.FS
// - false: uses os.PathSeparator for filesystem operations
func newBase(c tmpl.LoaderConfig, forEmbed bool) baseLoader {
	separator := string(filepath.Separator)
	if forEmbed {
		separator = "/"
	}
	return baseLoader{
		root:      c.Root,
		extension: c.Extension,
		separator: separator,
	}
}

// hasValidExtension checks if the given path matches the configured extension.
// If no extension is configured (empty string), all files are considered valid.
func (l *baseLoader) hasValidExtension(path string) bool {
	if l.extension == "" {
		return true
	}
	return strings.HasSuffix(path, l.extension)
}
