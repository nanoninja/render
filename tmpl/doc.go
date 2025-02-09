// Package tmpl provides a flexible template rendering system that supports both
// text and HTML templates. It builds upon Go's standard template packages while
// adding consistent rendering options, loader abstraction, and proper header handling.
//
// Basic usage:
//
//	// Create and configure a text template
//	t := tmpl.Text("example",
//	    tmpl.Load(loader),
//	    tmpl.SetFuncs(funcMap),
//	)
//
//	// Render with options
//	err := t.RenderContext(ctx, w, data,
//	    render.Name("content.txt"),
//	    render.WriteResponse(w),
//	)
//
// The package provides two main template types:
//
//   - TextTemplate: For plain text templates without HTML escaping
//   - HTMLTemplate: For HTML templates with proper escaping and content type
//
// Both types share common functionality through the Base type while providing
// type-specific behavior and proper content type handling.
package tmpl
