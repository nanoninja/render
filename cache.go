// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"
)

// CacheConfig holds configuration for the cache renderer.
type CacheConfig struct {
	TTL time.Duration // Duration before the cache expires. 0 means permanent.
}

type cacheRenderer struct {
	renderer Renderer
	ttl      time.Duration

	mu        sync.RWMutex
	cached    []byte
	expiresAt time.Time
	hasCache  bool
}

// Cache creates a renderer that caches the output permanently.
// Use it for truly static data that never changes.
//
// Example:
//
//	var cached = render.Cache(render.JSON())
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	    cached.Render(w, data, render.WriteResponse(w))
//	}
func Cache(r Renderer) Renderer {
	return NewCache(r, CacheConfig{})
}

// NewCache creates a renderer that caches the output with optional TTL.
// If TTL is zero, the cache never expires.
//
// Example:
//
//	renderer := render.NewCache(render.JSON(), render.CacheConfig{
//	    TTL: 5 * time.Minute,
//	})
func NewCache(r Renderer, cfg CacheConfig) Renderer {
	return &cacheRenderer{
		renderer: r,
		ttl:      cfg.TTL,
	}
}

func (r *cacheRenderer) ContentType() string {
	return r.renderer.ContentType()
}

func (r *cacheRenderer) Render(ctx context.Context, w io.Writer, data any, opts Options) error {
	if err := CheckContext(ctx); err != nil {
		return err
	}
	_, cancel := ApplyTimeout(ctx, opts)
	defer cancel()

	// Fast path — cache hit.
	r.mu.RLock()
	if r.hasCache && (r.ttl == 0 || time.Now().Before(r.expiresAt)) {
		cached := r.cached
		r.mu.RUnlock()
		if _, err := w.Write(cached); err != nil {
			return fmt.Errorf("cache: %w", err)
		}
		return nil
	}
	r.mu.RUnlock()

	// Slow path — cache miss or expired.
	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-check after acquiring write lock — another goroutine may have
	// populated the cache between RUnlock and Lock.
	if r.hasCache && (r.ttl == 0 || time.Now().Before(r.expiresAt)) {
		if _, err := w.Write(r.cached); err != nil {
			return fmt.Errorf("cache: %w", err)
		}
		return nil
	}

	var buf bytes.Buffer
	if err := r.renderer.Render(ctx, &buf, data, opts); err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	r.cached = buf.Bytes()
	r.hasCache = true
	if r.ttl > 0 {
		r.expiresAt = time.Now().Add(r.ttl)
	}
	if _, err := w.Write(r.cached); err != nil {
		return fmt.Errorf("cache: %w", err)
	}
	return nil
}
