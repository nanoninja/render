// Copyright 2025 The Nanoninja Authors. All rights reserved.
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

// Render writes the JSON representation of data to the writer.
// It uses a background context and forwards to RenderContext.
func (r jsonRenderer) Render(w io.Writer, data any, opts ...func(*Options)) error {
	return r.RenderContext(context.Background(), w, data, opts...)
}

// RenderContext writes the JSON representation of data with context support.
// It handles:
// - JSONP wrapping if padding is configured
// - HTML escaping based on configuration
// - Pretty printing with customizable indent and prefix
// - Content type setting to application/json
func (r *jsonRenderer) RenderContext(ctx context.Context, w io.Writer, data any, opts ...func(*Options)) error {
	if err := CheckContext(ctx); err != nil {
		return err
	}
	options := NewOptions().
		Use(MimeJSON()).
		Use(opts...)

	if r.config.Padding != "" {
		if _, err := fmt.Fprintf(w, "%s(", r.config.Padding); err != nil {
			return err
		}
		defer fmt.Fprintf(w, ")")
	}
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(r.config.EscapeHTML)

	if options.format.pretty {
		prefix := options.format.prefix
		indent := r.config.Indent

		if prefix == "" {
			prefix = r.config.Prefix
		}
		if options.format.indent != "" {
			indent = options.format.indent
		}
		encoder.SetIndent(prefix, indent)
	}
	return encoder.Encode(data)
}
