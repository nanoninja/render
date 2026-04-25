// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// JSONConfig defines configuration for JSON renderer.
// It provides JSON-specific settings that are set during initialization.
type JSONConfig struct {
	// Controls HTML character escaping
	EscapeHTML bool

	// Custom prefix for JSON output
	Prefix string

	// Custom indentation for JSON output
	Indent string

	// JSONP function name (if empty, standard JSON is used)
	Padding string
}

// jsonRenderer implements JSON and JSONP rendering with configurable formatting options.
type jsonRenderer struct {
	config JSONConfig
}

func (r *jsonRenderer) ContentType() string {
	return ContentTypeJSON
}

// JSON creates a new JSONRenderer with safe default configuration:
// - HTML escaping enabled for web safety
// - Standard 2-space indentation
// This is the recommended constructor for most use cases.
func JSON() Renderer {
	return NewJSON(JSONConfig{
		EscapeHTML: true, // Safe default for web contexts
		Indent:     "  ", // Standard 2-space indentation
	})
}

// NewJSON creates a JSONRenderer with custom configuration.
// Use this when you need specific JSON behaviors different from defaults.
func NewJSON(c JSONConfig) Renderer {
	return &jsonRenderer{config: c}
}

func (r *jsonRenderer) Render(ctx context.Context, w io.Writer, data any, opts Options) error {
	if err := CheckContext(ctx); err != nil {
		return err
	}
	if r.config.Padding != "" {
		if _, err := io.WriteString(w, r.config.Padding); err != nil {
			return err
		}
		if _, err := io.WriteString(w, "("); err != nil {
			return err
		}
		defer io.WriteString(w, ")") //nolint:errcheck
	}
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(r.config.EscapeHTML)

	if opts.Pretty {
		prefix := opts.Prefix
		indent := r.config.Indent

		if prefix == "" {
			prefix = r.config.Prefix
		}
		if opts.Indent != "" {
			indent = opts.Indent
		}
		encoder.SetIndent(prefix, indent)
	}
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("json: %w", err)
	}
	return nil
}
