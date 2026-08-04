// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forms

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"

	"github.com/nanoninja/render"
)

// Renderer produces HTML for fields and forms using one or more
// registered Theme instances. It implements render.Renderer, so it
// composes with render.Buffer, render.Cache, render.Gzip, and
// render.Multi like any other renderer.
//
// Public methods (Field, Label, Errors, Row, Form, Group...) each
// resolve the theme once and allocate a single buffer, delegating
// to the unexported write* helpers below for the actual template
// execution. Those helpers write directly into the caller's
// io.Writer/buffer instead of each allocating and returning their
// own — composed rendering (Row over Field, Form over Row, Group
// over Row) does not re-buffer or re-resolve the theme once per
// sub-piece.
type Renderer struct {
	themes       map[string]*template.Template
	defaultTheme string
}

// NewRenderer creates a Renderer with t registered as its default
// theme under name. It panics if t does not satisfy the Theme
// contract (see requiredBlocks) — a renderer's default theme is
// typically the caller's own, trusted, hard-coded choice, so a
// misconfiguration here is a programmer error meant to fail loudly
// at startup, the same way template.Must does for malformed
// template syntax. Use RegisterTheme for themes registered after
// construction, which returns an error instead.
func NewRenderer(name string, t Theme) *Renderer {
	r := &Renderer{
		themes:       make(map[string]*template.Template),
		defaultTheme: name,
	}
	if err := r.RegisterTheme(name, t); err != nil {
		panic(err)
	}
	return r
}

// RegisterTheme adds t to the renderer under name, making it
// selectable via render.Options.Name. It returns an error if t does
// not define every block required by the Theme contract.
func (r *Renderer) RegisterTheme(name string, t Theme) error {
	tmpl := t.Templates()
	if err := validateTheme(tmpl); err != nil {
		return fmt.Errorf("forms: theme %q: %w", name, err)
	}
	r.themes[name] = tmpl
	return nil
}

// theme resolves the template set to use for a given Options.Name,
// falling back to the default theme when empty.
func (r *Renderer) theme(name string) (*template.Template, error) {
	if name == "" {
		name = r.defaultTheme
	}
	t, ok := r.themes[name]
	if !ok {
		return nil, fmt.Errorf("forms: unknown theme %q", name)
	}
	return t, nil
}

// writeField writes f's widget alone into w using the already
// resolved tmpl.
func (r *Renderer) writeField(w io.Writer, tmpl *template.Template, f *Field) error {
	return tmpl.ExecuteTemplate(w, f.Type.blockName(), f)
}

// writeLabel writes f's <label> alone into w.
func (r *Renderer) writeLabel(w io.Writer, tmpl *template.Template, f *Field) error {
	return tmpl.ExecuteTemplate(w, "label", f)
}

// writeErrors writes f's error messages alone into w.
func (r *Renderer) writeErrors(w io.Writer, tmpl *template.Template, f *Field) error {
	return tmpl.ExecuteTemplate(w, "errors", f)
}

// writeFormErrors writes fm's form-level error messages into w.
func (r *Renderer) writeFormErrors(w io.Writer, tmpl *template.Template, fm *Form) error {
	return tmpl.ExecuteTemplate(w, "form_errors", fm)
}

// writeFormStart writes fm's opening <form ...> tag into w.
func (r *Renderer) writeFormStart(w io.Writer, tmpl *template.Template, fm *Form) error {
	return tmpl.ExecuteTemplate(w, "form_start", fm)
}

// writeFormEnd writes fm's closing </form> tag into w.
func (r *Renderer) writeFormEnd(w io.Writer, tmpl *template.Template, fm *Form) error {
	return tmpl.ExecuteTemplate(w, "form_end", fm)
}

// rowData is the data passed to the "row" template block: the
// field itself (for label, name, errors) plus its already-rendered
// widget HTML.
type rowData struct {
	Field  *Field
	Widget template.HTML
}

// writeRow writes f's label, widget, and errors into w. The widget
// still needs its own small buffer, since rowData requires it as
// an already-rendered template.HTML value — but nothing downstream
// (Form, Group) re-buffers or re-resolves the theme per field.
func (r *Renderer) writeRow(w io.Writer, tmpl *template.Template, f *Field) error {
	var widgetBuf bytes.Buffer
	if err := r.writeField(&widgetBuf, tmpl, f); err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(w, "row", rowData{Field: f, Widget: template.HTML(widgetBuf.String())})
}

// writeForm writes fm's complete <form>...</form> into w, iterating
// fields in declaration order through writeRow.
func (r *Renderer) writeForm(w io.Writer, tmpl *template.Template, fm *Form) error {
	if err := r.writeFormStart(w, tmpl, fm); err != nil {
		return err
	}
	if err := r.writeFormErrors(w, tmpl, fm); err != nil {
		return err
	}
	for _, f := range fm.Fields {
		if err := r.writeRow(w, tmpl, f); err != nil {
			return err
		}
	}
	return r.writeFormEnd(w, tmpl, fm)
}

// writeGroup writes, concatenated, the rows of every field in fm
// tagged with the given group name.
func (r *Renderer) writeGroup(w io.Writer, tmpl *template.Template, fm *Form, name string) error {
	for _, f := range fm.Group(name) {
		if err := r.writeRow(w, tmpl, f); err != nil {
			return err
		}
	}
	return nil
}

// Field renders a single field's widget, e.g. the bare <input>,
// without label or errors.
func (r *Renderer) Field(f *Field, opts render.Options) (template.HTML, error) {
	tmpl, err := r.theme(opts.Name)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := r.writeField(&buf, tmpl, f); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// Label renders a field's <label> element alone.
func (r *Renderer) Label(f *Field, opts render.Options) (template.HTML, error) {
	tmpl, err := r.theme(opts.Name)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := r.writeLabel(&buf, tmpl, f); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// Errors renders a field's error messages alone.
func (r *Renderer) Errors(f *Field, opts render.Options) (template.HTML, error) {
	tmpl, err := r.theme(opts.Name)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := r.writeErrors(&buf, tmpl, f); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// FormErrors renders a form's global (non-field) error messages.
func (r *Renderer) FormErrors(fm *Form, opts render.Options) (template.HTML, error) {
	tmpl, err := r.theme(opts.Name)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := r.writeFormErrors(&buf, tmpl, fm); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// ContentType implements render.Renderer.
func (r *Renderer) ContentType() string {
	return "text/html; charset=utf-8"
}

// Render implements render.Renderer. data must be a *Field or
// *Form. It writes directly into w — no intermediate buffer, since
// w is already the final destination.
func (r *Renderer) Render(ctx context.Context, w io.Writer, data any, opts render.Options) error {
	if err := render.CheckContext(ctx); err != nil {
		return err
	}
	tmpl, err := r.theme(opts.Name)
	if err != nil {
		return err
	}
	switch v := data.(type) {
	case *Field:
		return r.writeField(w, tmpl, v)
	case *Form:
		return r.writeForm(w, tmpl, v)
	default:
		return fmt.Errorf("forms: unsupported data type %T", data)
	}
}

// Row renders a field's label, widget, and errors together.
func (r *Renderer) Row(f *Field, opts render.Options) (template.HTML, error) {
	tmpl, err := r.theme(opts.Name)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := r.writeRow(&buf, tmpl, f); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// FormStart renders the opening <form ...> tag alone, for manual
// assembly (e.g. to insert a submit button before FormEnd).
func (r *Renderer) FormStart(fm *Form, opts render.Options) (template.HTML, error) {
	tmpl, err := r.theme(opts.Name)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := r.writeFormStart(&buf, tmpl, fm); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// FormEnd renders the closing </form> tag alone.
func (r *Renderer) FormEnd(fm *Form, opts render.Options) (template.HTML, error) {
	tmpl, err := r.theme(opts.Name)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := r.writeFormEnd(&buf, tmpl, fm); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// Form renders a complete <form>...</form>, iterating fields in
// declaration order.
func (r *Renderer) Form(fm *Form, opts render.Options) (template.HTML, error) {
	tmpl, err := r.theme(opts.Name)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := r.writeForm(&buf, tmpl, fm); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// Group renders, concatenated, the rows of every field tagged with
// the given group name via Field.WithGroup — for two-column layouts
// or multi-step wizards, where the caller controls the surrounding
// markup (columns, step containers, form_start/form_end...).
func (r *Renderer) Group(fm *Form, name string, opts render.Options) (template.HTML, error) {
	tmpl, err := r.theme(opts.Name)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := r.writeGroup(&buf, tmpl, fm, name); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}
