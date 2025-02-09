// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package render

import (
	"context"
	"io"
)

// Renderer defines a common interface for all renderers in the system.
// It provides a unified way to render content with or without context,
// using functional options for configuration.
//
// The interface is designed to be simple yet flexible, allowing different
// types of renderers (JSON, XML, HTML, etc.) to share the same API while
// implementing their specific rendering logic.
type Renderer interface {
	// Render writes the formatted representation of data to the writer.
	// It uses the provided options to configure the rendering process.
	// This is a convenience method that uses a background context.
	Render(w io.Writer, data any, opts ...func(*Options)) error

	// RenderContext is similar to Render but accepts a context for
	// cancellation, timeouts, and other context-based features.
	// The context allows controlling the rendering lifecycle.
	RenderContext(ctx context.Context, w io.Writer, data any, opts ...func(*Options)) error
}

// MimeTextPlain provides default text/plain content type options with UTF-8 encoding.
// Can be overridden using options in Render/RenderContext methods.
func MimeTextPlain() func(*Options) {
	return MimeUTF8("text/plain")
}

// MimeTextHTML provides default text/html content type options with UTF-8 encoding.
// Can be overridden using options in Render/RenderContext methods.
func MimeTextHTML() func(*Options) {
	return MimeUTF8("text/html")
}

// MimeJSON provides default application/json content type options with UTF-8 encoding.
// Can be overridden using options in Render/RenderContext methods.
func MimeJSON() func(*Options) {
	return MimeUTF8("application/json")
}

// MimeXML provides default application/xml content type options with UTF-8 encoding.
// Can be overridden using options in Render/RenderContext methods.
func MimeXML() func(*Options) {
	return MimeUTF8("application/xml")
}

// MimeCSV provides default text/csv content type options with UTF-8 encoding.
// Can be overridden using options in Render/RenderContext methods.
func MimeCSV() func(*Options) {
	return MimeUTF8("text/csv")
}

// MimeBinary provides default application/octet-stream content type options.
// Used for binary data or when the content type is unknown.
func MimeBinary() func(*Options) {
	return Mime("application/octet-stream")
}

// MimePDF provides default application/pdf content type options.
func MimePDF() func(*Options) {
	return Mime("application/pdf")
}

// MimeStream provides default application/octet-stream content type options.
// Used for streaming binary data.
func MimeStream() func(*Options) {
	return Mime("application/octet-stream")
}

// MimeYAML provides default application/yaml content type options with UTF-8 encoding.
// Can be overridden using options in Render/RenderContext methods.
func MimeYAML() func(*Options) {
	return MimeUTF8("application/yaml")
}

// CheckContext verifies if the context is still valid.
// It returns nil if the context is valid, or the context error if it's done.
func CheckContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
