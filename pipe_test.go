// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/nanoninja/assert"
)

var (
	_ Renderer = (*pipeRenderer)(nil)
	_ Renderer = Pipe()
)

type upperRenderer struct{}

func (r *upperRenderer) ContentType() string { return "" }

func (r *upperRenderer) Render(_ context.Context, w io.Writer, data any, _ Options) error {
	b, ok := data.([]byte)
	if !ok {
		return nil
	}
	_, err := io.WriteString(w, strings.ToUpper(string(b)))
	return err
}

func TestPipeRenderer(t *testing.T) {
	t.Run("ContentType", func(t *testing.T) {
		assert.Equal(t, "", Pipe().ContentType())
	})

	t.Run("ExecutesSingleRenderer", func(t *testing.T) {
		var w bytes.Buffer

		err := Pipe(Text()).Render(context.Background(), &w, "hello", NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "hello", w.String())
	})

	t.Run("ChainsRenderers", func(t *testing.T) {
		var w bytes.Buffer

		err := Pipe(Text(), &upperRenderer{}).Render(context.Background(), &w, "hello", NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "HELLO", w.String())
	})

	t.Run("PassesJSONOutputToNextRenderer", func(t *testing.T) {
		var w bytes.Buffer

		err := Pipe(JSON(), &upperRenderer{}).Render(context.Background(), &w, map[string]string{"key": "val"}, NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, w.String(), "KEY")
	})

	t.Run("ReturnsNilForEmptyPipe", func(t *testing.T) {
		var w bytes.Buffer

		err := Pipe().Render(context.Background(), &w, "data", NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "", w.String())
	})

	t.Run("PropagatesErrorWithPosition", func(t *testing.T) {
		var w bytes.Buffer

		err := Pipe(JSON(), &upperRenderer{}).Render(context.Background(), &w, make(chan int), NoOptions)

		assert.Error(t, err)
		assert.StringContains(t, err.Error(), "pipe[0]")
	})

	t.Run("PropagatesErrorOnLastRenderer", func(t *testing.T) {
		err := Pipe(&upperRenderer{}, JSON()).Render(context.Background(), &errorWriterTest{}, make(chan int), NoOptions)

		assert.Error(t, err)
		assert.StringContains(t, err.Error(), "pipe[1]")
	})

	t.Run("RespectsContextCancellation", func(t *testing.T) {
		var w bytes.Buffer

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := Pipe(JSON()).Render(ctx, &w, "data", NoOptions)

		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("RespectsTimeout", func(t *testing.T) {
		var w bytes.Buffer

		err := Pipe(JSON()).Render(context.Background(), &w, map[string]string{"key": "value"}, Options{Timeout: 1 * time.Hour})

		assert.NoError(t, err)
	})
}
