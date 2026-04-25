// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tmpl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"strings"
	"sync"
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

		assert.True(t, ok, "HTML() did not return *HTMLTemplate")
		assert.NotNil(t, htmlTpl.BaseTemplate)
		assert.NotNil(t, htmlTpl.Template)
	})

	t.Run("RenderWithContentType", func(t *testing.T) {
		tpl := HTML("test")
		htmlTpl := tpl.(*HTMLTemplate)

		_, err := htmlTpl.Parse("<h1>Hello {{ . }}</h1>")
		assert.NoError(t, err)

		var w bytes.Buffer

		err = tpl.Render(context.Background(), &w, "World", render.NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "<h1>Hello World</h1>", w.String())
		assert.Equal(t, "text/html; charset=utf-8", tpl.ContentType())
	})

	t.Run("LoadTemplates", func(t *testing.T) {
		loader := &mockLoader{
			templates: map[string]string{
				"page.html":   `<h1>Hello {{ .Name }}</h1>`,
				"layout.html": `<div> {{ template "content" . }}</div>`,
			},
		}

		tpl := HTML("test")
		htmlTpl := tpl.(*HTMLTemplate)
		err := htmlTpl.Load(loader)
		assert.NoError(t, err)

		names := []string{"page.html", "layout.html"}
		for _, name := range names {
			assert.NotNil(t, htmlTpl.Lookup(name))
		}

		var w bytes.Buffer
		data := struct{ Name string }{"World"}

		err = tpl.Render(context.Background(), &w, data, render.Options{Name: "page.html"})

		assert.NoError(t, err)
		assert.Equal(t, "<h1>Hello World</h1>", w.String())
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Load", func(t *testing.T) {
			tpl := HTML("test")
			err := tpl.(*HTMLTemplate).Load(&stepErrorLoader{step: 1})

			assert.Error(t, err)
			assert.StringContains(t, err.Error(), "Load error")
		})

		t.Run("Read", func(t *testing.T) {
			tpl := HTML("test")
			err := tpl.(*HTMLTemplate).Load(&stepErrorLoader{step: 2})
			assert.Error(t, err)
			assert.StringContains(t, err.Error(), "Read error")
		})

		t.Run("Parse", func(t *testing.T) {
			tpl := HTML("test")
			err := tpl.(*HTMLTemplate).Load(&stepErrorLoader{step: 3})
			assert.Error(t, err)
		})
	})

	t.Run("FunctionMap", func(t *testing.T) {
		funcMap := template.FuncMap{
			"upper": strings.ToUpper,
		}

		tpl := HTML("test", SetFuncsHTML(funcMap))
		htmlTpl := tpl.(*HTMLTemplate)

		_, err := htmlTpl.Parse(`<h1>{{ upper . }}</h1>`)

		assert.NoError(t, err)

		var w bytes.Buffer
		err = tpl.Render(context.Background(), &w, "hello", render.NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "<h1>HELLO</h1>", w.String())
	})

	t.Run("WithDefaultFuncs", func(t *testing.T) {
		tpl := HTML("test", WithDefaultFuncsHTML())
		htmlTpl := tpl.(*HTMLTemplate)

		_, err := htmlTpl.Parse(`<p>{{ upper . }}</p>`)
		assert.NoError(t, err)

		var w bytes.Buffer
		err = tpl.Render(context.Background(), &w, "hello", render.NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "<p>HELLO</p>", w.String())
	})

	t.Run("CustomDelimiter", func(t *testing.T) {
		tpl := HTML("test", SetDelimsHTML("[[", "]]"))
		htmlTpl := tpl.(*HTMLTemplate)

		_, err := htmlTpl.Parse(`<h1>[[ . ]]</h1>`)
		if err != nil {
			t.Errorf("failed to parse template: %v", err)
		}

		var w bytes.Buffer
		err = tpl.Render(context.Background(), &w, "Hello", render.NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "<h1>Hello</h1>", w.String())
	})

	t.Run("TemplateNotFound", func(t *testing.T) {
		tpl := HTML("test")
		htmlTpl := tpl.(*HTMLTemplate)

		_, err := htmlTpl.Parse(`<h1>{{ . }}</h1>`)
		assert.NoError(t, err)

		var w bytes.Buffer
		err = tpl.Render(context.Background(), &w, "data", render.Options{Name: "missing.html"})

		assert.ErrorIs(t, err, ErrTemplateNotFound)
	})

	t.Run("ContextCancellation", func(t *testing.T) {
		tpl := HTML("test")
		htmlTpl := tpl.(*HTMLTemplate)

		_, err := htmlTpl.Parse(`{{ . }}`)

		assert.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		var w bytes.Buffer
		err = tpl.Render(ctx, &w, "data", render.NoOptions)

		assert.Error(t, err)
	})

	t.Run("ExecuteError", func(t *testing.T) {
		// FuncMap must be registered before Parse in html/template.
		// A func returning an error triggers the Execute error path
		// where the clone must NOT be returned to the pool.
		failFn := template.FuncMap{
			"fail": func() (string, error) { return "", errors.New("execute failed") },
		}
		base := template.Must(template.New("test").Funcs(failFn).Parse(`{{fail}}`))
		tpl := NewHTML(base)
		var w bytes.Buffer
		err := tpl.Render(context.Background(), &w, nil, render.NoOptions)
		assert.Error(t, err)
	})

	t.Run("ConcurrentRenderIsolation", func(t *testing.T) {
		// Vérifie que deux goroutines concurrentes obtiennent
		// chacune le résultat de leurs propres données — pas celles de l'autre.
		tpl := HTML("test")
		htmlTpl := tpl.(*HTMLTemplate)

		_, err := htmlTpl.Parse(`<p>{{.Name}}</p>`)
		assert.NoError(t, err)

		const workers = 50
		var wg sync.WaitGroup
		wg.Add(workers)

		for i := range workers {
			go func(id int) {
				defer wg.Done()
				data := struct{ Name string }{Name: fmt.Sprintf("user-%d", id)}
				var w bytes.Buffer
				err := tpl.Render(context.Background(), &w, data, render.NoOptions)
				assert.NoError(t, err)
				assert.StringContains(t, w.String(), data.Name)
			}(i)
		}
		wg.Wait()
	})
}
