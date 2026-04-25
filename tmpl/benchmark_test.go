// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tmpl

import (
	"context"
	htmltmpl "html/template"
	"io"
	"testing"
	texttmpl "text/template"

	"github.com/nanoninja/render"
)

type benchTmplData struct {
	Name    string
	Message string
}

var benchData = benchTmplData{Name: "Gopher", Message: "Hello, World!"}

// ── HTML ──────────────────────────────────────────────────────────────────────

func BenchmarkHTMLTemplate(b *testing.B) {
	tpl := NewHTML(htmltmpl.Must(htmltmpl.New("test").Parse("<h1>{{.Name}}</h1><p>{{.Message}}</p>")))
	ctx := context.Background()
	for b.Loop() {
		_ = tpl.Render(ctx, io.Discard, benchData, render.NoOptions)
	}
}

func BenchmarkHTMLTemplateNamed(b *testing.B) {
	tpl := NewHTML(htmltmpl.Must(htmltmpl.New("base").Parse(`{{define "page"}}<h1>{{.Name}}</h1><p>{{.Message}}</p>{{end}}`)))
	ctx := context.Background()
	opts := render.Options{Name: "page"}
	for b.Loop() {
		_ = tpl.Render(ctx, io.Discard, benchData, opts)
	}
}

func BenchmarkHTMLTemplateParallel(b *testing.B) {
	tpl := NewHTML(htmltmpl.Must(htmltmpl.New("test").Parse("<h1>{{.Name}}</h1><p>{{.Message}}</p>")))
	ctx := context.Background()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = tpl.Render(ctx, io.Discard, benchData, render.NoOptions)
		}
	})
}

// ── Text ──────────────────────────────────────────────────────────────────────

func BenchmarkTextTemplate(b *testing.B) {
	tpl := NewText(texttmpl.Must(texttmpl.New("test").Parse("{{.Name}}: {{.Message}}")))
	ctx := context.Background()
	for b.Loop() {
		_ = tpl.Render(ctx, io.Discard, benchData, render.NoOptions)
	}
}

func BenchmarkTextTemplateParallel(b *testing.B) {
	tpl := NewText(texttmpl.Must(texttmpl.New("test").Parse("{{.Name}}: {{.Message}}")))
	ctx := context.Background()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = tpl.Render(ctx, io.Discard, benchData, render.NoOptions)
		}
	})
}
