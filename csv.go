// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package render

import (
	"context"
	"encoding/csv"
	"io"
)

// csvRenderer implements CSV data rendering using encoding/csv package.
// It supports writing [][]string data with configurable delimiter and line endings.
type csvRenderer struct{}

// CSV creates a new CSV renderer with default configuration:
// - Comma as delimiter
// - Standard line endings based on encoding/csv defaults
// This is the recommended constructor for most use cases.
func CSV() Renderer {
	return &csvRenderer{}
}

// Render writes CSV data using a background context.
// Data must be [][]string type. It uses default content type of text/csv.
func (r *csvRenderer) Render(w io.Writer, data any, opts ...func(*Options)) error {
	return r.RenderContext(context.Background(), w, data, opts...)
}

// RenderContext writes CSV data with context support.
// It accepts only [][]string data type and uses encoding/csv.Writer for output.
// The content type is set to text/csv by default but can be overridden through options.
func (r *csvRenderer) RenderContext(ctx context.Context, w io.Writer, data any, opts ...func(*Options)) error {
	if err := CheckContext(ctx); err != nil {
		return err
	}
	options := NewOptions().
		Use(MimeCSV()).
		Use(opts...)

	writer := csv.NewWriter(w)
	if sep := options.params["separator"]; sep != "" {
		writer.Comma = rune(sep[0])
	}
	if options.format.lineEnding != "" {
		writer.UseCRLF = options.format.lineEnding == "\r\n"
	}
	records, ok := data.([][]string)
	if !ok {
		return ErrInvalidData
	}
	return writer.WriteAll(records)
}
