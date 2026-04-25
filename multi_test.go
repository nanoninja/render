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
	_ Renderer = (*multiRenderer)(nil)
	_ Renderer = Multi(JSON())
)

func TestMultiRenderer(t *testing.T) {
	t.Run("ContentType", func(t *testing.T) {
		assert.Equal(t, ContentTypeJSON, Multi(JSON()).ContentType())
	})

	t.Run("WritesToPrimaryWriter", func(t *testing.T) {
		var w bytes.Buffer

		err := Multi(JSON()).Render(context.Background(), &w, map[string]string{"key": "value"}, NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, w.String(), `"key"`)
	})

	t.Run("WritesToAllWriters", func(t *testing.T) {
		var w1, w2, w3 bytes.Buffer

		err := Multi(JSON(), &w2, &w3).Render(context.Background(), &w1, map[string]string{"key": "value"}, NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, w2.String(), w1.String())
		assert.Equal(t, w3.String(), w1.String())
	})

	t.Run("WorksWithNoAdditionalWriters", func(t *testing.T) {
		var w bytes.Buffer

		err := Multi(JSON()).Render(context.Background(), &w, map[string]string{"msg": "hello"}, NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, w.String(), "hello")
	})

	t.Run("PropagatesRendererError", func(t *testing.T) {
		var w bytes.Buffer

		err := Multi(JSON()).Render(context.Background(), &w, make(chan int), NoOptions)

		assert.Error(t, err)
		assert.StringContains(t, err.Error(), "multi:")
	})

	t.Run("RespectsContextCancellation", func(t *testing.T) {
		var w bytes.Buffer

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := Multi(JSON()).Render(ctx, &w, map[string]string{"key": "value"}, NoOptions)

		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("RespectsTimeout", func(t *testing.T) {
		var w bytes.Buffer

		err := Multi(JSON()).Render(context.Background(), &w, map[string]string{"key": "value"}, Options{Timeout: 1 * time.Hour})

		assert.NoError(t, err)
	})
}
