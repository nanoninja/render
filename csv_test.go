// Copyright 2025 The Nanoninja Authors. All rights reserved.
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

		err := CSV().Render(&w, data)

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "name,age\nAlice,25\nBob,30\n")
	})

	t.Run("HandlesCustomDelimiter", func(t *testing.T) {
		var w bytes.Buffer

		data := [][]string{
			{"name", "age"},
			{"Alice", "25"},
		}

		err := CSV().Render(&w, data, Separator(";"))

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "name;age\nAlice;25\n")
	})

	t.Run("HandlesEmptyData", func(t *testing.T) {
		var w bytes.Buffer

		err := CSV().Render(&w, [][]string{})

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "")
	})

	t.Run("ReturnsErrorForInvalidDataType", func(t *testing.T) {
		var w bytes.Buffer
		data := "invalid data"

		err := CSV().Render(&w, data)

		assert.ErrorIs(t, err, ErrInvalidData)
	})

	t.Run("SetsCorrectContentType", func(t *testing.T) {
		var w bytes.Buffer
		data := [][]string{{"test"}}

		var opts *Options
		err := CSV().Render(&w, data, CaptureOptions(&opts))

		assert.Nil(t, err)
		assert.Equals(t, opts.ContentType(), "text/csv; charset=utf-8")
	})

	t.Run("RespectsCRLFFormatOption", func(t *testing.T) {
		var w bytes.Buffer

		data := [][]string{
			{"name", "age"},
			{"Alice", "25"},
		}

		err := CSV().Render(&w, data, UseCRLF())

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "name,age\r\nAlice,25\r\n")
	})

	t.Run("RespectsContextCancellation", func(t *testing.T) {
		var w bytes.Buffer
		data := [][]string{{"test"}}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := CSV().RenderContext(ctx, &w, data)

		assert.ErrorIs(t, err, context.Canceled)
	})
}
