// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package render

import (
	"context"
	"fmt"
	"io"
)

// textRenderer implements text rendering with support for various data types
// and formatting options. It can handle strings, fmt.Stringer types, errors
// and any other type that can be converted to string representation.
type textRenderer struct{}

// Text creates a new TextRenderer instance configured with default options.
// The renderer uses text/plain content type by default which can be overridden
// using render options.
func Text() Renderer {
	return &textRenderer{}
}

// Render writes the string representation of data to the writer.
// It converts input data to string based on its type:
// - string: used as-is
// - fmt.Stringer: String() method is called
// - error: Error() method is called
// - others: fmt.Sprintf("%v") is used
// This is a convenience method that uses a background context.
func (r *textRenderer) Render(w io.Writer, data any, opts ...func(*Options)) error {
	return r.RenderContext(context.Background(), w, data, opts...)
}

// RenderContext is similar to Render but accepts a context for cancellation
// and timeout control. It applies any formatting options before writing the
// text output to the provided writer.
// The content type defaults to text/plain but can be overridden through options.
// When pretty-printing is enabled, a newline is added after the text.
func (r *textRenderer) RenderContext(ctx context.Context, w io.Writer, data any, opts ...func(*Options)) error {
	if err := CheckContext(ctx); err != nil {
		return err
	}
	options := NewOptions().
		Use(MimeTextPlain()).
		Use(opts...)

	var text string
	switch v := data.(type) {
	case string:
		text = fmt.Sprintf(v, options.format.args...)
	case fmt.Stringer:
		text = v.String()
	case error:
		text = v.Error()
	default:
		text = fmt.Sprintf("%v", v)
	}
	if options.format.pretty {
		text += "\n"
	}
	_, err := io.WriteString(w, text)
	return err
}
