// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"context"
	"testing"

	"github.com/nanoninja/assert"
)

var (
	_ Renderer = (*csvRenderer)(nil)
	_ Renderer = CSV()
)

func TestCSVRenderer(t *testing.T) {
	t.Run("RendersSimpleStringArrayData", func(t *testing.T) {
		var w bytes.Buffer

		data := [][]string{
			{"name", "age"},
			{"Alice", "25"},
			{"Bob", "30"},
		}

		err := CSV().Render(context.Background(), &w, data, NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "name,age\nAlice,25\nBob,30\n", w.String())
	})

	t.Run("HandlesCustomDelimiter", func(t *testing.T) {
		var w bytes.Buffer

		data := [][]string{
			{"name", "age"},
			{"Alice", "25"},
		}

		err := CSV().Render(context.Background(), &w, data, Options{Separator: ';'})

		assert.NoError(t, err)
		assert.Equal(t, "name;age\nAlice;25\n", w.String())
	})

	t.Run("HandlesEmptyData", func(t *testing.T) {
		var w bytes.Buffer

		err := CSV().Render(context.Background(), &w, [][]string{}, NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "", w.String())
	})

	t.Run("ReturnsErrorForInvalidDataType", func(t *testing.T) {
		var w bytes.Buffer
		data := "invalid data"

		err := CSV().Render(context.Background(), &w, data, NoOptions)

		assert.ErrorIs(t, err, ErrInvalidData)
	})

	t.Run("SetsCorrectContentType", func(t *testing.T) {
		assert.Equal(t, "text/csv; charset=utf-8", CSV().ContentType())
	})

	t.Run("RespectsCRLFFormatOption", func(t *testing.T) {
		var w bytes.Buffer

		data := [][]string{
			{"name", "age"},
			{"Alice", "25"},
		}

		err := CSV().Render(context.Background(), &w, data, Options{LineEnding: "\r\n"})

		assert.NoError(t, err)
		assert.Equal(t, "name,age\r\nAlice,25\r\n", w.String())
	})

	t.Run("RespectsContextCancellation", func(t *testing.T) {
		var w bytes.Buffer
		data := [][]string{{"test"}}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := CSV().Render(ctx, &w, data, NoOptions)

		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("ReturnsErrorOnWriteFailure", func(t *testing.T) {
		data := make([][]string, 100)
		for i := range data {
			data[i] = []string{"name", "age", "city", "country"}
		}

		err := CSV().Render(context.Background(), &errorWriterTest{}, data, NoOptions)

		assert.Error(t, err)
		assert.StringContains(t, err.Error(), "csv:")
	})
}
