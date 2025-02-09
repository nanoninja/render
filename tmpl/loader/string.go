// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package loader

import (
	"fmt"

	"github.com/nanoninja/render/tmpl"
)

// StringLoader implements tmpl.Loader for templates defined directly in code as strings.
// This allows templates to be defined programmatically without needing external files.
type StringLoader struct {
	baseLoader
	templates map[string]string
}

// NewString creates a StringLoader with the given templates map and configuration.
// The templates map keys are used as template names and values are template content.
//
// Example:
//
//	src := loader.NewString(map[string]string{
//	    "welcome.html": "Hello {{ .Name }}",
//	    "layout.html": "<body>{{ template \"content\" . }}</body>",
//	}, tmpl.LoaderConfig{
//	    Extension: ".html",
//	})
func NewString(templates map[string]string, config tmpl.LoaderConfig) tmpl.Loader {
	return &StringLoader{
		baseLoader: newBase(config, false), // Using OS separator
		templates:  templates,
	}
}

// Load returns all template names that match the given pattern.
// Since templates are stored in memory, pattern matching is done against
// template names rather than filesystem paths.
func (l *StringLoader) Load(string) ([]string, error) {
	var templates []string
	for name := range l.templates {
		if !l.hasValidExtension(name) {
			continue
		}
		templates = append(templates, name)
	}
	return templates, nil
}

// Read returns the content of the named template.
// The name should be one of the keys provided in the templates map.
func (l *StringLoader) Read(name string) ([]byte, error) {
	content, ok := l.templates[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", tmpl.ErrTemplateNotFound, name)
	}
	return []byte(content), nil
}

// Extension returns the configured template file extension.
func (l *StringLoader) Extension() string {
	return l.extension
}
