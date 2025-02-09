// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package loader

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/nanoninja/render/tmpl"
)

// FSLoader implements the Loader interface for the local filesystem.
// It provides access to templates stored on disk while ensuring secure
// access within the configured root directory.
type FSLoader struct {
	baseLoader
}

// NewFS creates a new filesystem-based template loader with the given configuration.
// It initializes the loader to access templates from the local filesystem
// according to the specified configuration settings.
//
// Example:
//
//	src := NewFS(LoaderConfig{
//	    Root:      "templates",
//	    Extension: ".html",
//	})
//
// The root path is cleaned and normalized during initialization.
// If Root is empty in the config, it defaults to the current directory (".").
func NewFS(c tmpl.LoaderConfig) (tmpl.Loader, error) {
	// Validate root directory exists and is accessible
	info, err := os.Stat(c.Root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: root directory does not exist", ErrInvalidRoot)
		}
		return nil, fmt.Errorf("%w: cannot access root directory: %v", ErrInvalidRoot, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%w: root path is not a directory", ErrInvalidRoot)
	}
	return &FSLoader{baseLoader: newBase(c, false)}, nil
}

// Load retrieves all template paths from the filesystem under the root directory.
// It returns relative paths to all templates that match the configured extension.
//
// Example paths returned (with Root="/templates", Extension=".html"):
//   - layouts/base.html
//   - users/index.html
//   - users/edit.html
//
// All returned paths use forward slashes regardless of the operating system.
func (l *FSLoader) Load(string) ([]string, error) {
	var templates []string

	err := filepath.Walk(l.root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories in results
		if info.IsDir() {
			return nil
		}
		// Check if the file matches our extension
		if !l.hasValidExtension(path) {
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			realPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}
			if !strings.HasPrefix(filepath.Clean(realPath), filepath.Clean(l.root)) {
				return fmt.Errorf("%w: symlink points outside root directory", ErrPathTraversal)
			}
		}
		// Convert to relative path
		relativePath, err := filepath.Rel(l.root, path)
		if err != nil {
			return fmt.Errorf("%w: failed to resolve relative path", ErrInvalidPath)
		}
		// Use consistent path separators
		relativePath = filepath.ToSlash(relativePath)
		templates = append(templates, relativePath)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}
	return templates, nil
}

// Read retrieves the content of a specific template by its relative path.
// The path should be relative to the loader's root directory.
//
// Example:
//
//	loader.Read("users/index.html") // reads /templates/users/index.html
//
// The method prevents path traversal attacks by validating the final path
// is still under the root directory.
func (l *FSLoader) Read(name string) ([]byte, error) {
	// Ensure we have a clean path that's relative to root
	path := filepath.Clean(filepath.Join(l.root, name))

	if !strings.HasPrefix(path, l.root) {
		return nil, fmt.Errorf("%w: attempt to read outside root directory", ErrPathTraversal)
	}
	// Security check: ensure the path is still under our root directory
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", tmpl.ErrTemplateNotFound, name)
		}
		return nil, fmt.Errorf("failed to read template %s: %w", name, err)
	}
	return content, nil
}

// Extension returns the configured file extension for templates.
// This is the extension specified in the LoaderConfig during initialization.
// If no extension was configured, it returns an empty string, meaning
// all files matching the pattern will be considered as templates.
//
// The returned extension includes the dot prefix if one was configured
// (e.g., ".html", ".gohtml").
func (l *FSLoader) Extension() string {
	return l.extension
}
