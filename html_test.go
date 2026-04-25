// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/nanoninja/assert"
)

var (
	_ Renderer = (*htmlRenderer)(nil)
	_ Renderer = HTML()
)

func TestHTMLRenderer(t *testing.T) {
	t.Run("RendersStringData", func(t *testing.T) {
		var w bytes.Buffer

		err := HTML().Render(context.Background(), &w, "<h1>Hello</h1>", NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "<h1>Hello</h1>", w.String())
	})

	t.Run("RendersByteSlice", func(t *testing.T) {
		var w bytes.Buffer

		err := HTML().Render(context.Background(), &w, []byte("<p>Hello</p>"), NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "<p>Hello</p>", w.String())
	})

	t.Run("RendersDefaultType", func(t *testing.T) {
		var w bytes.Buffer

		err := HTML().Render(context.Background(), &w, 42, NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "42", w.String())
	})

	t.Run("RendersStringerType", func(t *testing.T) {
		var w bytes.Buffer
		data := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

		err := HTML().Render(context.Background(), &w, data, NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "2026-01-01 00:00:00 +0000 UTC", w.String())
	})

	t.Run("SetsDefaultContentType", func(t *testing.T) {
		assert.Equal(t, "text/html; charset=utf-8", HTML().ContentType())
	})

	t.Run("RespectsContextCancellation", func(t *testing.T) {
		var w bytes.Buffer

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := HTML().Render(ctx, &w, "<p>test</p>", NoOptions)

		assert.ErrorIs(t, err, context.Canceled)
	})
}
