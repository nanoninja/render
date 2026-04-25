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
	_ Renderer = (*binaryRenderer)(nil)
	_ Renderer = Binary()
)

func TestBinaryRenderer(t *testing.T) {
	t.Run("RendersByteSlice", func(t *testing.T) {
		var w bytes.Buffer

		data := []byte{0x89, 0x50, 0x4E, 0x47}

		err := Binary().Render(context.Background(), &w, data, NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, data, w.Bytes())
	})

	t.Run("RendersReader", func(t *testing.T) {
		var w bytes.Buffer

		data := strings.NewReader("binary content")

		err := Binary().Render(context.Background(), &w, data, NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "binary content", w.String())
	})

	t.Run("ReturnsErrorForUnsupportedType", func(t *testing.T) {
		var w bytes.Buffer

		err := Binary().Render(context.Background(), &w, 42, NoOptions)

		assert.Error(t, err)
		assert.StringContains(t, err.Error(), "binary: unsupported data type")
	})

	t.Run("SetsDefaultContentType", func(t *testing.T) {
		assert.Equal(t, "application/octet-stream", Binary().ContentType())
	})

	t.Run("RespectsContextCancellation", func(t *testing.T) {
		var w bytes.Buffer

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := Binary().Render(ctx, &w, []byte("data"), NoOptions)

		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("RespectsTimeout", func(t *testing.T) {
		var w bytes.Buffer

		err := Binary().Render(context.Background(), &w, []byte("data"), Options{Timeout: 1 * time.Hour})

		assert.NoError(t, err)
	})
}
