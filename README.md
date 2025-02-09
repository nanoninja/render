# Render

A flexible Go package for rendering content in different formats with configurable options and formatting.

[![Golang](https://img.shields.io/badge/Go-%3E%3D%201.18-%2300ADD8.svg)](https://go.dev/)
[![Tests](https://github.com/nanoninja/render/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/nanoninja/render/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nanoninja/render/branch/main/graph/badge.svg)](https://codecov.io/gh/nanoninja/render)
[![Go Report Card](https://goreportcard.com/badge/github.com/nanoninja/render)](https://goreportcard.com/report/github.com/nanoninja/render)
[![GoDoc](https://godoc.org/github.com/nanoninja/render?status.svg)](https://godoc.org/github.com/nanoninja/render)
[![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

## Features

- Multiple output formats (JSON, XML, text)
- Pretty printing and custom formatting
- Context support with cancellation
- Buffered rendering with post-processing
- Content type handling

## Installation

```bash
go get github.com/nanoninja/render
```

## Quick Start

```go
// Create renderer
renderer := render.JSON()

data := map[string]string{"message": "ping"}

// Simple render
renderer.Render(os.Stdout, data)

// Pretty printed with custom indent
renderer.Render(os.Stdout, data, render.Format(
    render.Pretty(),
    render.Indent("    "),
))
```

## Text Rendering

```go
// Simple text
render.Text().Render(os.Stdout, "Hello World")

// With formatting
render.Text().Render(os.Stdout, "Hello %s", render.Textf("Gopher"))
```

## Buffered Rendering

```go
// Create buffered JSON renderer
render.Buffer(render.JSON()).Render(os.Stdout, data)
```

## Template Rendering

```go
package main

import (
	"embed"
	"net/http"
	"strings"
	"text/template"

	"github.com/nanoninja/render"
	"github.com/nanoninja/render/tmpl"
	"github.com/nanoninja/render/tmpl/loader"
)

//go:embed templates
var templatesFS embed.FS

func main() {
	// Create a loader for your templates
	src := loader.NewEmbed(templatesFS, tmpl.LoaderConfig{
		Root:      "templates",
		Extension: ".html",
	})

	// Create template with configuration
	t := tmpl.HTML("",
		tmpl.SetFuncsHTML(template.FuncMap{
			"upper": strings.ToUpper,
		}),
		tmpl.LoadHTML(src),
	)

	// Use in an HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{"Title": "Welcome"}

		t.RenderContext(r.Context(), w, data,
			render.Name("index.html"), // Specify which template to render
			render.WriteResponse(w),   // Write headers to response
		)
	})

	http.ListenAndServe("localhost:8080", nil)
}

```

## Context Support

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

renderer.RenderContext(ctx, w, data)
```

## Configuration

### Format Options

Options can be configured using functional options:

```go
// Pretty printing with custom format
renderer.Render(w, data, render.Format(
    render.Pretty(),          // Enable pretty printing
    render.Indent("    "),    // Custom indentation
    render.LineEnding("\n"),  // Line ending style
    render.Prefix("// ")      // Line prefix
))

// CSV specific options
renderer.Render(w, data, render.Format(
    render.Separator(";"),    // Use semicolon as separator
    render.UseCRLF(),         // Use Windows-style line endings
))
```

### Other Options

```go
// Set template name (for template renderers)
renderer.Render(w, data, render.Name("index.html"))

// Set render timeout
renderer.Render(w, data, render.Timeout(5*time.Second))

// Custom parameters
renderer.Render(w, data, render.Param("key", "value"))

// Combining multiple options
renderer.Render(w, data, render.With(
    render.MimeJSON(),
    render.Format(render.Pretty()),
    render.Timeout(5*time.Second),
))
```

### In HTTP Context

```go
func handleJSON(w http.ResponseWriter, r *http.Request) {
    data := struct {
        Message string `json:"message"`
    }{
        Message: "Hello, World!",
    }

    renderer := render.JSON()
    renderer.Render(w, data,
        render.MimeJSON(),             // Set Content-Type
        render.WriteResponse(w),       // Write headers to response
        render.Format(render.Pretty()) // Pretty print output
    )
}
```

## Creating Custom Renderers

You can create your own renderers to support any output format. Here's a complete guide to implementing a custom renderer.

### Basic Implementation

Your renderer must implement the `Renderer` interface. Here's a minimal example:

```go
type CustomRenderer struct {}

func (r *CustomRenderer) Render(w io.Writer, data any, opts func(*render.Options)) error {
    return r.RenderContext(context.Background(), w, data, opts...)
}

func (r *CustomRenderer) RenderContext(ctx context.Context, w io.Writer, data any, opts ...func(*render.Options)) error {
    // Check context validity
    if err := render.CheckContext(ctx); err != nil {
        return err
    }

    // Get options with default values
    opt := render.NewOptions(opts...)

    // Your rendering logic here
    return nil
}
```

### Available Utilities

- `NewOptions`: Processes option functions and returns configured Options
- `CheckContext`: Verifies if the context is still valid
Common option handlers for content type, formatting, etc.


## License

This project is licensed under the BSD 3-Clause License.

It allows you to:
- Use the software commercially
- Modify the software
- Distribute the software
- Place warranty on the software
- Use the software privately

The only requirements are:
- Include the copyright notice
- Include the license text
- Not use the author's name to promote derived products without permission

For more details, see the LICENSE file in the project repository.