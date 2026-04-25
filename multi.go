// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"context"
	"fmt"
	"io"
)

// multiRenderer writes rendered output to multiple writers simultaneously.
type multiRenderer struct {
	renderer Renderer
	writers  []io.Writer
}

// Multi creates a renderer that writes to the primary writer and all additional
// writers simultaneously. It wraps an existing renderer and duplicates its output
// to every destination.
//
// Example:
//
//	// Write JSON to the response and a log buffer at the same time
//	var log bytes.Buffer
//	render.Multi(render.JSON(), &log).Render(w, data)
func Multi(renderer Renderer, writers ...io.Writer) Renderer {
	return &multiRenderer{
		renderer: renderer,
		writers:  writers,
	}
}

func (r *multiRenderer) ContentType() string {
	return r.renderer.ContentType()
}

// Render writes the rendered output to all writers with context support.
// The primary writer w is always included alongside any additional writers
// configured via Multi.
func (r *multiRenderer) Render(ctx context.Context, w io.Writer, data any, opts Options) error {
	if err := CheckContext(ctx); err != nil {
		return err
	}
	ctx, cancel := ApplyTimeout(ctx, opts)
	defer cancel()

	writer := io.MultiWriter(append([]io.Writer{w}, r.writers...)...)

	if err := r.renderer.Render(ctx, writer, data, opts); err != nil {
		return fmt.Errorf("multi: %w", err)
	}
	return nil
}
