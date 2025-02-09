// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package render

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
)

// XMLConfig defines configuration options for XML rendering.
// It allows customization of escaping, indentation, and XML header inclusion.
type XMLConfig struct {
	// Prefix specifies the string to prepend at the start of each line.
	// This is particularly useful when embedding XML within another format.
	Prefix string

	// Indent specifies the string used for each level of indentation.
	// Common values are spaces ("  ") or tabs ("\t").
	Indent string

	// Header controls whether to include the XML declaration at the start.
	// When true, adds <?xml version="1.0" encoding="UTF-8"?>.
	Header bool
}

// xmlRenderer implements XML rendering with configurable formatting options.
// It supports both compact and pretty-printed output formats.
type xmlRenderer struct {
	config XMLConfig
}

// XML creates a new XMLRenderer with safe default configuration.
// Default settings include:
// - Standard 2-space indentation
// - XML header included
// This is the recommended constructor for most use cases.
func XML() Renderer {
	return NewXML(XMLConfig{
		Indent: "  ", // Standard 2-space indentation
		Header: true, // Include XML header by default
	})
}

// NewXML creates a XMLRenderer with custom configuration.
// Use this when you need specific XML behaviors different from defaults.
func NewXML(c XMLConfig) Renderer {
	return &xmlRenderer{config: c}
}

// Render writes the XML representation of data to the writer.
// It uses a background context and forwards to RenderContext.
// The data must be a value that can be encoded as XML according to
// the encoding/xml package rules.
func (r *xmlRenderer) Render(w io.Writer, data any, opts ...func(*Options)) error {
	return r.RenderContext(context.Background(), w, data, opts...)
}

// RenderContext writes the XML representation of data with context support.
// It handles:
// - XML header inclusion based on configuration
// - Pretty printing with configurable prefix and indentation
// - Content type setting to application/xml
func (r *xmlRenderer) RenderContext(ctx context.Context, w io.Writer, data any, opts ...func(*Options)) error {
	if err := CheckContext(ctx); err != nil {
		return err
	}
	options := NewOptions().
		Use(MimeXML()).
		Use(opts...)

	if r.config.Header {
		if _, err := fmt.Fprint(w, xml.Header); err != nil {
			return err
		}
	}
	encoder := xml.NewEncoder(w)

	if options.format.pretty {
		prefix := r.config.Prefix
		indent := r.config.Indent

		if options.format.prefix != "" {
			prefix = options.format.prefix
		}
		if options.format.indent != "" {
			indent = options.format.indent
		}
		encoder.Indent(prefix, indent)
	}
	return encoder.Encode(data)
}
