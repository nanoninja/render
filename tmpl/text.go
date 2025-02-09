// Copyright 2025 The Nanoninja Authors. All rights reserved.
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
//	t := tmpl.Text("mytemplate",
//	    tmpl.SetFuncs(funcMap),
//	    tmpl.Load(loader),
//	)
func Text(name string, opts ...func(*TextTemplate)) render.Renderer {
	return NewText(template.New(name), opts...)
}

// Render executes the template with the given data.
// It provides a convenient wrapper around RenderContext using a background context.
func (t *TextTemplate) Render(w io.Writer, data any, opts ...func(*render.Options)) error {
	return t.RenderContext(context.Background(), w, data, opts...)
}

// RenderContext executes the template with context and data.
// It sets the appropriate content type (text/plain) and supports
// all render options. The context allows for cancellation and timeout control.
func (t *TextTemplate) RenderContext(ctx context.Context, w io.Writer, data any, opts ...func(*render.Options)) error {
	if err := render.CheckContext(ctx); err != nil {
		return err
	}
	options := render.NewOptions().
		Use(render.MimeTextPlain()).
		Use(opts...)

	tpl, err := t.Clone()
	if err != nil {
		return err
	}
	if options.Name() == "" {
		return tpl.Execute(w, data)
	}
	return tpl.ExecuteTemplate(w, options.Name(), data)
}

// Load configures a template loader and loads all templates.
// Templates are loaded relative to the configured root directory
// and can be referenced by their path during rendering.
//
// Example:
//
//	tmpl.Load(loader)
func Load(loader Loader) func(*TextTemplate) {
	return func(t *TextTemplate) {
		t.loader = loader
		templates, err := loader.Load("")
		if err != nil {
			panic(err)
		}
		for _, name := range templates {
			tpl := t.Template.New(name)
			content, err := loader.Read(name)
			if err != nil {
				panic(err)
			}
			tpl, err = tpl.Parse(string(content))
			if err != nil {
				panic(err)
			}
			t.Template = tpl
		}
	}
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
			t.BaseTemplate.funcMap[name] = fn
		}
		t.Template = t.Template.Funcs(t.funcMap)
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
		t.Template = t.Template.Delims(left, right)
	}
}
