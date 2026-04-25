// Copyright 2026 The Nanoninja Authors. All rights reserved.
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

func (r *textRenderer) ContentType() string {
	return ContentTypeText
}

func (r *textRenderer) Render(ctx context.Context, w io.Writer, data any, opts Options) error {
	if err := CheckContext(ctx); err != nil {
		return err
	}
	var text string
	switch v := data.(type) {
	case string:
		if len(opts.Args) > 0 {
			text = fmt.Sprintf(v, opts.Args...)
		} else {
			text = v
		}
	case fmt.Stringer:
		text = v.String()
	case error:
		text = v.Error()
	default:
		text = fmt.Sprintf("%v", v)
	}
	if opts.Pretty {
		text += "\n"
	}
	_, err := io.WriteString(w, text)
	return err
}
