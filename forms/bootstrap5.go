// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forms

import (
	"embed"
	"html/template"
)

//go:embed themes/bootstrap5/*.html
var bootstrap5ThemeFS embed.FS

type bootstrap5Theme struct{}

// ThemeBootstrap5 returns a theme styled for Bootstrap 5
// (form-control, form-label, is-invalid/invalid-feedback...).
func ThemeBootstrap5() Theme {
	return bootstrap5Theme{}
}

func (bootstrap5Theme) Templates() *template.Template {
	return template.Must(template.ParseFS(bootstrap5ThemeFS, "themes/bootstrap5/*.html"))
}
