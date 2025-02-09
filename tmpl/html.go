// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package tmpl

import (
	"context"
	"html/template"
	"io"

	"github.com/nanoninja/render"
)

// HTMLTemplate combines Go's html/template with additional rendering capabilities.
// It provides content type handling, context support, and template loading while
// preserving all standard template functionality.
type HTMLTemplate struct {
	*BaseTemplate
	*template.Template
}

// NewHTML creates a text template with the given template processor.
// This constructor allows for advanced template configuration by accepting
// a pre-configured template.Template and additional functional options.
func NewHTML(tmpl *template.Template, opts ...func(*HTMLTemplate)) *HTMLTemplate {
	t := &HTMLTemplate{
		BaseTemplate: New(),
		Template:     tmpl,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// HTML creates a new html template with the given name.
// It returns a render.Renderer interface for standardized rendering operations.
// This is the recommended constructor for most use cases.
//
// Example:
//
//	t := tmpl.HTML("mytemplate",
//	    tmpl.LoadHTML(loader),
//	    tmpl.SetFuncsHTML(funcMap),
//	)
func HTML(name string, opts ...func(*HTMLTemplate)) render.Renderer {
	return NewHTML(template.New(name), opts...)
}

// Render executes the template with the given data.
// It provides a convenient wrapper around RenderContext using a background context.
func (t *HTMLTemplate) Render(w io.Writer, data any, opts ...func(*render.Options)) error {
	return t.RenderContext(context.Background(), w, data, opts...)
}

// RenderContext executes the template with context and data.
// It sets the appropriate content type (text/html) and supports
// all render options. The context allows for cancellation and timeout control.
func (t *HTMLTemplate) RenderContext(ctx context.Context, w io.Writer, data any, opts ...func(*render.Options)) error {
	if err := render.CheckContext(ctx); err != nil {
		return err
	}
	options := render.NewOptions().
		Use(render.MimeTextHTML()).
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

// LoadHTML configures a template loader and loads all templates.
// Templates are loaded relative to the configured root directory
// and can be referenced by their path during rendering.
//
// Example:
//
//	tmpl.LoadHTML(loader)
func LoadHTML(loader Loader) func(*HTMLTemplate) {
	return func(t *HTMLTemplate) {
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

// SetFuncsHTML adds the provided functions to the template's function map.
// These functions become available in all templates handled by this renderer.
//
// Example:
//
//	tmpl.SetFuncsHTML(template.FuncMap{
//	    "upper": strings.ToUpper,
//	})
func SetFuncsHTML(funcMap map[string]any) func(*HTMLTemplate) {
	return func(t *HTMLTemplate) {
		for name, fn := range funcMap {
			t.BaseTemplate.funcMap[name] = fn
		}
		t.Template = t.Template.Funcs(funcMap)
	}
}

// SetDelimsHTML sets the template delimiters to the specified strings.
// This is useful when the default delimiters ({{ and }}) conflict with
// the template content.
//
// Example:
//
//	tmpl.SetDelimsHTML("[[", "]]")
func SetDelimsHTML(left, right string) func(*HTMLTemplate) {
	return func(t *HTMLTemplate) {
		t.Template = t.Template.Delims(left, right)
	}
}
