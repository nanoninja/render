// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package tmpl

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/nanoninja/assert"
	"github.com/nanoninja/render"
)

func TestNewText(t *testing.T) {
	tpl := NewText(template.New("test"))

	assert.NotNil(t, tpl)
	assert.NotNil(t, tpl.BaseTemplate)
	assert.NotNil(t, tpl.Template)
}

func TestTextTemplate(t *testing.T) {
	t.Run("BasicInitialisation", func(t *testing.T) {
		tpl := Text("test")
		assert.NotNil(t, tpl)

		textTpl, ok := tpl.(*TextTemplate)

		assert.True(t, ok)
		assert.NotNil(t, textTpl.BaseTemplate)
		assert.NotNil(t, textTpl.Template)
	})

	t.Run("RenderWithContenType", func(t *testing.T) {
		tpl := Text("test")
		textTpl := tpl.(*TextTemplate)

		_, err := textTpl.Template.Parse("Hello {{ . }}")
		assert.Nil(t, err)

		var w bytes.Buffer
		var contentType string

		err = tpl.Render(&w, "World", func(o *render.Options) {
			contentType = o.Header().Get("Content-Type")
		})

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "Hello World")
		assert.Equals(t, contentType, "text/plain; charset=utf-8")
	})

	t.Run("LoadTemplates", func(t *testing.T) {
		loader := &mockLoader{
			templates: map[string]string{
				"page.tmpl":   "Hello {{ .Name }}",
				"layout.tmpl": `Base: {{ template "content" . }}`,
			},
		}

		tpl := Text("test", Load(loader))
		textTpl := tpl.(*TextTemplate)

		// Verify templates were loaded
		names := []string{"page.tmpl", "layout.tmpl"}
		for _, name := range names {
			assert.NotNil(t, textTpl.Lookup(name))
		}

		// Test rendering loaded template
		var w bytes.Buffer
		data := struct{ Name string }{"World"}

		err := tpl.Render(&w, data, render.Name("page.tmpl"))

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "Hello World")
	})

	t.Run("LoadError", func(t *testing.T) {
		t.Run("load error", func(t *testing.T) {
			loader := &stepErrorLoader{step: 1}

			defer func() {
				r := recover()

				if r == nil {
					t.Error("expected panic from Load error")
				}

				if !strings.Contains(fmt.Sprint(r), "Load error") {
					t.Errorf("unexpected panic message: %v", r)
				}
			}()

			Text("test", Load(loader))
		})

		t.Run("ReadError", func(t *testing.T) {
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

			Text("test", Load(loader))
		})

		t.Run("parse error", func(t *testing.T) {
			loader := &stepErrorLoader{step: 3}

			defer func() {
				r := recover()

				if r == nil {
					t.Errorf("expected panic from Parse error")
				}

				if !strings.Contains(fmt.Sprint(r), "template") {
					t.Errorf("unexpected panic message: %v", r)
				}
			}()

			Text("test", Load(loader))
		})
	})

	t.Run("FunctionMap", func(t *testing.T) {
		funcMap := template.FuncMap{
			"upper": strings.ToUpper,
		}

		tpl := Text("test", SetFuncs(funcMap))
		textTpl := tpl.(*TextTemplate)

		// Parse template using custom function
		_, err := textTpl.Parse(`{{ upper . }}`)
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		var w bytes.Buffer
		err = tpl.Render(&w, "hello")

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "HELLO")
	})

	t.Run("CustomDelimiters", func(t *testing.T) {
		tpl := Text("test", SetDelims("[[", "]]"))
		textTpl := tpl.(*TextTemplate)

		_, err := textTpl.Parse(`[[ . ]]`)
		assert.Nil(t, err)

		var w bytes.Buffer
		err = tpl.Render(&w, "Hello")

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "Hello")
	})

	t.Run("ContextCancellation", func(t *testing.T) {
		tpl := Text("test")
		textTpl := tpl.(*TextTemplate)

		_, err := textTpl.Parse(`{{ . }}`)
		assert.Nil(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		var w bytes.Buffer
		err = tpl.(*TextTemplate).RenderContext(ctx, &w, "data")

		assert.NotNil(t, err)
	})
}
