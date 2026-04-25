// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"context"
	"encoding/csv"
	"fmt"
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

func (r *csvRenderer) ContentType() string {
	return ContentTypeCSV
}

// Render writes CSV data with context support.
// It accepts only [][]string data type and uses encoding/csv.Writer for output.
// The content type is set to text/csv by default but can be overridden through options.
func (r *csvRenderer) Render(ctx context.Context, w io.Writer, data any, opts Options) error {
	if err := CheckContext(ctx); err != nil {
		return err
	}
	writer := csv.NewWriter(w)
	if opts.Separator != 0 {
		writer.Comma = opts.Separator
	}
	if opts.LineEnding != "" {
		writer.UseCRLF = opts.LineEnding == "\r\n"
	}
	records, ok := data.([][]string)
	if !ok {
		return ErrInvalidData
	}
	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("csv: %w", err)
	}
	return nil
}
