// Copyright 2026 The Nanoninja Authors. All rights reserved.
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

func (r xmlRenderer) ContentType() string {
	return ContentTypeXML
}

func (r *xmlRenderer) Render(ctx context.Context, w io.Writer, data any, opts Options) error {
	if err := CheckContext(ctx); err != nil {
		return err
	}
	if r.config.Header {
		if _, err := io.WriteString(w, xml.Header); err != nil {
			return err
		}
	}
	encoder := xml.NewEncoder(w)

	if opts.Pretty {
		prefix := r.config.Prefix
		indent := r.config.Indent

		if opts.Prefix != "" {
			prefix = opts.Prefix
		}
		if opts.Indent != "" {
			indent = opts.Indent
		}
		encoder.Indent(prefix, indent)
	}
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("xml: %w", err)
	}
	return nil
}
