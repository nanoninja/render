// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render_test

import (
	"bytes"
	"context"
	"os"

	"github.com/nanoninja/render"
)

func ExampleBinary() {
	_ = render.Binary().Render(context.Background(), os.Stdout, []byte("binary content"), render.NoOptions)

	// Output:
	// binary content
}

func ExampleText() {
	_ = render.Text().Render(context.Background(), os.Stdout, "Hello, Gopher!", render.NoOptions)

	// Output:
	// Hello, Gopher!
}

func ExampleJSON() {
	data := map[string]string{"message": "hello"}

	_ = render.JSON().Render(context.Background(), os.Stdout, data, render.NoOptions)

	// Output:
	// {"message":"hello"}
}

func ExampleYAML() {
	data := map[string]string{"message": "hello"}

	_ = render.YAML().Render(context.Background(), os.Stdout, data, render.NoOptions)

	// Output:
	// message: hello
}

func ExampleHTML() {
	_ = render.HTML().Render(context.Background(), os.Stdout, "<h1>Hello</h1>", render.NoOptions)

	// Output:
	// <h1>Hello</h1>
}

func ExampleMulti() {
	var log bytes.Buffer

	_ = render.Multi(render.JSON(), &log).Render(context.Background(), os.Stdout, map[string]string{
		"key": "value",
	}, render.NoOptions)

	// Output:
	// {"key":"value"}
}

func ExamplePipe() {
	_ = render.Pipe(
		render.Text(),
	).Render(context.Background(), os.Stdout, "Hello", render.NoOptions)

	// Output:
	// Hello
}

func ExampleMarkdown() {
	_ = render.Markdown().Render(context.Background(), os.Stdout, "# Hello", render.NoOptions)

	// Output:
	// <h1>Hello</h1>
}

func ExampleGzip() {
	var buf bytes.Buffer

	_ = render.Pipe(
		render.JSON(),
		render.Gzip(),
	).Render(context.Background(), &buf, map[string]string{"message": "hello"}, render.NoOptions)
}

func ExampleCache() {
	r := render.Cache(render.JSON())

	_ = r.Render(context.Background(), os.Stdout, "hello", render.NoOptions)
	_ = r.Render(context.Background(), os.Stdout, "hello", render.NoOptions)

	// Output:
	// "hello"
	// "hello"
}

func ExampleNewCache() {
	r := render.NewCache(render.JSON(), render.CacheConfig{
		TTL: 5 * 60 * 1000000000,
	})

	_ = r.Render(context.Background(), os.Stdout, map[string]string{"status": "ok"}, render.NoOptions)

	// Output:
	// {"status":"ok"}
}
