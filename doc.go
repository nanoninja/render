// Package render provides a lightweight, allocation-free rendering system for Go.
//
// All renderers implement a single two-method interface:
//
//	type Renderer interface {
//	    ContentType() string
//	    Render(ctx context.Context, w io.Writer, data any, opts Options) error
//	}
//
// Options is a plain value struct — pass it by value, no allocations on the hot path:
//
//	opts := render.Options{Pretty: true, Indent: "  "}
//	err := render.JSON().Render(ctx, w, data, opts)
//
// Use render.NoOptions (the zero value) when no configuration is needed:
//
//	err := render.JSON().Render(ctx, w, data, render.NoOptions)
//
// Built-in renderers: JSON, XML, CSV, Text, HTML, Binary, YAML, Markdown,
// Gzip, Buffer, Cache, Multi, Pipe.
//
// For HTML and text template rendering, see the tmpl sub-package.
package render
