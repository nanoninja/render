// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

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

		err := JSON().Render(context.Background(), &w, data, NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, jsonTest, strings.TrimSpace(w.String()))
	})

	t.Run("RendersWithPrettyPrinting", func(t *testing.T) {
		var w bytes.Buffer

		data := map[string]string{"message": "JSON render test"}

		err := JSON().Render(context.Background(), &w, data, Options{Pretty: true})

		assert.NoError(t, err)
		assert.Equal(t, jsonPrettyTest, strings.TrimSpace(w.String()))
	})

	t.Run("CustomConfigurationWithJSONP", func(t *testing.T) {
		var w bytes.Buffer

		data := map[string]string{"message": "JSON render test"}
		config := JSONConfig{Padding: "callback"}

		err := NewJSON(config).Render(context.Background(), &w, data, NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, `callback({"message":"JSON render test"}
)`, w.String())
	})

	t.Run("CustomIndentAndPrefix", func(t *testing.T) {
		var w bytes.Buffer

		data := map[string]string{"message": "JSON prefix test"}
		config := JSONConfig{
			Prefix: ">>",
			Indent: "  ",
		}

		err := NewJSON(config).Render(context.Background(), &w, data, Options{Pretty: true})

		assert.NoError(t, err)
		assert.Equal(t, jsonPrefixTest, strings.TrimSpace(w.String()))
	})

	t.Run("UseConfigIndentWithPretty", func(t *testing.T) {
		var w bytes.Buffer

		data := map[string]string{"key": "value"}
		config := JSONConfig{Indent: "    "}

		err := NewJSON(config).Render(context.Background(), &w, data, Options{Pretty: true})

		expected := "{\n    \"key\": \"value\"\n}\n"

		assert.NoError(t, err)
		assert.Equal(t, expected, w.String())
	})

	t.Run("UseOptionIndentWhenExplicitlySet", func(t *testing.T) {
		var w bytes.Buffer

		data := map[string]string{"key": "value"}
		config := JSONConfig{Indent: "    "}

		err := NewJSON(config).Render(context.Background(), &w, data, Options{Pretty: true, Indent: "  "})

		assert.NoError(t, err)
		assert.Equal(t, "{\n  \"key\": \"value\"\n}\n", w.String())
	})

	t.Run("EscapesHTMLByDefault", func(t *testing.T) {
		var w bytes.Buffer

		data := map[string]string{"html": "<script>alert('test')</script>"}

		err := JSON().Render(context.Background(), &w, data, NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, w.String(), `\u003cscript\u003e`)
	})

	t.Run("CanDisableHTMLEscaping", func(t *testing.T) {
		var w bytes.Buffer

		data := map[string]string{"html": "<script>alert('test')</script>"}
		config := JSONConfig{EscapeHTML: false}

		err := NewJSON(config).Render(context.Background(), &w, data, NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, w.String(), "<script>")
	})

	t.Run("SetsDefaultContentType", func(t *testing.T) {
		assert.Equal(t, "application/json; charset=utf-8", JSON().ContentType())
	})

	t.Run("RespectsContextCancellation", func(t *testing.T) {
		var w bytes.Buffer

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := JSON().Render(ctx, &w, "tests", NoOptions)

		assert.Error(t, err)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("HandlesWriterError", func(t *testing.T) {
		config := JSONConfig{Padding: "callback"}

		err := NewJSON(config).Render(context.Background(), &errorWriterTest{}, nil, NoOptions)

		assert.Error(t, err)
	})

	t.Run("HandlesWriterErrorOnOpeningParen", func(t *testing.T) {
		// First write (padding name) succeeds, second write ("(") fails.
		config := JSONConfig{Padding: "callback"}
		w := &countedErrorWriter{remaining: 1}

		err := NewJSON(config).Render(context.Background(), w, nil, NoOptions)

		assert.Error(t, err)
	})

	t.Run("RespectsTimeout", func(t *testing.T) {
		var w bytes.Buffer

		err := JSON().Render(context.Background(), &w, map[string]string{"key": "value"}, Options{Timeout: 1 * time.Hour})

		assert.NoError(t, err)
	})
}
