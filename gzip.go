// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
)

type gzipRenderer struct{}

// Gzip creates a renderer that compresses input data using gzip.
// It expects []byte as input — use it as the last step in a Pipe.
//
// Example:
//
//	render.Pipe(
//	    render.JSON(),
//	    render.Gzip(),
//	).Render(w, data,
//	    render.EncodingGzip(),
//	    render.WriteResponse(w),
//	)
func Gzip() Renderer {
	return &gzipRenderer{}
}

func (r *gzipRenderer) ContentType() string {
	return ContentTypeBinary
}

func (r *gzipRenderer) Render(ctx context.Context, w io.Writer, data any, _ Options) error {
	if err := CheckContext(ctx); err != nil {
		return err
	}
	b, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("gzip: expected []byte, got %T", data)
	}
	gz := gzip.NewWriter(w)
	if _, err := gz.Write(b); err != nil {
		return fmt.Errorf("gzip: %w", err)
	}
	if err := gz.Close(); err != nil {
		return fmt.Errorf("gzip: %w", err)
	}
	return nil
}
