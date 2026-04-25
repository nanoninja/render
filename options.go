// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import "time"

// Options configures a single render call. It is passed by value so no
// allocation occurs on the hot path. Zero value is valid and means "use
// renderer defaults".
type Options struct {
	// Timeout limits how long rendering may take. Zero means no limit.
	Timeout time.Duration

	// Name selects the template to execute (tmpl renderers only).
	Name string

	// Pretty enables human-readable, indented output.
	Pretty bool

	// Indent sets the indentation string for JSON and XML output.
	// Ignored when Pretty is false.
	Indent string

	// Prefix sets the per-line prefix string for JSON and XML output.
	// Ignored when Pretty is false.
	Prefix string

	// LineEnding overrides the default line ending character(s).
	// Used by the CSV renderer. Common values: "\n", "\r\n".
	LineEnding string

	// Args holds printf-style arguments for text rendering.
	// Used by the Text renderer when the data string contains verbs.
	Args []any

	// Separator sets the CSV field separator rune. Default is ','.
	Separator rune

	// Padding sets the JSONP callback function name (JSON renderer only).
	// When non-empty the output is wrapped: callback(json).
	Padding string

	// Headers contains additional response headers to set.
	// Nil by default — allocate only when custom headers are needed.
	Headers map[string][]string
}

// NoOptions is the reusable zero-value Options preset.
var NoOptions = Options{}
