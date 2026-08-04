// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forms

import (
	"fmt"
	"html/template"
	"strings"
)

// Theme provides the HTML templates used to render fields and
// forms. Every theme must define the same set of named blocks, so
// switching themes never requires changing calling code:
//
//	field_input       — <input> widget (Type == InputType)
//	field_textarea    — <textarea> widget
//	field_select      — <select> widget
//	field_checkbox    — <input type="checkbox"> widget
//	field_checkboxes  — a group of independent checkboxes sharing one name
//	field_radio       — radio group widget
//	field_hidden      — <input type="hidden"> widget
//	field_submit      — <button type="submit"> widget
//	label             — a field's <label> alone
//	errors            — a field's error messages alone
//	form_errors       — a form's form-level error messages alone
//	row               — label + widget + errors for one field
//	form_start        — opening <form ...> tag
//	form_end          — closing </form> tag
//
// RegisterTheme validates that a theme defines every block above
// before accepting it, so a missing block fails fast at startup
// rather than lazily the first time a field of that type is
// rendered.
//
// Templates returns a full *template.Template, entirely owned by
// the theme implementation — a custom theme is free to call
// [html/template.Template.Funcs] with its own helpers (including
// [github.com/nanoninja/render/tmpl.DefaultFuncs]) before parsing,
// if its markup needs them. The built-in [ThemeDefault] and
// [ThemeBootstrap5] deliberately don't: theme templates stay
// mechanical, structural HTML — any value that needs formatting
// (dates, numbers, translated strings...) should already be in its
// final display form by the time it reaches Field.Value, set by the
// caller before rendering, not computed inside the theme.
type Theme interface {
	Templates() *template.Template
}

// requiredBlocks lists every named template block a Theme must
// define, per the contract documented on Theme.
var requiredBlocks = []string{
	"field_input",
	"field_textarea",
	"field_select",
	"field_checkbox",
	"field_checkboxes",
	"field_radio",
	"field_hidden",
	"field_submit",
	"label",
	"errors",
	"form_errors",
	"row",
	"form_start",
	"form_end",
}

// validateTheme reports an error naming every block from
// requiredBlocks that tmpl does not define.
func validateTheme(tmpl *template.Template) error {
	var missing []string
	for _, name := range requiredBlocks {
		if tmpl.Lookup(name) == nil {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required block(s): %s", strings.Join(missing, ", "))
	}
	return nil
}
