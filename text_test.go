// Copyright 2025 The Nanoninja Authors. All rights reserved.
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

		err := Text().Render(&w, "Hello world!")

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "Hello world!")
	})

	t.Run("RendersWithFormatArguments", func(t *testing.T) {
		var w bytes.Buffer

		err := Text().Render(&w, "Hello, %s", Textf("Gophers"))

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "Hello, Gophers")
	})

	t.Run("RendersStringerTypes", func(t *testing.T) {
		var w bytes.Buffer
		data := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

		err := Text().Render(&w, data)

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "2025-01-01 00:00:00 +0000 UTC")
	})

	t.Run("RendersErrorTypes", func(t *testing.T) {
		var w bytes.Buffer
		data := errors.New("test error")

		err := Text().Render(&w, data)

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "test error")
	})

	t.Run("RendersDefaultType", func(t *testing.T) {
		var w bytes.Buffer
		data := 20

		err := Text().Render(&w, data)

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "20")
	})

	t.Run("RendersWithPrettyPrinting", func(t *testing.T) {
		var w bytes.Buffer

		err := Text().Render(&w, "hello", Format(Pretty()))

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "hello\n")
	})

	t.Run("RespectsContextCancellation", func(t *testing.T) {
		var w bytes.Buffer

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := Text().RenderContext(ctx, &w, "test")

		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("SetsDefaultContentTypeInOptions", func(t *testing.T) {
		var w bytes.Buffer

		var opts *Options
		err := Text().Render(&w, "test", CaptureOptions(&opts))

		assert.Nil(t, err)
		assert.Equals(t, opts.ContentType(), "text/plain; charset=utf-8")
	})
}
