// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package tmpl

import "text/template"

// BaseTemplate provides common template functionality.
// It serves as a foundation for both HTML and text templates.
type BaseTemplate struct {
	loader  Loader
	funcMap template.FuncMap
}

// New creates a new base template with the given template implementation.
// It initializes common components like function maps and renderer.
func New() *BaseTemplate {
	return &BaseTemplate{
		funcMap: make(template.FuncMap),
	}
}
