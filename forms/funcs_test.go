// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forms

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/nanoninja/assert"
)

func TestFuncs(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())
	fns := Funcs(r)

	t.Run("ReturnsExactlyTheDocumentedFunctions", func(t *testing.T) {
		want := []string{"form", "form_start", "form_end", "form_row", "form_label", "form_widget", "form_errors", "form_global_errors", "form_group"}

		assert.Equal(t, len(want), len(fns))
		for _, name := range want {
			assert.HasKey(t, fns, name)
		}
	})
}

// renderTemplate executes src as a template with forms' FuncMap and
// data, returning the produced HTML.
func renderTemplate(t *testing.T, r *Renderer, src string, data any) string {
	t.Helper()

	tmpl, err := template.New("test").Funcs(Funcs(r)).Parse(src)
	assert.NoError(t, err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	assert.NoError(t, err)

	return buf.String()
}

func TestFuncs_Form(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())
	f := New("login").WithAction("/login").Add(Text("username"))

	got := renderTemplate(t, r, `{{ form . }}`, f)

	assert.StringContains(t, got, `<form`)
	assert.StringContains(t, got, `action="/login"`)
	assert.StringContains(t, got, `name="username"`)
}

func TestFuncs_FormStart(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())
	f := New("login").WithAction("/login")

	got := renderTemplate(t, r, `{{ form_start . }}`, f)

	assert.StringContains(t, got, `<form`)
	assert.StringContains(t, got, `action="/login"`)
}

func TestFuncs_FormEnd(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())
	f := New("login")

	got := renderTemplate(t, r, `{{ form_end . }}`, f)

	assert.Equal(t, "</form>", got)
}

func TestFuncs_FormRow(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())
	f := Text("username").WithLabel("Nom d'utilisateur")

	got := renderTemplate(t, r, `{{ form_row . }}`, f)

	assert.StringContains(t, got, "Nom d&#39;utilisateur")
	assert.StringContains(t, got, `type="text"`)
}

func TestFuncs_FormLabel(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())
	f := Text("username").WithLabel("Nom d'utilisateur")

	got := renderTemplate(t, r, `{{ form_label . }}`, f)

	assert.StringContains(t, got, `<label`)
	assert.StringContains(t, got, "Nom d&#39;utilisateur")
}

func TestFuncs_FormErrors(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())
	f := Text("username")
	f.Errors = []string{"required"}

	got := renderTemplate(t, r, `{{ form_errors . }}`, f)

	assert.StringContains(t, got, "required")
}

func TestFuncs_FormGlobalErrors(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())
	f := New("login")
	f.Errors = []string{"too many attempts"}

	got := renderTemplate(t, r, `{{ form_global_errors . }}`, f)

	assert.StringContains(t, got, "too many attempts")
}

func TestFuncs_FormGroup(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())
	f := New("register").Add(
		Text("first_name").WithGroup("left"),
		Email("email").WithGroup("right"),
	)

	got := renderTemplate(t, r, `{{ form_group . "left" }}`, f)

	assert.StringContains(t, got, `name="first_name"`)
}

func TestFuncs_FormWidget(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())

	t.Run("RendersBareWidget", func(t *testing.T) {
		f := Text("username")

		got := renderTemplate(t, r, `{{ form_widget . }}`, f)

		assert.StringContains(t, got, `<input`)
		assert.False(t, bytes.Contains([]byte(got), []byte("<label")))
	})

	t.Run("AppliesAttrOverridesWithoutMutatingOriginalField", func(t *testing.T) {
		f := Text("email").WithClass("form-control")

		got := renderTemplate(t, r, `{{ form_widget . "class" "email-input" }}`, f)

		assert.StringContains(t, got, `class="email-input"`)
		assert.Equal(t, "form-control", f.Attrs["class"])
	})
}
