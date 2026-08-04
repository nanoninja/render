// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forms

import (
	"html/template"

	"github.com/nanoninja/render"
)

// Funcs returns the template functions exposing r for use in HTML
// templates via tmpl.SetFuncsHTML(forms.Funcs(r)):
//
//	form                — renders a complete *Form
//	form_start          — renders the opening <form ...> tag alone
//	form_end            — renders the closing </form> tag alone
//	form_row            — renders a *Field's label, widget, and errors
//	form_label          — renders a *Field's <label> alone
//	form_widget         — renders a *Field's bare widget alone
//	form_errors         — renders a *Field's error messages alone
//	form_global_errors  — renders a *Form's form-level error messages
//	form_group          — renders the rows of every field in a group
//
// All functions use r's default theme; theme selection per call is
// not exposed here — construct a dedicated Renderer if you need a
// different theme in a given template set.
func Funcs(r *Renderer) map[string]any {
	return map[string]any{
		"form": func(f *Form) (template.HTML, error) {
			return r.Form(f, render.NoOptions)
		},
		"form_start": func(f *Form) (template.HTML, error) {
			return r.FormStart(f, render.NoOptions)
		},
		"form_end": func(f *Form) (template.HTML, error) {
			return r.FormEnd(f, render.NoOptions)
		},
		"form_row": func(f *Field) (template.HTML, error) {
			return r.Row(f, render.NoOptions)
		},
		"form_label": func(f *Field) (template.HTML, error) {
			return r.Label(f, render.NoOptions)
		},
		"form_widget": func(f *Field, attrs ...string) (template.HTML, error) {
			return r.Field(withOverrides(f, attrs...), render.NoOptions)
		},
		"form_errors": func(f *Field) (template.HTML, error) {
			return r.Errors(f, render.NoOptions)
		},
		"form_global_errors": func(f *Form) (template.HTML, error) {
			return r.FormErrors(f, render.NoOptions)
		},
		"form_group": func(f *Form, name string) (template.HTML, error) {
			return r.Group(f, name, render.NoOptions)
		},
	}
}
