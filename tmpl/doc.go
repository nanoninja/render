// Package tmpl provides Go template rendering on top of the render.Renderer interface.
//
// Two template types are available:
//
//   - [HTML] — wraps html/template with automatic XSS escaping, recommended for web pages.
//   - [Text] — wraps text/template, suitable for emails, CLI output, and plain text.
//
// Both implement [render.Renderer] and share the same loading and function-map API.
//
// Basic usage:
//
//	t := tmpl.HTML("myapp")
//
//	if err := t.Load(src); err != nil {
//	    log.Fatal(err)
//	}
//
//	err := t.Render(ctx, w, data, render.Options{Name: "index.html"})
//
// Use [SetFuncsHTML] / [SetFuncs] to inject custom template functions,
// or [WithDefaultFuncsHTML] / [WithDefaultFuncs] for the built-in helper set
// (upper, lower, date, add, nl2br, …).
package tmpl
