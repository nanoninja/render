// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forms

import (
	"embed"
	"html/template"
)

//go:embed themes/default/*.html
var defaultThemeFS embed.FS

type defaultTheme struct{}

// ThemeDefault returns the built-in dependency-free theme.
func ThemeDefault() Theme {
	return defaultTheme{}
}

func (defaultTheme) Templates() *template.Template {
	return template.Must(template.ParseFS(defaultThemeFS, "themes/default/*.html"))
}
