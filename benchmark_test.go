// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package render benchmarks.
package render

import (
	"bytes"
	"context"
	"io"
	"testing"
)

// benchStruct is the realistic benchmark payload — reflects what a web handler
// typically sends. Struct encoding avoids reflect.MapIter allocs (map variant
// is kept separately to document that overhead).
type benchStruct struct {
	Name    string `json:"name"    xml:"name"    yaml:"name"`
	Age     int    `json:"age"     xml:"age"     yaml:"age"`
	Active  bool   `json:"active"  xml:"active"  yaml:"active"`
	Message string `json:"message" xml:"message" yaml:"message"`
}

var (
	benchStructData = benchStruct{
		Name:    "Gopher",
		Age:     10,
		Active:  true,
		Message: "Hello, World!",
	}

	// benchMapData documents reflect.MapIter overhead — not the normal case.
	benchMapData = map[string]any{
		"name":    "Gopher",
		"age":     10,
		"active":  true,
		"message": "Hello, World!",
	}

	benchCSVData = [][]string{
		{"name", "age", "city"},
		{"Alice", "30", "Paris"},
		{"Bob", "25", "Lyon"},
	}

	benchBinaryData = bytes.Repeat([]byte{0xFF}, 1024)
)

// ── JSON ──────────────────────────────────────────────────────────────────────

func BenchmarkJSON(b *testing.B) {
	r := JSON()
	ctx := context.Background()
	for b.Loop() {
		_ = r.Render(ctx, io.Discard, benchStructData, NoOptions)
	}
}

// BenchmarkJSONMap documents the reflect.MapIter overhead vs a struct.
func BenchmarkJSONMap(b *testing.B) {
	r := JSON()
	ctx := context.Background()
	for b.Loop() {
		_ = r.Render(ctx, io.Discard, benchMapData, NoOptions)
	}
}

func BenchmarkJSONPretty(b *testing.B) {
	r := JSON()
	ctx := context.Background()
	opts := Options{Pretty: true}
	for b.Loop() {
		_ = r.Render(ctx, io.Discard, benchStructData, opts)
	}
}

func BenchmarkJSONParallel(b *testing.B) {
	r := JSON()
	ctx := context.Background()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = r.Render(ctx, io.Discard, benchStructData, NoOptions)
		}
	})
}

// ── XML ───────────────────────────────────────────────────────────────────────

func BenchmarkXML(b *testing.B) {
	r := XML()
	ctx := context.Background()
	for b.Loop() {
		_ = r.Render(ctx, io.Discard, benchStructData, NoOptions)
	}
}

func BenchmarkXMLPretty(b *testing.B) {
	r := XML()
	ctx := context.Background()
	opts := Options{Pretty: true}
	for b.Loop() {
		_ = r.Render(ctx, io.Discard, benchStructData, opts)
	}
}

// ── YAML ──────────────────────────────────────────────────────────────────────

func BenchmarkYAML(b *testing.B) {
	r := YAML()
	ctx := context.Background()
	for b.Loop() {
		_ = r.Render(ctx, io.Discard, benchStructData, NoOptions)
	}
}

// ── CSV ───────────────────────────────────────────────────────────────────────

func BenchmarkCSV(b *testing.B) {
	r := CSV()
	ctx := context.Background()
	for b.Loop() {
		_ = r.Render(ctx, io.Discard, benchCSVData, NoOptions)
	}
}

// ── Text / HTML ───────────────────────────────────────────────────────────────

func BenchmarkText(b *testing.B) {
	r := Text()
	ctx := context.Background()
	for b.Loop() {
		_ = r.Render(ctx, io.Discard, "Hello, Gopher!", NoOptions)
	}
}

func BenchmarkHTML(b *testing.B) {
	r := HTML()
	ctx := context.Background()
	for b.Loop() {
		_ = r.Render(ctx, io.Discard, "<h1>Hello</h1>", NoOptions)
	}
}

// ── Binary ────────────────────────────────────────────────────────────────────

func BenchmarkBinaryBytes(b *testing.B) {
	r := Binary()
	ctx := context.Background()
	for b.Loop() {
		_ = r.Render(ctx, io.Discard, benchBinaryData, NoOptions)
	}
}

// ── Buffer ────────────────────────────────────────────────────────────────────

// BenchmarkBuffer uses bytes.Buffer (not io.Discard) because BufferRenderer
// needs to write to a real writer to exercise its copy path.
func BenchmarkBuffer(b *testing.B) {
	r := Buffer(JSON())
	ctx := context.Background()
	var w bytes.Buffer
	for b.Loop() {
		w.Reset()
		_ = r.Render(ctx, &w, benchStructData, NoOptions)
	}
}

// ── Gzip ──────────────────────────────────────────────────────────────────────

func BenchmarkGzip(b *testing.B) {
	r := Gzip()
	ctx := context.Background()
	for b.Loop() {
		_ = r.Render(ctx, io.Discard, benchBinaryData, NoOptions)
	}
}

// ── Pipe ──────────────────────────────────────────────────────────────────────

func BenchmarkPipeJSONGzip(b *testing.B) {
	r := Pipe(JSON(), Gzip())
	ctx := context.Background()
	for b.Loop() {
		_ = r.Render(ctx, io.Discard, benchStructData, NoOptions)
	}
}

// ── Cache ─────────────────────────────────────────────────────────────────────

func BenchmarkCacheHit(b *testing.B) {
	r := Cache(JSON())
	ctx := context.Background()
	// Warm the cache before measuring.
	_ = r.Render(ctx, io.Discard, benchStructData, NoOptions)
	for b.Loop() {
		_ = r.Render(ctx, io.Discard, benchStructData, NoOptions)
	}
}

func BenchmarkCacheParallel(b *testing.B) {
	r := Cache(JSON())
	ctx := context.Background()
	_ = r.Render(ctx, io.Discard, benchStructData, NoOptions)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = r.Render(ctx, io.Discard, benchStructData, NoOptions)
		}
	})
}
