// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nanoninja/assert"
)

var (
	_ Renderer = (*textRenderer)(nil)
	_ Renderer = Text()
)

func TestTextRenderer(t *testing.T) {
	t.Run("RendersStringData", func(t *testing.T) {
		var w bytes.Buffer

		err := Text().Render(context.Background(), &w, "Hello world!", NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "Hello world!", w.String())
	})

	t.Run("RendersWithFormatArguments", func(t *testing.T) {
		var w bytes.Buffer

		err := Text().Render(context.Background(), &w, "Hello, %s", Options{Args: []any{"Gophers"}})

		assert.NoError(t, err)
		assert.Equal(t, "Hello, Gophers", w.String())
	})

	t.Run("RenderStringLiteralWithoutFormatInterpolation", func(t *testing.T) {
		var w bytes.Buffer

		err := Text().Render(context.Background(), &w, "100% done", NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "100% done", w.String())
	})

	t.Run("RendersStringerTypes", func(t *testing.T) {
		var w bytes.Buffer
		data := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

		err := Text().Render(context.Background(), &w, data, NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "2026-01-01 00:00:00 +0000 UTC", w.String())
	})

	t.Run("RendersErrorTypes", func(t *testing.T) {
		var w bytes.Buffer
		data := errors.New("test error")

		err := Text().Render(context.Background(), &w, data, NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "test error", w.String())
	})

	t.Run("RendersDefaultType", func(t *testing.T) {
		var w bytes.Buffer
		data := 20

		err := Text().Render(context.Background(), &w, data, NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "20", w.String())
	})

	t.Run("RendersWithPrettyPrinting", func(t *testing.T) {
		var w bytes.Buffer

		err := Text().Render(context.Background(), &w, "hello", Options{Pretty: true})

		assert.NoError(t, err)
		assert.Equal(t, "hello\n", w.String())
	})

	t.Run("RespectsContextCancellation", func(t *testing.T) {
		var w bytes.Buffer

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := Text().Render(ctx, &w, "test", NoOptions)

		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("SetsDefaultContentType", func(t *testing.T) {
		assert.Equal(t, "text/plain; charset=utf-8", Text().ContentType())
	})
}
