// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package render

import "errors"

// Package level errors that can be returned by Options validation
var (
	// ErrInvalidContentType indicates that the content type format is invalid
	ErrInvalidContentType = errors.New("invalid content type")

	// ErrNegativeTimeout indicates that a negative duration was provided for timeout
	ErrNegativeTimeout = errors.New("timeout cannot be negative")

	// ErrInvalidParam indicates that a parameter value is invalid
	ErrInvalidParam = errors.New("invalid parameter value")

	// ErrRenderFailed indicates that the rendering process failed.
	// This can happen due to various reasons like I/O errors or formatting issues.
	ErrRenderFailed = errors.New("render failed")

	// ErrTemplateNotFound indicates that the requested template does not exist
	// in the template store or cannot be accessed.
	ErrTemplateNotFound = errors.New("template not found")

	// ErrInvalidData indicates that the provided data cannot be processed
	// by the renderer. This might happen when the data type is incompatible
	// with the chosen renderer.
	ErrInvalidData = errors.New("invalid data for renderer")
)
