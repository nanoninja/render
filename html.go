// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"context"
	"fmt"
	"io"
)

// htmlRenderer implements raw HTML rendering.
// For template-based HTML rendering, use the tmpl sub-package instead.
type htmlRenderer struct{}

// HTML creates a renderer that writes raw HTML content.
// It sets the Content-Type to text/html automatically.
// For template-based rendering, use the tmpl sub-package instead.
func HTML() Renderer {
	return &htmlRenderer{}
}

func (r *htmlRenderer) ContentType() string {
	return ContentTypeHTML
}

// Render writes the HTML representation of data with context support.
// It handles:
// - string and []byte written as-is
// - any other type converted via fmt.Fprint
// - Content type setting to text/html
func (r *htmlRenderer) Render(ctx context.Context, w io.Writer, data any, _ Options) error {
	if err := CheckContext(ctx); err != nil {
		return err
	}
	switch v := data.(type) {
	case string:
		_, err := io.WriteString(w, v)
		return err
	case []byte:
		_, err := w.Write(v)
		return err
	default:
		_, err := fmt.Fprint(w, v)
		return err
	}
}
