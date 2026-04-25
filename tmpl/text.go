// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tmpl

import (
	"context"
	"io"
	"text/template"

	"github.com/nanoninja/render"
)

// TextTemplate combines Go's text/template with additional rendering capabilities.
// It provides content type handling, context support, and template loading while
// preserving all standard template functionality.
type TextTemplate struct {
	*BaseTemplate
	*template.Template
}

// NewText creates a text template with the given template processor.
// This constructor allows for advanced template configuration by accepting
// a pre-configured template.Template and additional functional options.
func NewText(tmpl *template.Template, opts ...func(*TextTemplate)) *TextTemplate {
	t := &TextTemplate{
		BaseTemplate: New(),
		Template:     tmpl,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Text creates a new text template with the given name.
// It returns a render.Renderer interface for standardized rendering operations.
// This is the recommended constructor for most use cases.
//
// Example:
//
//	t := tmpl.Text("mytemplate", tmpl.SetFuncs(funcMap))
//	if err := t.Load(loader); err != nil {
//	    log.Fatal(err)
//	}
func Text(name string, opts ...func(*TextTemplate)) Renderer {
	return NewText(template.New(name), opts...)
}

// Load loads all templates from the given loader.
// It returns an error if any template cannot be loaded or parsed.
func (t *TextTemplate) Load(loader Loader) error {
	t.loader = loader
	templates, err := loader.Load("")
	if err != nil {
		return err
	}
	for _, name := range templates {
		content, err := loader.Read(name)
		if err != nil {
			return err
		}
		tpl, err := t.Template.New(name).Parse(string(content))
		if err != nil {
			return err
		}
		t.Template = tpl
	}
	return nil
}

// ContentType returns the plain-text MIME type produced by the template renderer.
func (t *TextTemplate) ContentType() string {
	return render.ContentTypeText
}

// Render executes the template with context and data.
// It sets the appropriate content type (text/plain) and supports
// The context is checked for cancellation before rendering begins.
func (t *TextTemplate) Render(ctx context.Context, w io.Writer, data any, opts render.Options) error {
	if err := render.CheckContext(ctx); err != nil {
		return err
	}
	tpl, err := t.Clone()
	if err != nil {
		return err
	}
	if opts.Name == "" {
		return tpl.Execute(w, data)
	}
	if tpl.Lookup(opts.Name) == nil {
		return ErrTemplateNotFound
	}
	return tpl.ExecuteTemplate(w, opts.Name, data)
}

// SetFuncs adds the provided functions to the template's function map.
// These functions become available in all templates handled by this renderer.
//
// Example:
//
//	tmpl.SetFuncs(template.FuncMap{
//	    "upper": strings.ToUpper,
//	})
func SetFuncs(funcMap template.FuncMap) func(*TextTemplate) {
	return func(t *TextTemplate) {
		for name, fn := range funcMap {
			t.funcMap[name] = fn
		}
		t.Template = t.Funcs(t.funcMap)
	}
}

// SetDelims sets the template delimiters to the specified strings.
// This is useful when the default delimiters ({{ and }}) conflict with
// the template content.
//
// Example:
//
//	tmpl.SetDelims("[[", "]]")
func SetDelims(left, right string) func(*TextTemplate) {
	return func(t *TextTemplate) {
		t.Template = t.Delims(left, right)
	}
}

// WithDefaultFuncs applies the default function map to the template.
// It includes common utilities for string manipulation, date formatting
// and arithmetic operations.
func WithDefaultFuncs() func(*TextTemplate) {
	return SetFuncs(DefaultFuncs())
}
