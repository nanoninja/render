// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"context"
	"io"
	"time"
)

// ContentTypes groups the MIME type constants used by the built-in renderers.
const (
	ContentTypeJSON   = "application/json; charset=utf-8"
	ContentTypeXML    = "application/xml; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeCSV    = "text/csv; charset=utf-8"
	ContentTypeBinary = "application/octet-stream"
	ContentTypeYAML   = "application/yaml; charset=utf-8"
)

// Renderer defines a common interface for all renderers in the system.
// It provides a unified way to render content with or without context,
// using functional options for configuration.
//
// The interface is designed to be simple yet flexible, allowing different
// types of renderers (JSON, XML, HTML, etc.) to share the same API while
// implementing their specific rendering logic.
type Renderer interface {
	ContentType() string

	// Render writes the rendered content to w using the provided context and options.
	// The context is checked for cancellation before rendering begins.
	Render(ctx context.Context, w io.Writer, data any, opts Options) error
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

// WithTimeout returns a copy of opts with Timeout set to d.
func WithTimeout(opts Options, d time.Duration) Options {
	opts.Timeout = d
	return opts
}

// WithPretty returns a copy of opts with Pretty enabled.
func WithPretty(opts Options) Options {
	opts.Pretty = true
	return opts
}

// WithIndent returns a copy of opts with Indent set to s.
func WithIndent(opts Options, s string) Options {
	opts.Indent = s
	return opts
}

// WithHeader returns a copy of opts with the given header key/value added.
func WithHeader(opts Options, key, value string) Options {
	if opts.Headers == nil {
		opts.Headers = make(map[string][]string)
	}
	opts.Headers[key] = append(opts.Headers[key], value)
	return opts
}

// ApplyTimeout returns a derived context with a timeout if one is configured in opts.
// If no timeout is set, it returns the original context with a no-op cancel function.
// The cancel function must always be called to release resources, typically via defer.
func ApplyTimeout(ctx context.Context, opts Options) (context.Context, context.CancelFunc) {
	if opts.Timeout > 0 {
		return context.WithTimeout(ctx, opts.Timeout)
	}
	return ctx, func() {}
}
