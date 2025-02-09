// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package loader

import (
	"errors"
	"fmt"

	"github.com/nanoninja/render/tmpl"
)

// CompositeLoader combines multiple loaders with priority order.
// When looking for templates, loaders are queried in order until a match is found.
// This allows for layered template organization, like having default templates
// that can be overridden by custom ones.
type CompositeLoader struct {
	baseLoader
	loaders []tmpl.Loader
}

// NewComposite creates a composite loader that combines multiple loaders.
// Loaders are queried in the order they are provided - first loader
// that has a matching template takes precedence.
//
// Example:
//
//	src := loader.NewComposite(
//	    customLoader,  // Custom templates have priority
//	    defaultLoader, // Fallback to default templates
//	    tmpl.LoaderConfig{
//	        Extension: ".html",
//	    },
//	)
func NewComposite(loaders []tmpl.Loader, config tmpl.LoaderConfig) tmpl.Loader {
	return &CompositeLoader{
		baseLoader: newBase(config, false),
		loaders:    loaders,
	}
}

// Load returns all template names from all configured loaders.
// If multiple loaders have templates with the same name, they are
// de-duplicated with priority given to earlier loaders.
func (l *CompositeLoader) Load(pattern string) ([]string, error) {
	seen := make(map[string]bool)
	var templates []string

	// Query each loader in order
	for _, loader := range l.loaders {
		names, err := loader.Load(pattern)
		if err != nil {
			return nil, fmt.Errorf("composite load error: %w", err)
		}
		// Add only unseen templates to preserve priority order
		for _, name := range names {
			if !seen[name] {
				seen[name] = true
				templates = append(templates, name)
			}
		}
	}
	return templates, nil
}

// Read returns the content of the named template from the first loader
// that contains it. This respects the priority order of the loaders.
func (l *CompositeLoader) Read(name string) ([]byte, error) {
	var lastErr error

	// Try each loader in order
	for _, loader := range l.loaders {
		content, err := loader.Read(name)
		if err == nil {
			return content, nil
		}
		lastErr = err
	}
	// If no loader found the template, return the last error
	if errors.Is(lastErr, tmpl.ErrTemplateNotFound) {
		return nil, fmt.Errorf("%w in any loader: %s", tmpl.ErrTemplateNotFound, name)
	}
	return nil, fmt.Errorf("composite read error: %w", lastErr)
}

// Extension returns the configured template file extension.
func (l *CompositeLoader) Extension() string {
	return l.extension
}
