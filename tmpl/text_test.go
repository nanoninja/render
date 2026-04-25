// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tmpl

import (
	"bytes"
	"context"
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

	t.Run("RenderWithContentType", func(t *testing.T) {
		tpl := Text("test")
		textTpl := tpl.(*TextTemplate)

		_, err := textTpl.Parse("Hello {{ . }}")
		assert.NoError(t, err)

		var w bytes.Buffer

		err = tpl.Render(context.Background(), &w, "World", render.NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "Hello World", w.String())
		assert.Equal(t, "text/plain; charset=utf-8", tpl.ContentType())
	})

	t.Run("LoadTemplates", func(t *testing.T) {
		loader := &mockLoader{
			templates: map[string]string{
				"page.tmpl":   "Hello {{ .Name }}",
				"layout.tmpl": `Base: {{ template "content" . }}`,
			},
		}

		tpl := Text("test")
		textTpl := tpl.(*TextTemplate)
		err := textTpl.Load(loader)
		assert.NoError(t, err)

		names := []string{"page.tmpl", "layout.tmpl"}
		for _, name := range names {
			assert.NotNil(t, textTpl.Lookup(name))
		}

		var w bytes.Buffer
		data := struct{ Name string }{"World"}

		err = tpl.Render(context.Background(), &w, data, render.Options{Name: "page.tmpl"})

		assert.NoError(t, err)
		assert.Equal(t, "Hello World", w.String())
	})

	t.Run("LoadError", func(t *testing.T) {
		t.Run("Load", func(t *testing.T) {
			tpl := Text("test")
			err := tpl.(*TextTemplate).Load(&stepErrorLoader{step: 1})
			assert.Error(t, err)
			assert.StringContains(t, err.Error(), "Load error")
		})

		t.Run("Read", func(t *testing.T) {
			tpl := Text("test")
			err := tpl.(*TextTemplate).Load(&stepErrorLoader{step: 2})
			assert.Error(t, err)
			assert.StringContains(t, err.Error(), "Read error")
		})

		t.Run("Parse", func(t *testing.T) {
			tpl := Text("test")
			err := tpl.(*TextTemplate).Load(&stepErrorLoader{step: 3})
			assert.Error(t, err)
		})
	})

	t.Run("FunctionMap", func(t *testing.T) {
		funcMap := template.FuncMap{
			"upper": strings.ToUpper,
		}

		tpl := Text("test", SetFuncs(funcMap))
		textTpl := tpl.(*TextTemplate)

		_, err := textTpl.Parse(`{{ upper . }}`)
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		var w bytes.Buffer
		err = tpl.Render(context.Background(), &w, "hello", render.NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "HELLO", w.String())
	})

	t.Run("WithDefaultFuncs", func(t *testing.T) {
		tpl := Text("test", WithDefaultFuncs())
		textTpl := tpl.(*TextTemplate)

		_, err := textTpl.Parse(`{{ upper . }}`)
		assert.NoError(t, err)

		var w bytes.Buffer
		err = tpl.Render(context.Background(), &w, "hello", render.NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "HELLO", w.String())
	})

	t.Run("TemplateNotFound", func(t *testing.T) {
		tpl := Text("test")
		textTpl := tpl.(*TextTemplate)

		_, err := textTpl.Parse(`{{ . }}`)
		assert.NoError(t, err)

		var w bytes.Buffer
		err = tpl.Render(context.Background(), &w, "data", render.Options{Name: "missing.txt"})

		assert.ErrorIs(t, err, ErrTemplateNotFound)
	})

	t.Run("CustomDelimiters", func(t *testing.T) {
		tpl := Text("test", SetDelims("[[", "]]"))
		textTpl := tpl.(*TextTemplate)

		_, err := textTpl.Parse(`[[ . ]]`)
		assert.NoError(t, err)

		var w bytes.Buffer
		err = tpl.Render(context.Background(), &w, "Hello", render.NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "Hello", w.String())
	})

	t.Run("ContextCancellation", func(t *testing.T) {
		tpl := Text("test")
		textTpl := tpl.(*TextTemplate)

		_, err := textTpl.Parse(`{{ . }}`)
		assert.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		var w bytes.Buffer
		err = tpl.(*TextTemplate).Render(ctx, &w, "data", render.NoOptions)

		assert.Error(t, err)
	})
}
