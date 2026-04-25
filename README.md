# Render

A lightweight Go package for rendering content in multiple formats.  
Zero allocations on the hot path — configuration is a plain value struct, not closures.

[![Golang](https://img.shields.io/badge/Go-%3E%3D%201.24-%2300ADD8.svg)](https://go.dev/)
[![Tests](https://github.com/nanoninja/render/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/nanoninja/render/actions/workflows/ci.yaml)
[![codecov](https://codecov.io/gh/nanoninja/render/graph/badge.svg?token=FNR16B9429)](https://codecov.io/gh/nanoninja/render)
[![Go Report Card](https://goreportcard.com/badge/github.com/nanoninja/render)](https://goreportcard.com/report/github.com/nanoninja/render)
[![Go Reference](https://pkg.go.dev/badge/github.com/nanoninja/render.svg)](https://pkg.go.dev/github.com/nanoninja/render)
[![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

## Features

- Built-in renderers: JSON, XML, CSV, Text, HTML, Binary, YAML, Markdown, Gzip
- Composable renderers: Buffer, Cache, Multi, Pipe
- Template rendering via the `tmpl` sub-package (`html/template` and `text/template`)
- Context support — cancellation checked on every render call; timeout supported by composable renderers (Buffer, Cache, Pipe, Multi)
- Zero-allocation options — `Options` is a plain struct passed by value
- Thread-safe by design — renderers hold no mutable state after construction

## Installation

```bash
go get github.com/nanoninja/render
```

## Interface

All renderers implement a two-method interface:

```go
type Renderer interface {
    ContentType() string
    Render(ctx context.Context, w io.Writer, data any, opts Options) error
}
```

`Options` is a plain value struct — set only the fields you need:

```go
type Options struct {
    Timeout    time.Duration // zero means no limit
    Name       string        // template name (tmpl renderers)
    Pretty     bool          // enable indented output
    Indent     string        // indentation string (JSON, XML)
    Prefix     string        // per-line prefix (JSON, XML)
    LineEnding string        // CSV line ending ("\n" or "\r\n")
    Args       []any         // printf-style args (Text renderer)
    Separator  rune          // CSV field separator (default ',')
    Padding    string        // JSONP callback name (JSON renderer)
    Headers    map[string][]string
}

var NoOptions = Options{} // zero value, no allocation
```

## Quick Start

```go
// Prefer structs over maps — structs encode with 1 alloc, maps require
// reflect.MapIter and allocate one value per key (9 allocs for 4 fields).
type Response struct {
    Message string `json:"message"`
}

ctx := context.Background()
data := Response{Message: "hello"}

// Compact JSON
render.JSON().Render(ctx, os.Stdout, data, render.NoOptions)

// Pretty-printed JSON
render.JSON().Render(ctx, os.Stdout, data, render.Options{Pretty: true})
```

## JSON

```go
// Default (HTML-safe, 2-space indent when pretty)
render.JSON().Render(ctx, w, data, render.NoOptions)

// Pretty with custom indent
render.JSON().Render(ctx, w, data, render.Options{Pretty: true, Indent: "    "})

// JSONP
render.NewJSON(render.JSONConfig{Padding: "callback"}).Render(ctx, w, data, render.NoOptions)

// Disable HTML escaping
render.NewJSON(render.JSONConfig{EscapeHTML: false}).Render(ctx, w, data, render.NoOptions)
```

## XML

```go
// Default (with XML header, 2-space indent when pretty)
render.XML().Render(ctx, w, data, render.NoOptions)

// Pretty with custom indent and prefix
render.XML().Render(ctx, w, data, render.Options{Pretty: true, Indent: "    ", Prefix: "//"})

// Without XML declaration header
render.NewXML(render.XMLConfig{Header: false}).Render(ctx, w, data, render.NoOptions)
```

## CSV

```go
data := [][]string{
    {"name", "age"},
    {"Alice", "30"},
}

// Default (comma-separated, LF line endings)
render.CSV().Render(ctx, w, data, render.NoOptions)

// Custom delimiter
render.CSV().Render(ctx, w, data, render.Options{Separator: ';'})

// Windows line endings
render.CSV().Render(ctx, w, data, render.Options{LineEnding: "\r\n"})
```

## Text

```go
// Plain string
render.Text().Render(ctx, w, "Hello, World!", render.NoOptions)

// With printf-style arguments
render.Text().Render(ctx, w, "Hello, %s!", render.Options{Args: []any{"Gopher"}})

// With trailing newline
render.Text().Render(ctx, w, "Hello", render.Options{Pretty: true})
```

## Binary

Accepts `[]byte` or `io.Reader`.

```go
render.Binary().Render(ctx, w, fileBytes, render.NoOptions)
render.Binary().Render(ctx, w, file, render.NoOptions)
```

## YAML

```go
// Default (2-space indent)
render.YAML().Render(ctx, w, data, render.NoOptions)

// Custom indent
render.NewYAML(render.YAMLConfig{Indent: 4}).Render(ctx, w, data, render.NoOptions)
```

## Markdown

Converts Markdown to HTML using [goldmark](https://github.com/yuin/goldmark).

```go
// Standard CommonMark
render.Markdown().Render(ctx, w, "# Hello", render.NoOptions)

// GitHub Flavored Markdown (tables, strikethrough, task lists)
render.NewMarkdown(render.MarkdownConfig{GFM: true}).Render(ctx, w, content, render.NoOptions)

// Allow raw HTML (disabled by default for XSS safety)
render.NewMarkdown(render.MarkdownConfig{Unsafe: true}).Render(ctx, w, content, render.NoOptions)
```

## Gzip

Compresses `[]byte` input. Typically the last step in a `Pipe`.

```go
render.Pipe(
    render.JSON(),
    render.Gzip(),
).Render(ctx, w, data, render.NoOptions)
```

## Multi

Writes the same rendered output to multiple destinations simultaneously.

```go
var log bytes.Buffer

render.Multi(render.JSON(), &log).Render(ctx, w, data, render.NoOptions)
```

## Pipe

Chains renderers — the output of each becomes the input of the next.

```go
render.Pipe(
    render.JSON(),
    render.Gzip(),
).Render(ctx, w, data, render.NoOptions)
```

## Buffer

Renders to an in-memory buffer first, then writes to the destination. Supports post-processing.

```go
// Simple buffered render
render.Buffer(render.JSON()).Render(ctx, w, data, render.NoOptions)

// With post-processing
render.NewBuffer(render.JSON(), render.BufferConfig{
    PostRender: func(b []byte) ([]byte, error) {
        return bytes.ToUpper(b), nil
    },
}).Render(ctx, w, data, render.NoOptions)
```

## Cache

Renders once, stores the output in memory, and replays it on subsequent calls.  
Safe for concurrent use.

```go
// Permanent cache — renders once, reused forever
var countries = render.Cache(render.JSON())

// Cache with TTL — refreshed after expiry
var catalog = render.NewCache(render.JSON(), render.CacheConfig{
    TTL: 5 * time.Minute,
})
```

> **Warning:** the cache stores rendered bytes, not the data. Never cache user-specific or permission-restricted responses.

## Template Rendering

The `tmpl` sub-package wraps Go's standard template packages:

- `tmpl.HTML` — `html/template` with automatic XSS escaping, for web pages
- `tmpl.Text` — `text/template`, for emails, CLI output, or plain text

```go
import (
    "github.com/nanoninja/render"
    "github.com/nanoninja/render/tmpl"
    "github.com/nanoninja/render/tmpl/loader"
)

t := tmpl.HTML("myapp", tmpl.WithDefaultFuncsHTML())

if err := t.Load(src); err != nil {
    log.Fatal(err)
}

// Render a named template
t.Render(ctx, w, data, render.Options{Name: "index.html"})
```

### Built-in template functions

`WithDefaultFuncsHTML` / `WithDefaultFuncs` provide:
`upper`, `lower`, `trim`, `truncate`, `replace`, `split`, `join`,
`hasPrefix`, `hasSuffix`, `default`, `ternary`, `first`, `last`,
`int`, `float`, `now`, `date`, `add`, `sub`, `mul`, `div`, `sum`, `avg`, `nl2br`

## Context and Timeout

Every `Render` call accepts a `context.Context`. Cancellation is checked before rendering begins in all renderers.

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

render.JSON().Render(ctx, w, data, render.NoOptions)
```

Timeout is enforced by composable renderers (`Buffer`, `Cache`, `Pipe`, `Multi`, `Markdown`) via `opts.Timeout`, or by passing a context with a deadline:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

render.Buffer(render.JSON()).Render(ctx, w, data, render.NoOptions)

// Or set the timeout inside Options (composable renderers only)
render.Buffer(render.JSON()).Render(context.Background(), w, data, render.Options{
    Timeout: 5 * time.Second,
})
```

## Custom Renderer

Implement two methods:

```go
type MyRenderer struct{}

func (r *MyRenderer) ContentType() string { return "application/x-custom" }

func (r *MyRenderer) Render(ctx context.Context, w io.Writer, data any, opts render.Options) error {
    if err := render.CheckContext(ctx); err != nil {
        return err
    }
    _, cancel := render.ApplyTimeout(ctx, opts)
    defer cancel()

    _, err := io.WriteString(w, fmt.Sprint(data))
    return err
}
```

## Benchmarks

Run locally:

```bash
go test -bench=. -benchmem ./...
```

Results on Apple M4, Go 1.26 (`-count=3`, struct payload, renderer pre-created):

| Renderer | ns/op | B/op | allocs/op |
|---|---|---|---|
| Text | 7 | 0 | 0 |
| HTML | 7 | 0 | 0 |
| Binary | 16 | 24 | 1 |
| JSON | 110 | 48 | 1 |
| JSON (parallel) | 31 | 48 | 1 |
| XML | 707 | 4 512 | 8 |
| CSV | 390 | 4 096 | 1 |
| YAML | 2 060 | 6 912 | 32 |
| Buffer(JSON) | 161 | 176 | 3 |
| Cache hit | 30 | 48 | 1 |
| HTML template | 630 | 384 | 13 |
| HTML template (parallel) | 215 | 385 | 13 |

The single alloc in JSON is `json.NewEncoder`. Text, HTML, and Binary are zero-allocation on the hot path. HTML templates use `sync.Pool` on clones — parallel throughput scales linearly.

## License

This project is licensed under the BSD 3-Clause License.
See the [LICENSE](LICENSE) file for details.
