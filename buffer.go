// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package render

import (
	"bytes"
	"context"
	"fmt"
	"io"
)

// BufferConfig holds configuration for the BufferRenderer.
// It allows customization of initial buffer size and post-processing functionality.
type BufferConfig struct {
	// InitialSize specifies the initial size of the buffer in bytes.
	// Setting this to a value close to the expected output size can improve
	// performance by reducing buffer reallocations.
	// If 0, the default buffer size is used.
	InitialSize int

	// PostRender defines a function to transform or validate the rendered content
	// before it is written to the output. This can be used for tasks like:
	// - Content minification
	// - Data compression
	// - Schema validation
	// - Content transformation
	// - Metadata injection
	// If nil, the content is written as-is.
	PostRender func([]byte) ([]byte, error)
}

// BufferRenderer provides buffered rendering capabilities with optional post-processing.
// It wraps another renderer and buffers its output before writing, which allows:
// - Complete content validation before writing
// - Content transformation through post-processing
// - Error handling without partial writes
type BufferRenderer struct {
	renderer    Renderer                     // The underlying renderer to buffer
	initialSize int                          // Initial buffer size if specified
	postRender  func([]byte) ([]byte, error) // Optional post-processing function
}

// Buffer creates a new BufferRenderer with default configuration.
// This is the recommended constructor for simple buffering without post-processing.
//
// Example:
//
//	// Create a buffered JSON renderer
//	renderer := Buffer(JSON())
//	file, _ := os.Create("output.json")
//	renderer.Render(file, data)
func Buffer(r Renderer) *BufferRenderer {
	return NewBuffer(r, BufferConfig{})
}

// NewBuffer creates a BufferRenderer with custom configuration.
// Use this when you need to specify buffer size or add post-processing.
//
// Example:
//
//	// Create a renderer with post-processing
//	renderer := NewBuffer(JSON(), BufferConfig{
//	    InitialSize: 4096,
//	    PostRender: func(data []byte) ([]byte, error) {
//	        // Perform content transformation
//	        return transform(data)
//	    },
//	})
func NewBuffer(r Renderer, c BufferConfig) *BufferRenderer {
	return &BufferRenderer{
		renderer:    r,
		initialSize: c.InitialSize,
		postRender:  c.PostRender,
	}
}

// Render implements buffered rendering using a background context.
// See RenderContext for details on the rendering process.
func (r *BufferRenderer) Render(w io.Writer, data any, opts ...func(*Options)) error {
	return r.RenderContext(context.Background(), w, data, opts...)
}

// RenderContext implements buffered rendering with context support.
// The rendering process follows these steps:
// 1. Creates a buffer (pre-allocated if InitialSize > 0)
// 2. Renders content to the buffer using the wrapped renderer
// 3. If configured, applies post-processing to the buffered content
// 4. Writes the final content to the provided writer
//
// This approach ensures that no partial content is written if an error
// occurs during rendering or post-processing.
func (r *BufferRenderer) RenderContext(ctx context.Context, w io.Writer, data any, opts ...func(*Options)) error {
	if err := CheckContext(ctx); err != nil {
		return err
	}
	// Create a buffer with initial size if configured
	var buf bytes.Buffer
	if r.initialSize > 0 {
		buf.Grow(r.initialSize)
	}
	// Render to buffer
	if err := r.renderer.RenderContext(ctx, &buf, data, opts...); err != nil {
		return fmt.Errorf("buffer render: %w", err)
	}
	// Get buffered content
	content := buf.Bytes()

	// Apply post-processing if configured
	if r.postRender != nil {
		if err := CheckContext(ctx); err != nil {
			return err
		}
		transformed, err := r.postRender(content)
		if err != nil {
			return fmt.Errorf("post-processing: %w", err)
		}
		content = transformed
	}
	// Check context one last time before writing
	if err := CheckContext(ctx); err != nil {
		return err
	}
	// Write final content
	_, err := w.Write(content)
	return err
}
