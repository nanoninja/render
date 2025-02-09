// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package render

import (
	"bytes"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nanoninja/assert"
)

func TestOptions(t *testing.T) {
	t.Run("NewOptionsReturnsInitializedInstance", func(t *testing.T) {
		opts := NewOptions()

		assert.Equals(t, opts.Name(), "")
		assert.Equals(t, opts.Timeout(), time.Duration(0))
		assert.Equals(t, opts.Format().Indent(), "")
		assert.False(t, opts.Format().Pretty())
		assert.NotNil(t, opts.Params())
		assert.NotNil(t, opts.Header())
	})

	t.Run("CloneCreatesDeepCopy", func(t *testing.T) {
		original := NewOptions()
		original.name = "test"
		original.timeout = 5 * time.Second
		original.format.pretty = true
		original.params["key"] = "value"
		original.header.Set("Content-Type", "text/plain")

		clone := original.Clone()

		assert.Equals(t, clone.Name(), original.Name())
		assert.Equals(t, clone.timeout, original.timeout)
		assert.Equals(t, clone.Format().Pretty(), original.Format().Pretty())

		clone.params["newKey"] = "newValue"
		assert.NotEquals(t, len(clone.params), len(original.Params()))

		clone.header.Set("New-Header", "value")
		assert.NotEquals(t, len(clone.Header()), len(original.Header()))
	})

	t.Run("NameReturnsConfiguredValue", func(t *testing.T) {
		opts := NewOptions()
		opts.name = "template.tmpl"

		assert.Equals(t, opts.Name(), "template.tmpl")
	})

	t.Run("ContentTypeReturnsHeaderValue", func(t *testing.T) {
		opts := NewOptions()
		opts.header.Set("Content-Type", "text/plain")

		assert.Equals(t, opts.ContentType(), "text/plain")
	})

	t.Run("FormatReturnsFormatOptions", func(t *testing.T) {
		opts := NewOptions()
		opts.format.pretty = true
		opts.format.indent = "    "

		format := opts.Format()

		assert.True(t, format.Pretty())
		assert.Equals(t, format.Indent(), "    ")
	})

	t.Run("ParamsReturnsParmetersMap", func(t *testing.T) {
		opts := NewOptions()
		opts.params["key1"] = "value1"
		opts.params["key2"] = "value2"

		params := opts.Params()

		assert.Len(t, params, 2)
		assert.Equals(t, params["key1"], "value1")
		assert.Equals(t, params["key2"], "value2")
	})

	t.Run("HeaderReturnsHeaderOptions", func(t *testing.T) {
		opts := NewOptions()
		opts.header.Set("X-Test", "value")

		headers := opts.Header()

		assert.Equals(t, headers.Get("X-Test"), "value")
	})

	t.Run("ResetRestoresDefaultValues", func(t *testing.T) {
		opts := NewOptions()

		opts.name = "test"
		opts.timeout = 5 * time.Second
		opts.format.pretty = true
		opts.params["key"] = "value"

		opts.Reset()

		assert.Equals(t, opts.Name(), "")
		assert.Equals(t, opts.Timeout(), time.Duration(0))
		assert.Equals(t, opts.Format().Indent(), "")
		assert.False(t, opts.Format().Pretty())
		assert.Len(t, opts.Params(), 0)
		assert.Len(t, opts.Header(), 0)
	})

	t.Run("TimeoutReturnsConfiguredDuration", func(t *testing.T) {
		opts := NewOptions()
		opts.timeout = 5 * time.Second
		assert.Equals(t, opts.Timeout(), 5*time.Second)
	})

	t.Run("UseAppliesOptionFunctions", func(t *testing.T) {
		opts := NewOptions()

		opts.Use(func(o *Options) {
			o.name = "test"
			o.timeout = 3 * time.Second
			o.params["key"] = "value"
		})

		assert.Equals(t, opts.Name(), "test")
		assert.Equals(t, opts.Timeout(), 3*time.Second)
		assert.Equals(t, opts.Params()["key"], "value")
	})

	t.Run("UseAllowsChaining", func(t *testing.T) {
		opts := NewOptions()

		result := opts.Use(
			func(o *Options) { o.name = "test1" },
		).Use(
			func(o *Options) { o.timeout = time.Second },
		)

		assert.Equals(t, result.Name(), "test1")
		assert.Equals(t, result.Timeout(), time.Second)
	})
}

func TestOptionFunctions(t *testing.T) {
	t.Run("CaptureOptionsClonesOptions", func(t *testing.T) {
		original := NewOptions()
		original.name = "test"
		var captured *Options

		CaptureOptions(&captured)(original)

		assert.NotNil(t, captured)
		assert.Equals(t, captured.Name(), original.Name())
	})

	t.Run("StringProvidesFormattedRepresentation", func(t *testing.T) {
		opts := NewOptions()
		opts.name = "template.tmpl"
		opts.timeout = 5 * time.Second
		opts.format.pretty = true
		opts.format.indent = "    "
		opts.header.Set("Content-Type", "text/plain")
		opts.params["key"] = "value"

		expected := `Template Name: "template.tmpl"
Timeout: 5s
Format:
  Pretty: true
  Indent: "    "
Headers:
  Content-Type: text/plain
Parameters:
  key: value
`
		assert.Equals(t, opts.String(), expected)
	})

	t.Run("DumpWritesFormattedOptionsToWriter", func(t *testing.T) {
		opts := NewOptions()
		opts.name = "template.tmpl"
		opts.timeout = 5 * time.Second
		opts.format.pretty = true
		opts.format.indent = "    "
		opts.header.Set("Content-Type", "text/plain")
		opts.params["key"] = "value"

		buf := new(bytes.Buffer)

		Dump(buf)(opts)

		expected := `
=== Template Render Options ===
Template Name: "template.tmpl"
Timeout: 5s
Format:
  Pretty: true
  Indent: "    "
Headers:
  Content-Type: text/plain
Parameters:
  key: value
================================
`
		assert.Equals(t, buf.String(), expected)
	})

	t.Run("DumpHandlesEmptyOptions", func(t *testing.T) {
		opts := NewOptions()
		buf := new(bytes.Buffer)

		Dump(buf)(opts)

		expected := `
=== Template Render Options ===
Template Name: ""
Timeout: 0s
Format:
  Pretty: false
  Indent: ""
Headers:
Parameters:
================================
`
		assert.Equals(t, buf.String(), expected)
	})

	t.Run("MimeSetsContentType", func(t *testing.T) {
		opts := NewOptions()

		Mime("text/plain", "utf-8")(opts)

		assert.Equals(t, opts.ContentType(), "text/plain; charset=utf-8")
	})

	t.Run("MimeUTF9SetsContentTypeWithUTF8", func(t *testing.T) {
		opts := NewOptions()

		MimeUTF8("text/plain")(opts)

		assert.Equals(t, opts.ContentType(), "text/plain; charset=utf-8")
	})

	t.Run("NameSetTemplateName", func(t *testing.T) {
		opts := NewOptions()

		Name("template.tmpl")(opts)

		assert.Equals(t, opts.Name(), "template.tmpl")
	})

	t.Run("ParamAddsCustomParameter", func(t *testing.T) {
		opts := NewOptions()

		Param("key", "value")(opts)

		assert.Equals(t, opts.Params()["key"], "value")
	})

	t.Run("ParamInitializedMapIfNil", func(t *testing.T) {
		opts := &Options{}

		Param("key", "value")(opts)

		assert.NotNil(t, opts.Params())
		assert.Equals(t, opts.Params()["key"], "value")
	})

	t.Run("SeparatorSetsSeparatorParameter", func(t *testing.T) {
		opts := NewOptions()

		Separator(";")(opts)

		assert.Equals(t, opts.Params()["separator"], ";")
	})

	t.Run("TimeoutSetsRenderTimeout", func(t *testing.T) {
		opts := NewOptions()

		Timeout(5 * time.Second)(opts)

		assert.Equals(t, opts.Timeout(), 5*time.Second)
	})

	t.Run("WithCombinedsMultipleOptions", func(t *testing.T) {
		opts := NewOptions()
		combinedOpt := With(
			Name("test"),
			Timeout(5*time.Second),
			Param("key", "value"),
		)
		combinedOpt(opts)

		assert.Equals(t, opts.Name(), "test")
		assert.Equals(t, opts.Timeout(), 5*time.Second)
		assert.Equals(t, opts.Params()["key"], "value")
	})

	t.Run("WriteResponseSetsResponseHeaders", func(t *testing.T) {
		opts := NewOptions()
		opts.header.Set("Content-Type", "text/plain")
		opts.header.Set("X-Test", "value")

		recorder := httptest.NewRecorder()

		WriteResponse(recorder)(opts)

		assert.Equals(t, recorder.Header().Get("Content-Type"), "text/plain")
		assert.Equals(t, recorder.Header().Get("X-Test"), "value")
	})
}
