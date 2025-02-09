// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package tmpl

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strings"
	"testing"

	"github.com/nanoninja/assert"
	"github.com/nanoninja/render"
)

func TestNewHTML(t *testing.T) {
	tpl := NewHTML(template.New("test"))

	assert.NotNil(t, tpl)
	assert.NotNil(t, tpl.BaseTemplate)
	assert.NotNil(t, tpl.Template)
}

func TestHTMLemplate(t *testing.T) {
	t.Run("BasicInitialisation", func(t *testing.T) {
		tpl := HTML("test")
		assert.NotNil(t, tpl)

		htmlTpl, ok := tpl.(*HTMLTemplate)

		assert.True(t, ok, "HTML() dit not return *HTMLTemplate")
		assert.NotNil(t, htmlTpl.BaseTemplate)
		assert.NotNil(t, htmlTpl.Template)
	})

	t.Run("RenderWithContentType", func(t *testing.T) {
		tpl := HTML("test")
		htmlTpl := tpl.(*HTMLTemplate)

		_, err := htmlTpl.Template.Parse("<h1>Hello {{ . }}</h1>")
		assert.Nil(t, err)

		var w bytes.Buffer
		var contentType string

		err = tpl.Render(&w, "World", func(o *render.Options) {
			contentType = o.Header().Get("Content-Type")
		})

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "<h1>Hello World</h1>")
		assert.Equals(t, contentType, "text/html; charset=utf-8")
	})

	t.Run("LoadTemplates", func(t *testing.T) {
		loader := &mockLoader{
			templates: map[string]string{
				"page.html":   `<h1>Hello {{ .Name }}</h1>`,
				"layout.html": `<div> {{ template "content" . }}</div>`,
			},
		}

		tpl := HTML("test", LoadHTML(loader))
		htmlTpl := tpl.(*HTMLTemplate)

		// Verify template werer loaded
		names := []string{"page.html", "layout.html"}
		for _, name := range names {
			assert.NotNil(t, htmlTpl.Lookup(name))
		}

		// Test rendering loaded template
		var w bytes.Buffer
		data := struct{ Name string }{"World"}

		err := tpl.Render(&w, data, render.Name("page.html"))

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "<h1>Hello World</h1>")
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Load", func(t *testing.T) {
			loader := &stepErrorLoader{step: 1}

			defer func() {
				r := recover()

				if r == nil {
					t.Errorf("expected panic from load error")
				}

				if !strings.Contains(fmt.Sprint(r), "Load error") {
					t.Errorf("unexpected panic message: %v", r)
				}
			}()

			HTML("test", LoadHTML(loader))
		})

		t.Run("Read", func(t *testing.T) {
			loader := &stepErrorLoader{step: 2}

			defer func() {
				r := recover()

				if r == nil {
					t.Error("expected panic from Read error")
				}

				if !strings.Contains(fmt.Sprint(r), "Read error") {
					t.Errorf("unexpected panic message: %v", r)
				}
			}()

			HTML("test", LoadHTML(loader))
		})

		t.Run("Parse", func(t *testing.T) {
			loader := &stepErrorLoader{step: 3}

			defer func() {
				r := recover()

				if r == nil {
					t.Error("expected panic from Parse error")
				}

				if !strings.Contains(fmt.Sprint(r), "template") {
					t.Errorf("unexpected panic message: %v", r)
				}
			}()

			HTML("test", LoadHTML(loader))
		})
	})

	t.Run("FunctionMap", func(t *testing.T) {
		funcMap := template.FuncMap{
			"upper": strings.ToUpper,
		}

		tpl := HTML("test", SetFuncsHTML(funcMap))
		htmlTpl := tpl.(*HTMLTemplate)

		// Parse template using custom function
		_, err := htmlTpl.Parse(`<h1>{{ upper . }}</h1>`)

		assert.Nil(t, err)

		var w bytes.Buffer
		err = tpl.Render(&w, "hello")

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "<h1>HELLO</h1>")
	})

	t.Run("CustomDelimiter", func(t *testing.T) {
		tpl := HTML("test", SetDelimsHTML("[[", "]]"))
		htmlTpl := tpl.(*HTMLTemplate)

		_, err := htmlTpl.Parse(`<h1>[[ . ]]</h1>`)
		if err != nil {
			t.Errorf("failed to parse template: %v", err)
		}

		var w bytes.Buffer
		err = tpl.Render(&w, "Hello")

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "<h1>Hello</h1>")
	})

	t.Run("ContextCancellation", func(t *testing.T) {
		tpl := HTML("test")
		htmlTpl := tpl.(*HTMLTemplate)

		_, err := htmlTpl.Parse(`{{ . }}`)

		assert.Nil(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		var w bytes.Buffer
		err = tpl.RenderContext(ctx, &w, "data")

		assert.NotNil(t, err)
	})
}
