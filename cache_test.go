// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/nanoninja/assert"
)

var (
	_ Renderer = (*cacheRenderer)(nil)
	_ Renderer = Cache(JSON())
	_ Renderer = NewCache(JSON(), CacheConfig{})
)

func TestCacheRenderer(t *testing.T) {
	t.Run("RendersContentType", func(t *testing.T) {
		assert.Equal(t, Cache(JSON()).ContentType(), ContentTypeJSON)
	})

	t.Run("RendersCachedResult", func(t *testing.T) {
		var w1, w2 bytes.Buffer

		data := map[string]string{"key": "value"}
		r := Cache(JSON())

		assert.NoError(t, r.Render(context.Background(), &w1, data, NoOptions))
		assert.NoError(t, r.Render(context.Background(), &w2, data, NoOptions))
		assert.Equal(t, w2.String(), w1.String())
	})

	t.Run("RendersOnlyOnce", func(t *testing.T) {
		calls := 0
		r := NewCache(&callCountRenderer{onRender: func() { calls++ }}, CacheConfig{})

		var w bytes.Buffer
		assert.NoError(t, r.Render(context.Background(), &w, nil, NoOptions))
		assert.NoError(t, r.Render(context.Background(), &w, nil, NoOptions))
		assert.NoError(t, r.Render(context.Background(), &w, nil, NoOptions))
		assert.Equal(t, 1, calls)
	})

	t.Run("PermanentCacheNeverExpires", func(t *testing.T) {
		calls := 0
		r := Cache(&callCountRenderer{onRender: func() { calls++ }})

		var w bytes.Buffer
		assert.NoError(t, r.Render(context.Background(), &w, nil, NoOptions))
		assert.NoError(t, r.Render(context.Background(), &w, nil, NoOptions))
		assert.Equal(t, 1, calls)
	})

	t.Run("RefreshesAfterTTLExpires", func(t *testing.T) {
		calls := 0
		r := NewCache(&callCountRenderer{onRender: func() { calls++ }}, CacheConfig{
			TTL: 10 * time.Millisecond,
		})

		var w bytes.Buffer
		assert.NoError(t, r.Render(context.Background(), &w, nil, NoOptions))
		assert.Equal(t, 1, calls)

		time.Sleep(20 * time.Millisecond)

		assert.NoError(t, r.Render(context.Background(), &w, nil, NoOptions))
		assert.Equal(t, 2, calls)
	})

	t.Run("DoesNotRefreshBeforeTTLExpires", func(t *testing.T) {
		calls := 0
		r := NewCache(&callCountRenderer{onRender: func() { calls++ }}, CacheConfig{
			TTL: 1 * time.Hour,
		})

		var w bytes.Buffer
		assert.NoError(t, r.Render(context.Background(), &w, nil, NoOptions))
		assert.NoError(t, r.Render(context.Background(), &w, nil, NoOptions))
		assert.Equal(t, 1, calls)
	})

	t.Run("HandlesConcurrentAccess", func(_ *testing.T) {
		data := map[string]string{"key": "value"}
		r := Cache(JSON())

		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var w bytes.Buffer
				_ = r.Render(context.Background(), &w, data, NoOptions)
			}()
		}
		wg.Wait()
	})

	t.Run("RespectsContextCancellation", func(t *testing.T) {
		var w bytes.Buffer

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := Cache(JSON()).Render(ctx, &w, nil, NoOptions)

		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("ReturnsErrorOnRendererFailure", func(t *testing.T) {
		var w bytes.Buffer
		err := Cache(JSON()).Render(context.Background(), &w, make(chan int), NoOptions)
		assert.Error(t, err)
	})

	t.Run("ReturnsErrorOnWriteFailure", func(t *testing.T) {
		err := Cache(Text()).Render(context.Background(), &errorWriterTest{}, "data", NoOptions)
		assert.Error(t, err)
	})

	t.Run("WritesFromCacheToWriter", func(t *testing.T) {
		var w bytes.Buffer
		r := Cache(Text())
		assert.NoError(t, r.Render(context.Background(), &w, "data", NoOptions))

		err := r.Render(context.Background(), &errorWriterTest{}, "data", NoOptions)
		assert.Error(t, err)
	})
}

type callCountRenderer struct {
	onRender func()
}

func (r *callCountRenderer) ContentType() string { return "" }

func (r *callCountRenderer) Render(_ context.Context, _ io.Writer, _ any, _ Options) error {
	r.onRender()
	return nil
}
