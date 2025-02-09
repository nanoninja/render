// doc.go
// Package render provides a flexible and extensible rendering system for Go applications.
//
// This package is designed to handle various output formats (JSON, Text, HTML templates, etc.)
// in a consistent way, while maintaining type safety and providing rich configuration options.
//
// Basic usage:
//
//	// Create a renderer
//	renderer := render.JSON()
//
//	// Render with options
//	err := renderer.Render(w, data,
//	    render.Format(render.Pretty()),
//	)
//
// The package supports multiple renderer types:
//   - Text rendering for plain text output
//   - JSON rendering with support for pretty printing and JSONP
//   - HTML template rendering
//   - And more...
//
// Key features:
//   - Multiple output formats (JSON, XML, text)
//   - Configurable formatting (indentation, pretty printing)
//   - Context-aware rendering with cancellation support
//   - Buffered rendering with post-processing
//   - Content type handling
//
// Each renderer implements the Renderer interface:
//
//	type Renderer interface {
//	    Render(w io.Writer, data any, opts ...func(*Options)) error
//	    RenderContext(ctx context.Context, w io.Writer, data any, opts ...func(*Options)) error
//	}
//
// The package uses the options pattern for configuration, making it easy to extend
// and customize rendering behavior without breaking existing code.
package render
