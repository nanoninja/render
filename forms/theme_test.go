// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forms

import (
	"fmt"
	"html/template"
	"strings"
	"testing"

	"github.com/nanoninja/assert"

	"github.com/nanoninja/render"
)

// incompleteTheme defines only a subset of the blocks required by
// the Theme contract, to exercise validateTheme's failure path.
type incompleteTheme struct{}

func (incompleteTheme) Templates() *template.Template {
	return template.Must(template.New("incomplete").Parse(`{{define "field_input"}}<input>{{end}}`))
}

func TestValidateTheme(t *testing.T) {
	t.Run("PassesForACompleteTheme", func(t *testing.T) {
		err := validateTheme(ThemeDefault().Templates())
		assert.NoError(t, err)
	})

	t.Run("PassesForBootstrap5", func(t *testing.T) {
		err := validateTheme(ThemeBootstrap5().Templates())
		assert.NoError(t, err)
	})

	t.Run("FailsAndNamesEveryMissingBlock", func(t *testing.T) {
		err := validateTheme(incompleteTheme{}.Templates())

		assert.Error(t, err)
		assert.StringContains(t, err.Error(), "field_select")
		assert.StringContains(t, err.Error(), "row")
		assert.StringContains(t, err.Error(), "form_start")
	})
}

func TestRenderer_RegisterTheme(t *testing.T) {
	t.Run("RejectsAnIncompleteTheme", func(t *testing.T) {
		r := NewRenderer("default", ThemeDefault())

		err := r.RegisterTheme("broken", incompleteTheme{})

		assert.Error(t, err)
		assert.StringContains(t, err.Error(), `"broken"`)
	})

	t.Run("DoesNotRegisterAnInvalidTheme", func(t *testing.T) {
		r := NewRenderer("default", ThemeDefault())

		_ = r.RegisterTheme("broken", incompleteTheme{})
		_, err := r.theme("broken")

		assert.Error(t, err)
	})
}

func TestNewRenderer_PanicsOnIncompleteTheme(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r)
	}()

	NewRenderer("broken", incompleteTheme{})

	t.Fatal("expected NewRenderer to panic on an incomplete theme")
}

// explodingTheme defines every block required by the Theme contract
// (so it passes validateTheme and RegisterTheme), except brokenBlock,
// which references a field that doesn't exist on whatever data it's
// executed with — a generic way to force an ExecuteTemplate failure
// at render time for any single named block, regardless of the Go
// type that block receives (Field, Form, or the internal rowData).
type explodingTheme struct {
	brokenBlock string
}

func (e explodingTheme) Templates() *template.Template {
	blocks := map[string]string{}
	for _, name := range requiredBlocks {
		blocks[name] = fmt.Sprintf(`{{define %q}}ok{{end}}`, name)
	}
	blocks[e.brokenBlock] = fmt.Sprintf(`{{define %q}}{{.Bogus}}{{end}}`, e.brokenBlock)

	var src strings.Builder
	for _, block := range blocks {
		src.WriteString(block)
	}
	return template.Must(template.New("exploding").Parse(src.String()))
}

func TestRenderer_UnknownTheme(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())
	opts := render.Options{Name: "bogus"}

	f := Text("username")
	fm := New("login").Add(f)

	cases := []struct {
		name string
		fn   func() error
	}{
		{"Label", func() error { _, err := r.Label(f, opts); return err }},
		{"Errors", func() error { _, err := r.Errors(f, opts); return err }},
		{"FormErrors", func() error { _, err := r.FormErrors(fm, opts); return err }},
		{"Row", func() error { _, err := r.Row(f, opts); return err }},
		{"FormStart", func() error { _, err := r.FormStart(fm, opts); return err }},
		{"FormEnd", func() error { _, err := r.FormEnd(fm, opts); return err }},
		{"Form", func() error { _, err := r.Form(fm, opts); return err }},
		{"Group", func() error { _, err := r.Group(fm, "left", opts); return err }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Error(t, c.fn())
		})
	}
}

func TestRenderer_TemplateExecutionErrors(t *testing.T) {
	f := Text("username")
	fm := New("login").Add(Text("username"))

	cases := []struct {
		name  string
		block string
		fn    func(r *Renderer, opts render.Options) error
	}{
		{"Field", "field_input", func(r *Renderer, opts render.Options) error {
			_, err := r.Field(f, opts)
			return err
		}},
		{"Label", "label", func(r *Renderer, opts render.Options) error {
			_, err := r.Label(f, opts)
			return err
		}},
		{"Errors", "errors", func(r *Renderer, opts render.Options) error {
			_, err := r.Errors(f, opts)
			return err
		}},
		{"FormErrors", "form_errors", func(r *Renderer, opts render.Options) error {
			_, err := r.FormErrors(fm, opts)
			return err
		}},
		{"RowWidgetFails", "field_input", func(r *Renderer, opts render.Options) error {
			_, err := r.Row(f, opts)
			return err
		}},
		{"RowTemplateFails", "row", func(r *Renderer, opts render.Options) error {
			_, err := r.Row(f, opts)
			return err
		}},
		{"FormStart", "form_start", func(r *Renderer, opts render.Options) error {
			_, err := r.FormStart(fm, opts)
			return err
		}},
		{"FormEnd", "form_end", func(r *Renderer, opts render.Options) error {
			_, err := r.FormEnd(fm, opts)
			return err
		}},
		{"FormViaFormStart", "form_start", func(r *Renderer, opts render.Options) error {
			_, err := r.Form(fm, opts)
			return err
		}},
		{"FormViaFormErrors", "form_errors", func(r *Renderer, opts render.Options) error {
			_, err := r.Form(fm, opts)
			return err
		}},
		{"FormViaRow", "row", func(r *Renderer, opts render.Options) error {
			_, err := r.Form(fm, opts)
			return err
		}},
		{"FormViaFormEnd", "form_end", func(r *Renderer, opts render.Options) error {
			_, err := r.Form(fm, opts)
			return err
		}},
		{"GroupViaRow", "row", func(r *Renderer, opts render.Options) error {
			_, err := r.Group(fm, "", opts)
			return err
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := NewRenderer("broken", explodingTheme{brokenBlock: c.block})
			err := c.fn(r, render.NoOptions)
			assert.Error(t, err)
		})
	}
}
