// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"context"
	"fmt"
	"io"
)

// pipeRenderer chains multiple renderers where the output of each
// becomes the input of the next.
type pipeRenderer struct {
	renderers []Renderer
}

// Pipe creates a renderer that chains multiple renderers in sequence.
// The output of each renderer becomes the input data for the next one.
// The final renderer writes directly to w.
//
// Example:
//
//	render.Pipe(
//	    render.JSON(),
//	    myGzipRenderer,
//	).Render(w, data)
func Pipe(renderers ...Renderer) Renderer {
	return &pipeRenderer{renderers: renderers}
}

func (r *pipeRenderer) ContentType() string {
	return ""
}

// Render executes the pipeline with context support.
// Each renderer in the chain receives the output of the previous one as input.
// If any renderer in the chain fails, the error is returned with its position.
func (r *pipeRenderer) Render(ctx context.Context, w io.Writer, data any, opts Options) error {
	if err := CheckContext(ctx); err != nil {
		return err
	}
	if len(r.renderers) == 0 {
		return nil
	}
	ctx, cancel := ApplyTimeout(ctx, opts)
	defer cancel()

	current := data

	for i, renderer := range r.renderers {
		if i < len(r.renderers)-1 {
			// Intermediate renderer — capture output as input for the next step.
			var buf bytes.Buffer
			if err := renderer.Render(ctx, &buf, current, opts); err != nil {
				return fmt.Errorf("pipe[%d]: %w", i, err)
			}
			current = buf.Bytes()
		} else {
			// Last renderer — write directly to the destination writer.
			if err := renderer.Render(ctx, w, current, opts); err != nil {
				return fmt.Errorf("pipe[%d]: %w", i, err)
			}
		}
	}
	return nil
}
