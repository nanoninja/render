// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"context"
	"fmt"
	"io"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// MarkdownConfig defines configuration for the Markdown renderer.
type MarkdownConfig struct {
	// GFM enables GitHub Flavored Markdown extensions:
	// tables, strikethrough, autolinks, and task lists.
	GFM bool

	// Unsafe allows raw HTML inside Markdown content.
	// Disabled by default to prevent XSS vulnerabilities.
	Unsafe bool
}

// markdownRenderer converts Markdown content to HTML using goldmark.
type markdownRenderer struct {
	config MarkdownConfig
	md     goldmark.Markdown
}

// Markdown creates a new Markdown renderer with default configuration.
// It converts Markdown input to HTML and sets the Content-Type to text/html.
func Markdown() Renderer {
	return NewMarkdown(MarkdownConfig{})
}

// NewMarkdown creates a Markdown renderer with custom configuration.
// Use MarkdownConfig to enable GFM extensions or allow raw HTML.
func NewMarkdown(c MarkdownConfig) Renderer {
	var gopts []goldmark.Option

	if c.GFM {
		gopts = append(gopts, goldmark.WithExtensions(extension.GFM))
	}
	if c.Unsafe {
		gopts = append(gopts, goldmark.WithRendererOptions(html.WithUnsafe()))
	}

	return &markdownRenderer{
		config: c,
		md:     goldmark.New(gopts...),
	}
}

func (r *markdownRenderer) ContentType() string { return ContentTypeHTML }

// Render converts Markdown data to HTML with context support.
// It accepts string or []byte as data.
// The Content-Type is set to text/html automatically.
func (r *markdownRenderer) Render(ctx context.Context, w io.Writer, data any, opts Options) error {
	if err := CheckContext(ctx); err != nil {
		return err
	}

	_, cancel := ApplyTimeout(ctx, opts)
	defer cancel()

	var src []byte

	switch v := data.(type) {
	case string:
		src = []byte(v)
	case []byte:
		src = v
	default:
		return fmt.Errorf("markdown: unsupported data type %T", data)
	}

	if err := r.md.Convert(src, w); err != nil {
		return fmt.Errorf("markdown: %w", err)
	}

	return nil
}
