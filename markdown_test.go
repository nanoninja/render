// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/nanoninja/assert"
)

var (
	_ Renderer = (*markdownRenderer)(nil)
	_ Renderer = Markdown()
	_ Renderer = NewMarkdown(MarkdownConfig{})
)

func TestMarkdownRenderer(t *testing.T) {
	t.Run("RendersHeading", func(t *testing.T) {
		var w bytes.Buffer

		err := Markdown().Render(context.Background(), &w, "# Hello", NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, w.String(), "<h1>Hello</h1>")
	})

	t.Run("RendersParagraph", func(t *testing.T) {
		var w bytes.Buffer

		err := Markdown().Render(context.Background(), &w, "Hello World", NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, w.String(), "<p>Hello World</p>")
	})

	t.Run("RendersByteSlice", func(t *testing.T) {
		var w bytes.Buffer

		err := Markdown().Render(context.Background(), &w, []byte("**bold**"), NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, w.String(), "<strong>bold</strong>")
	})

	t.Run("RendersLink", func(t *testing.T) {
		var w bytes.Buffer

		err := Markdown().Render(context.Background(), &w, "[Go](https://go.dev)", NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, w.String(), "<a href=")
		assert.StringContains(t, w.String(), "Go</a>")
	})

	t.Run("OutputEndsWithNewline", func(t *testing.T) {
		var w bytes.Buffer

		err := Markdown().Render(context.Background(), &w, "hello", NoOptions)

		assert.NoError(t, err)
		assert.True(t, strings.HasSuffix(w.String(), "\n"))
	})

	t.Run("SetsDefaultContentType", func(t *testing.T) {
		assert.Equal(t, "text/html; charset=utf-8", Markdown().ContentType())
	})

	t.Run("ReturnsErrorForUnsupportedType", func(t *testing.T) {
		var w bytes.Buffer

		err := Markdown().Render(context.Background(), &w, 42, NoOptions)

		assert.Error(t, err)
		assert.StringContains(t, err.Error(), "markdown: unsupported data type")
	})

	t.Run("RespectsContextCancellation", func(t *testing.T) {
		var w bytes.Buffer

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := Markdown().Render(ctx, &w, "# Hello", NoOptions)

		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("RespectsTimeout", func(t *testing.T) {
		var w bytes.Buffer

		err := Markdown().Render(context.Background(), &w, "# Hello", Options{Timeout: 1 * time.Hour})

		assert.NoError(t, err)
	})

	t.Run("ReturnsErrorOnWriteFailure", func(t *testing.T) {
		err := Markdown().Render(context.Background(), &errorWriterTest{}, "# Hello", NoOptions)

		assert.Error(t, err)
		assert.StringContains(t, err.Error(), "markdown:")
	})

	t.Run("GFMRendersTable", func(t *testing.T) {
		var w bytes.Buffer

		md := "| A | B |\n|---|---|\n| 1 | 2 |"
		err := NewMarkdown(MarkdownConfig{GFM: true}).Render(context.Background(), &w, md, NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, w.String(), "<table>")
	})

	t.Run("GFMRendersStrikethrough", func(t *testing.T) {
		var w bytes.Buffer

		err := NewMarkdown(MarkdownConfig{GFM: true}).Render(context.Background(), &w, "~~deleted~~", NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, w.String(), "<del>deleted</del>")
	})

	t.Run("UnsafeAllowsRawHTML", func(t *testing.T) {
		var w bytes.Buffer

		err := NewMarkdown(MarkdownConfig{Unsafe: true}).Render(context.Background(), &w, "<div>raw</div>", NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, w.String(), "<div>raw</div>")
	})
}
