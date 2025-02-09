// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package tmpl

// Loader defines how templates are loaded and read from different sources.
// This interface allows for flexible template loading from various sources like
// embedded files, local filesystem, or even remote storage.
type Loader interface {
	// Load retrieves all template paths matching the given pattern.
	// It returns a list of template names that can be used for reading.
	Load(pattern string) ([]string, error)

	// Read retrieves the content of a specific template by its name.
	// The name should be one of those returned by Load.
	Read(name string) ([]byte, error)

	// Extension returns the template file extension (e.g., ".html", ".gohtml")
	Extension() string
}

// LoaderConfig provides configuration options for template loaders.
// It defines how templates are located, filtered, and loaded from the source.
// This configuration can be used with both filesystem and embedded template loaders.
type LoaderConfig struct {
	// Root specifies the base directory for template loading.
	// For filesystem loader: can be relative or absolute path (e.g., "./templates" or "/app/templates")
	// For embedded loader: must be the directory name in the embedded filesystem
	// If empty, defaults to the current directory (".")
	Root string

	// Extension defines the file extension for template files (e.g., ".html", ".gohtml").
	// When set, only files with this extension will be considered as templates.
	// The extension should include the dot prefix.
	// If empty, all files matching the pattern will be loaded.
	Extension string
}
