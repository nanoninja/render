// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"context"
	"encoding/xml"
	"testing"

	"github.com/nanoninja/assert"
)

var (
	_ Renderer = (*xmlRenderer)(nil)
	_ Renderer = XML()
	_ Renderer = NewXML(XMLConfig{})
)

const xmlPrettyTest = xml.Header + `<root>
    <message>Pretty test</message>
</root>`

const xmlPrefixTest = xml.Header + `>><root>
>>  <message>Prefix test</message>
>></root>`

func TestXMLRenderer(t *testing.T) {
	t.Run("RenderSimpleXML", func(t *testing.T) {
		var w bytes.Buffer

		data := struct {
			XMLName xml.Name `xml:"root"`
			Message string   `xml:"message"`
		}{
			Message: "test",
		}

		expected := xml.Header + "<root><message>test</message></root>"

		err := XML().Render(context.Background(), &w, data, NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, expected, w.String())
	})

	t.Run("RenderWithoutHeader", func(t *testing.T) {
		var w bytes.Buffer

		data := struct {
			XMLName xml.Name `xml:"root"`
			Message string   `xml:"message"`
		}{
			Message: "Hello",
		}

		config := XMLConfig{Header: false}

		err := NewXML(config).Render(context.Background(), &w, data, NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, `<root><message>Hello</message></root>`, w.String())
	})

	t.Run("RenderPrettyPrintWithConfigIndent", func(t *testing.T) {
		var w bytes.Buffer

		data := struct {
			XMLName xml.Name `xml:"root"`
			Message string   `xml:"message"`
		}{
			Message: "Pretty test",
		}

		config := XMLConfig{
			Indent: "    ",
			Header: true,
		}

		err := NewXML(config).Render(context.Background(), &w, data, Options{Pretty: true})

		assert.NoError(t, err)
		assert.Equal(t, xmlPrettyTest, w.String())
	})

	t.Run("RenderInvalidData", func(t *testing.T) {
		var w bytes.Buffer

		data := make(chan int)
		err := XML().Render(context.Background(), &w, data, NoOptions)

		assert.Error(t, err)
	})

	t.Run("RenderWithCustomPrefixAndIndent", func(t *testing.T) {
		var w bytes.Buffer

		data := struct {
			XMLName xml.Name `xml:"root"`
			Message string   `xml:"message"`
		}{
			Message: "Prefix test",
		}

		err := XML().Render(context.Background(), &w, data, Options{Pretty: true, Prefix: ">>", Indent: "  "})

		assert.NoError(t, err)
		assert.Equal(t, xmlPrefixTest, w.String())
	})

	t.Run("HandleHeaderWriteError", func(t *testing.T) {
		errWriter := &errorWriterTest{}

		err := XML().Render(context.Background(), errWriter, struct{}{}, NoOptions)

		assert.Error(t, err)
	})

	t.Run("RespectsContextCancellation", func(t *testing.T) {
		var w bytes.Buffer

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := XML().Render(ctx, &w, struct{}{}, NoOptions)

		assert.ErrorIs(t, err, context.Canceled)
	})
}
