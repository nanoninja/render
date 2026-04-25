// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/nanoninja/assert"
)

type gzipCloseErrorWriter struct {
	headerWritten bool
}

func (w *gzipCloseErrorWriter) Write(p []byte) (int, error) {
	if !w.headerWritten {
		w.headerWritten = true
		return len(p), nil
	}
	return 0, errors.New("write error")
}

var (
	_ Renderer = (*gzipRenderer)(nil)
	_ Renderer = Gzip()
)

func TestGzipRenderer(t *testing.T) {
	t.Run("ContentType", func(t *testing.T) {
		assert.Equal(t, ContentTypeBinary, Gzip().ContentType())
	})

	t.Run("CompressByteSlice", func(t *testing.T) {
		var w bytes.Buffer

		data := []byte("hello gzip")

		err := Gzip().Render(context.Background(), &w, data, NoOptions)

		assert.NoError(t, err)

		r, err := gzip.NewReader(&w)
		assert.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, r.Close()) })

		got, err := io.ReadAll(r)
		assert.NoError(t, err)
		assert.Equal(t, data, got)
	})

	t.Run("ReturnsErrorForNonByteSlice", func(t *testing.T) {
		var w bytes.Buffer

		err := Gzip().Render(context.Background(), &w, "not a []byte", NoOptions)

		assert.Error(t, err)
		assert.StringContains(t, err.Error(), "gzip:")
	})

	t.Run("RespectsContextCancellation", func(t *testing.T) {
		var w bytes.Buffer

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := Gzip().Render(ctx, &w, []byte("data"), NoOptions)

		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("RespectsTimeout", func(t *testing.T) {
		var w bytes.Buffer

		err := Gzip().Render(context.Background(), &w, []byte("data"), Options{Timeout: 1 * time.Hour})

		assert.NoError(t, err)
	})

	t.Run("ReturnsErrorOnWriteFailure", func(t *testing.T) {
		err := Gzip().Render(context.Background(), &errorWriterTest{}, []byte("data"), NoOptions)

		assert.Error(t, err)
		assert.StringContains(t, err.Error(), "gzip:")
	})

	t.Run("ReturnsErrorOnCloseFailure", func(t *testing.T) {
		err := Gzip().Render(context.Background(), &gzipCloseErrorWriter{}, []byte("data"), NoOptions)

		assert.Error(t, err)
		assert.StringContains(t, err.Error(), "gzip:")
	})
}
