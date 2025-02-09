// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package render

// FormatOptions defines common formatting settings for all renderers.
// It provides basic text formatting capabilities that can be used
// across different types of output (JSON, XML, Text, etc.).
type FormatOptions struct {
	// prefix is added at the beginning of each line.
	// It can be used for comments, indentation, or any line-starting content.
	prefix string

	// lineEnding specifies the character(s) used for line termination.
	// Common values are "\n" (Unix) or "\r\n" (Windows).
	// If empty, the renderer will use its default line ending.
	lineEnding string

	// indent specifies the string used for indentation levels.
	// Typically spaces or tabs, e.g., "  " for 2-space indent.
	indent string

	// pretty enables formatted output with proper spacing and line breaks.
	// When true, the output will be human-readable rather than compact.
	pretty bool

	// args contains format arguments used with Printf-style formatting.
	// For example, with format string "Hello, %s", args would contain ["World"].
	args []any
}

// Clone creates a deep copy of FormatOptions.
// It returns a new FormatOptions instance with all fields copied from the original:
// - prefix for line starting content
// - lineEnding for line termination characters
// - indent for indentation string
// - pretty flag for human-readable output
// - args for format arguments
// This ensures that modifications to the clone don't affect the original options.
func (o FormatOptions) Clone() FormatOptions {
	return FormatOptions{
		prefix:     o.prefix,
		lineEnding: o.lineEnding,
		indent:     o.indent,
		pretty:     o.pretty,
		args:       append([]any(nil), o.args...),
	}
}

// Indent returns the string used for indentation.
func (o FormatOptions) Indent() string {
	return o.indent
}

// LineEnding returns the string used for line endings.
func (o FormatOptions) LineEnding() string {
	if o.lineEnding == "" {
		return "\n"
	}
	return o.lineEnding
}

// Prefix returns the string to prepend to each line.
func (o FormatOptions) Prefix() string {
	return o.prefix
}

// Pretty returns whether pretty formatting is enabled.
func (o FormatOptions) Pretty() bool {
	return o.pretty
}

// Args sets format arguments for text formatting operations.
// Example: Format(Args("Gopher"), Pretty())
func Args(args ...any) func(*FormatOptions) {
	return func(o *FormatOptions) {
		o.args = args
	}
}

// Format applies formatting options to the render operation.
func Format(formatters ...func(*FormatOptions)) func(*Options) {
	return func(o *Options) {
		for _, formatter := range formatters {
			formatter(&o.format)
		}
	}
}

// Comment adds a comment marker to each line.
func Comment(marker string) func(*FormatOptions) {
	return Prefix(marker + " ")
}

// Pretty enables pretty printing with default settings.
func Pretty() func(*FormatOptions) {
	return func(f *FormatOptions) { f.pretty = true }
}

// Indent sets the indentation string for formatting.
func Indent(indent string) func(*FormatOptions) {
	return func(f *FormatOptions) { f.indent = indent }
}

// LineEnding sets the line ending style.
func LineEnding(ending string) func(*FormatOptions) {
	return func(f *FormatOptions) { f.lineEnding = ending }
}

// Prefix sets the prefix string for formatted output.
func Prefix(prefix string) func(*FormatOptions) {
	return func(f *FormatOptions) { f.prefix = prefix }
}

// Textf provides a convenient way to set text format arguments.
// It internally uses Format and Args functions to maintain consistency
// with the existing formatting system.
// Example:
//
//	txt := render.Text()
//	txt.Render(os.Stdout, "hello %s", render.Textf("Gopher"))
func Textf(args ...any) func(*Options) {
	return Format(Args(args...))
}

// UseCRLF returns an option function that sets line endings to CRLF (\r\n).
// It combines Format and LineEnding functions to provide a convenient way
// to configure CRLF line endings, commonly used in CSV files for Windows
// compatibility and RFC 4180 compliance.
//
// Example:
//
//	renderer.Render(w, data, UseCRLF())
//
// This is equivalent to:
//
//	renderer.Render(w, data, Format(LineEnding("\r\n")))
func UseCRLF() func(*Options) {
	return Format(LineEnding("\r\n"))
}
