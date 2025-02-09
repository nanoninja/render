// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package render

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/nanoninja/assert"
)

var (
	_ Renderer = (*jsonRenderer)(nil)
	_ Renderer = JSON()
	_ Renderer = NewJSON(JSONConfig{})
)

const jsonTest = `{"message":"JSON render test"}`

const jsonPrettyTest = `{
  "message": "JSON render test"
}`

const jsonPrefixTest = `{
>>  "message": "JSON prefix test"
>>}`

func TestJSONRenderer(t *testing.T) {
	t.Run("RendersSimpleData", func(t *testing.T) {
		var w bytes.Buffer

		data := map[string]string{"message": "JSON render test"}

		err := JSON().Render(&w, data)

		assert.Nil(t, err)
		assert.Equals(t, strings.TrimSpace(w.String()), jsonTest)
	})

	t.Run("RendersWithPrettyPrinting", func(t *testing.T) {
		var w bytes.Buffer

		data := map[string]string{"message": "JSON render test"}

		err := JSON().Render(&w, data, Format(Pretty()))

		assert.Nil(t, err)
		assert.Equals(t, strings.TrimSpace(w.String()), jsonPrettyTest)
	})

	t.Run("CustomConfigurationWithJSONP", func(t *testing.T) {
		var w bytes.Buffer

		data := map[string]string{"message": "JSON render test"}
		config := JSONConfig{Padding: "callback"}

		err := NewJSON(config).Render(&w, data)

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "callback({\"message\":\"JSON render test\"}\n)")
	})

	t.Run("CustomIndentAndPrefix", func(t *testing.T) {
		var w bytes.Buffer

		data := map[string]string{"message": "JSON prefix test"}
		config := JSONConfig{
			Prefix: ">>",
			Indent: "  ",
		}

		err := NewJSON(config).Render(&w, data, Format(Pretty()))

		assert.Nil(t, err)
		assert.Equals(t, strings.TrimSpace(w.String()), jsonPrefixTest)
	})

	t.Run("UseConfigIndentWithPretty", func(t *testing.T) {
		var w bytes.Buffer

		data := map[string]string{"key": "value"}
		config := JSONConfig{Indent: "    "}

		err := NewJSON(config).Render(&w, data, Format(Pretty()))

		expected := "{\n    \"key\": \"value\"\n}\n"

		assert.Nil(t, err)
		assert.Equals(t, w.String(), expected)
	})

	t.Run("UseOptionIndentWhenExplicitlySet", func(t *testing.T) {
		var w bytes.Buffer

		data := map[string]string{"key": "value"}
		config := JSONConfig{Indent: "    "}

		err := NewJSON(config).Render(&w, data, Format(
			Pretty(),
			Indent("  "),
		))

		assert.Nil(t, err)
		assert.Equals(t, w.String(), "{\n  \"key\": \"value\"\n}\n")
	})

	t.Run("EscapesHTMLByDefault", func(t *testing.T) {
		var w bytes.Buffer

		data := map[string]string{"html": "<script>alert('test')</script>"}

		err := JSON().Render(&w, data)

		assert.Nil(t, err)
		assert.StringContains(t, w.String(), `\u003cscript\u003e`)
	})

	t.Run("CanDisableHTMLEscaping", func(t *testing.T) {
		var w bytes.Buffer

		data := map[string]string{"html": "<script>alert('test')</script>"}
		config := JSONConfig{EscapeHTML: false}

		err := NewJSON(config).Render(&w, data)

		assert.Nil(t, err)
		assert.StringContains(t, w.String(), "<script>")
	})

	t.Run("SetsDefaultContentType", func(t *testing.T) {
		var w bytes.Buffer
		var opts *Options

		err := JSON().Render(&w, "test", CaptureOptions(&opts))

		assert.Nil(t, err)
		assert.Equals(t, opts.ContentType(), "application/json; charset=utf-8")
	})

	t.Run("RespectsContextCancellation", func(t *testing.T) {
		var w bytes.Buffer

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := JSON().RenderContext(ctx, &w, "tests")

		assert.NotNil(t, err)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("HandlesWriterError", func(t *testing.T) {
		config := JSONConfig{Padding: "callback"}

		err := NewJSON(config).Render(&errorWriterTest{}, nil)

		assert.NotNil(t, err)
	})
}
