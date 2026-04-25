// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tmpl

import (
	"context"
	"html/template"
	"io"
	"sync"

	"github.com/nanoninja/render"
)

// HTMLTemplate combines Go's html/template with additional rendering capabilities.
// It provides content type handling, context support, and template loading while
// preserving all standard template functionality.
type HTMLTemplate struct {
	*BaseTemplate
	*template.Template
	pool sync.Pool
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
	t.pool = sync.Pool{
		New: func() any {
			clone, err := t.Clone()
			if err != nil {
				panic("tmpl: clone failed: " + err.Error())
			}
			return clone
		},
	}
	return t
}

// HTML creates a new html template with the given name.
// It returns a render.Renderer interface for standardized rendering operations.
// This is the recommended constructor for most use cases.
//
// Example:
//
//	t := tmpl.HTML("mytemplate", tmpl.SetFuncsHTML(funcMap))
//	if err := t.Load(loader); err != nil {
//	    log.Fatal(err)
//	}
func HTML(name string, opts ...func(*HTMLTemplate)) Renderer {
	return NewHTML(template.New(name), opts...)
}

// Load loads all templates from the given loader.
// It returns an error if any template cannot be loaded or parsed.
func (t *HTMLTemplate) Load(loader Loader) error {
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

// ContentType returns the HTML MIME type produced by the template renderer.
func (t *HTMLTemplate) ContentType() string {
	return render.ContentTypeHTML
}

// Render executes the template with context and data.
// It sets the appropriate content type (text/html) and supports
// The context is checked for cancellation before rendering begins.
func (t *HTMLTemplate) Render(ctx context.Context, w io.Writer, data any, opts render.Options) error {
	if err := render.CheckContext(ctx); err != nil {
		return err
	}
	tpl := t.pool.Get().(*template.Template)

	var err error
	if opts.Name == "" {
		err = tpl.Execute(w, data)
	} else if tpl.Lookup(opts.Name) == nil {
		t.pool.Put(tpl)
		return ErrTemplateNotFound
	} else {
		err = tpl.ExecuteTemplate(w, opts.Name, data)
	}
	if err != nil {
		return err
	}
	t.pool.Put(tpl)
	return nil
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
			t.funcMap[name] = fn
		}
		t.Template = t.Funcs(funcMap)
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
		t.Template = t.Delims(left, right)
	}
}

// WithDefaultFuncsHTML applies the default function map to the template.
// It includes common utilities for string manipulation, date formatting,
// arithmetic operations, and HTML helpers (nl2br).
func WithDefaultFuncsHTML() func(*HTMLTemplate) {
	return SetFuncsHTML(DefaultFuncs())
}
