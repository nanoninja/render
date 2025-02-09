// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package loader

import (
	"errors"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/nanoninja/render/tmpl"
)

// EmbedLoader implements Loader for embedded file systems.
// It provides secure and efficient access to templates embedded in the binary.
type EmbedLoader struct {
	baseLoader
	fs fs.ReadFileFS
}

// NewEmbed creates a new template loader that works with Go's embed.FS.
// It provides access to templates that have been embedded into the binary
// at compile time.
//
// The fs parameter must be an embed.FS that contains the template files.
// The configuration specifies how templates should be located and filtered
// within the embedded filesystem.
//
// Example:
//
//	//go:embed templates/*
//	var templateFS embed.FS
//
//	src := NewEmbed(templateFS, LoaderConfig{
//	    Root:      "templates",
//	    Extension: ".html",
//	})
//
// Note that embed.FS always uses forward slashes ('/') regardless of the operating system.
func NewEmbed(fs fs.ReadFileFS, c tmpl.LoaderConfig) tmpl.Loader {
	return &EmbedLoader{
		baseLoader: newBase(c, true),
		fs:         fs,
	}
}

// Load retrieves all template paths from the embedded filesystem under the root directory.
// It returns relative paths to all templates that match the configured extension.
//
// Example paths returned (with Root="templates", Extension=".html"):
//   - layouts/base.html
//   - users/index.html
//   - users/edit.html
//
// Note: embed.FS always uses forward slashes ('/') in paths.
func (l *EmbedLoader) Load(_ string) ([]string, error) {
	// For embed.FS, we need to read the entire directory structure
	var templates []string

	err := fs.WalkDir(l.fs, l.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !l.hasValidExtension(path) {
			return nil
		}
		// Convert to relative path by removing root prefix
		relativePath := strings.TrimPrefix(path, l.root)
		relativePath = strings.TrimPrefix(relativePath, "/")
		templates = append(templates, relativePath)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk embedded filesystem: %w", err)
	}
	return templates, nil
}

// Read retrieves the content of a specific template from the embedded filesystem.
// The path should be relative to the loader's root directory.
//
// Note: embed.FS requires forward slashes ('/') in paths, so this method
// converts any OS-specific path separators to forward slashes.
func (l *EmbedLoader) Read(name string) ([]byte, error) {
	// Convert path separators to forward slashes for embed.FS
	name = filepath.ToSlash(name)

	content, err := l.fs.ReadFile(path.Join(l.root, name))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("%w: %s", tmpl.ErrTemplateNotFound, name)
		}
		return nil, fmt.Errorf("failed to read embedded template %s: %w", name, err)
	}
	return content, nil
}

// Extension returns the configured file extension for templates in the embedded filesystem.
// This is the extension specified in the LoaderConfig during initialization.
// If no extension was configured, it returns an empty string, meaning
// all files matching the pattern will be considered as templates.
//
// The returned extension includes the dot prefix if one was configured
// (e.g., ".html", ".gohtml").
func (l *EmbedLoader) Extension() string {
	return l.extension
}
