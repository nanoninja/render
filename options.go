// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package render

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/textproto"
	"strings"
	"time"
)

// HeaderOptions is a type alias for textproto.MIMEHeader, providing standard
// HTTP header management with support for multiple values per key.
type HeaderOptions = textproto.MIMEHeader

// Options defines configuration settings for rendering operations.
// It provides control over formatting, content type, encoding and other
// rendering aspects.
type Options struct {
	name    string            // Template name or identifier
	timeout time.Duration     // Maximum duration for rendering
	format  FormatOptions     // Formatting configuration
	header  HeaderOptions     // All headers including content type and charset
	params  map[string]string // Additional parameters
}

// NewOptions creates a new Options instance with default values.
func NewOptions() *Options {
	return (&Options{}).Reset()
}

// Clone creates a deep copy of the Options instance.
// It returns a new Options pointer with all fields deeply copied to ensure
// modifications to the clone do not affect the original:
//   - Simple fields (name, timeout) are copied directly
//   - Format options are cloned using FormatOptions.Clone()
//   - Parameters map is recreated with copied key-value pairs
//   - Headers are recreated with copied values for each key
//
// This method is useful when you need to modify options without affecting
// the original configuration or when sharing options across different renders.
func (o *Options) Clone() *Options {
	clone := &Options{
		name:    o.name,
		timeout: o.timeout,
		format:  o.format.Clone(),
		params:  make(map[string]string, len(o.params)),
	}
	for k, v := range o.params {
		clone.params[k] = v
	}
	clone.header = make(HeaderOptions)
	for k, v := range o.header {
		clone.header[k] = append([]string{}, v...)
	}
	return clone
}

// Name returns the configured template name or identifier.
// An empty string means no specific template name is set.
func (o *Options) Name() string {
	return o.name
}

// ContentType returns the content type from headers, if any.
func (o *Options) ContentType() string {
	return o.header.Get("Content-Type")
}

// Format returns the formatting options configured for this renderer.
// It provides access to formatting settings like:
//   - Indentation
//   - Line endings
//   - Pretty printing configuration
//   - Format arguments for template rendering
//
// The returned FormatOptions controls how the rendered output will be formatted,
// allowing consistent formatting across different render operations.
func (o *Options) Format() FormatOptions {
	return o.format
}

// Header returns the header options for the rendering operation.
// Headers can be used to set content type, charset, and other HTTP headers.
func (o *Options) Header() HeaderOptions {
	return o.header
}

// Params returns the map of additional parameters.
// These parameters can be used to pass custom configuration to renderers.
func (o *Options) Params() map[string]string {
	return o.params
}

// Reset restores all options to their default values.
// It ensures that all fields have appropriate defaults and
// allocates necessary resources like maps.
func (o *Options) Reset() *Options {
	o.name = ""
	o.timeout = 0
	o.format = FormatOptions{
		indent: "",
		pretty: false,
	}
	o.params = make(map[string]string)
	o.header = make(HeaderOptions)
	return o
}

// String returns a human-readable representation of the current options configuration.
// This includes template name, timeout, format settings, headers, and parameters.
func (o *Options) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Template Name: %q\n", o.name))
	b.WriteString(fmt.Sprintf("Timeout: %v\n", o.timeout))

	b.WriteString("Format:\n")
	b.WriteString(fmt.Sprintf("  Pretty: %v\n", o.format.pretty))
	b.WriteString(fmt.Sprintf("  Indent: %q\n", o.format.indent))

	b.WriteString("Headers:\n")
	for key, values := range o.header {
		b.WriteString(fmt.Sprintf("  %s: %s\n", key, strings.Join(values, ",")))
	}
	b.WriteString("Parameters:\n")
	for k, v := range o.params {
		b.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
	}
	return b.String()
}

// Timeout returns the configured rendering timeout duration.
// A zero duration means no timeout is set.
func (o *Options) Timeout() time.Duration {
	return o.timeout
}

// Use applies one or more option functions to the Options instance.
// It validates the options after all functions are applied.
func (o *Options) Use(opts ...func(*Options)) *Options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// CaptureOptions provides a way to access the options used during rendering.
// It returns an option function that captures the fully configured Options object.
//
// Example:
//
//	var opts *Options
//	renderer.Render(w, data, CaptureOptions(&opts))
//	fmt.Println(opts.ContentType())
func CaptureOptions(target **Options) func(*Options) {
	return func(o *Options) { *target = o.Clone() }
}

// Dump returns an option that prints the current rendering configuration
// to the specified writer. This allows for flexible debugging output.
//
// Example:
//
//	// Print options to stdout
//	t.Render(w, data, render.Dump(os.Stdout))
//
//	// Write options to a log file
//	logFile, _ := os.Create("render.log")
//	t.Render(w, data, render.Dump(logFile))
func Dump(w io.Writer) func(*Options) {
	return func(o *Options) {
		fmt.Fprintln(w, "\n=== Template Render Options ===")
		fmt.Fprint(w, o.String())
		fmt.Fprintln(w, "================================")
	}
}

// Header creates an option function that applies multiple header modifications.
// It allows configuring multiple headers in a single option using functional parameters.
func Header(headers ...func(HeaderOptions)) func(*Options) {
	return func(o *Options) {
		for _, header := range headers {
			header(o.header)
		}
	}
}

// Mime returns an option function that sets the Content-Type header with an optional charset.
// It configures the MIME type for the rendered content, which defines how the content
// should be interpreted by clients. The charset parameter specifies the character encoding.
//
// The mediatype parameter should be a valid MIME type like "text/plain" or "application/json".
// If a charset is provided, it's added to the Content-Type header (e.g., "text/plain; charset=utf-8").
//
// Example:
//
//	// Set content type without charset
//	renderer.Render(w, data, Mime("text/plain"))
//
//	// Set content type with charset
//	renderer.Render(w, data, Mime("text/html", "utf-8"))
func Mime(mediatype string, charset ...string) func(*Options) {
	return Header(func(o HeaderOptions) {
		if len(charset) > 0 {
			mediatype = mime.FormatMediaType(mediatype, map[string]string{
				"charset": charset[0],
			})
		}
		o.Set("Content-Type", mediatype)
	})
}

// MimeUTF8 returns an option function that sets the Content-Type header with UTF-8 charset.
// This is a convenience function that combines Mime with the UTF-8 charset, which is
// the most common character encoding for web content. It's equivalent to calling
// Mime(mediatype, "utf-8").
//
// Example:
//
//	// Both lines are equivalent:
//	renderer.Render(w, data, MimeUTF8("text/plain"))
//	renderer.Render(w, data, Mime("text/plain", "utf-8"))
func MimeUTF8(mediatype string) func(*Options) {
	return Mime(mediatype, "utf-8")
}

// Name sets the template name for template renderers.
func Name(name string) func(*Options) {
	return func(o *Options) { o.name = name }
}

// Param adds a custom parameter with the given key and value.
// It's useful for passing additional configuration to renderers.
func Param(key, value string) func(*Options) {
	return func(o *Options) {
		if o.params == nil {
			o.params = make(map[string]string)
		}
		o.params[key] = value
	}
}

// Separator returns an option function that sets the CSV field separator.
// The first character of the string is used as separator.
// Example:
//
//	renderer.Render(w, data, Separator(";"))
func Separator(sep string) func(*Options) {
	return Param("separator", sep)
}

// Timeout sets a timeout duration for the rendering operation.
func Timeout(d time.Duration) func(*Options) {
	return func(o *Options) { o.timeout = d }
}

// With creates a reusable set of options that can be applied together.
// It combines multiple option functions into a single function, making it
// easier to manage and reuse common configuration patterns.
func With(opts ...func(*Options)) func(*Options) {
	return func(o *Options) {
		for _, opt := range opts {
			opt(o)
		}
	}
}

// WriteResponse creates an option function that copies all headers from the Options
// to an http.ResponseWriter. This is useful in HTTP handlers when you need to apply
// the configured headers to the HTTP response.
//
// Example:
//
//	func handleTemplate(w http.ResponseWriter, r *http.Request) {
//	    render.JSON().Render(w, data, render.WriteResponse(w))
//	}
func WriteResponse(w http.ResponseWriter) func(*Options) {
	return func(o *Options) {
		for k, v := range o.header {
			for _, value := range v {
				w.Header().Set(k, value)
			}
		}
	}
}
