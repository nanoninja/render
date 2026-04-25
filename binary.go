// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"context"
	"fmt"
	"io"
)

type binaryRenderer struct{}

// Binary creates a renderer for raw binary content such as files, images or PDFs.
// It sets the Content-Type to application/octet-stream by default.
// Use the Attachment option to set a Content-Disposition header for file downloads.
func Binary() Renderer {
	return &binaryRenderer{}
}

func (r binaryRenderer) ContentType() string {
	return ContentTypeBinary
}

func (r *binaryRenderer) Render(ctx context.Context, w io.Writer, data any, _ Options) error {
	if err := CheckContext(ctx); err != nil {
		return err
	}
	switch v := data.(type) {
	case []byte:
		_, err := w.Write(v)
		return err
	case io.Reader:
		_, err := io.Copy(w, v)
		return err
	default:
		return fmt.Errorf("binary: unsupported data type %T", data)
	}
}
