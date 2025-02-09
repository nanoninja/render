// Copyright 2025 The Nanoninja Authors. All rights reserved.
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

		err := XML().Render(&w, data)

		assert.Nil(t, err)
		assert.Equals(t, w.String(), expected)
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

		err := NewXML(config).Render(&w, data)

		assert.Nil(t, err)
		assert.Equals(t, w.String(), `<root><message>Hello</message></root>`)
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

		err := NewXML(config).Render(&w, data, Format(Pretty()))

		assert.Nil(t, err)
		assert.Equals(t, w.String(), xmlPrettyTest)
	})

	t.Run("RenderInvalidData", func(t *testing.T) {
		var w bytes.Buffer

		data := make(chan int)
		err := XML().Render(&w, data)

		assert.NotNil(t, err)
	})

	t.Run("RenderWithCustomPrefixAndIndent", func(t *testing.T) {
		var w bytes.Buffer

		data := struct {
			XMLName xml.Name `xml:"root"`
			Message string   `xml:"message"`
		}{
			Message: "Prefix test",
		}

		err := XML().Render(&w, data, Format(
			Pretty(),
			Prefix(">>"),
			Indent("  "),
		))

		assert.Nil(t, err)
		assert.Equals(t, w.String(), xmlPrefixTest)
	})

	t.Run("HandleHeaderWriteError", func(t *testing.T) {
		errWriter := &errorWriterTest{}

		err := XML().Render(errWriter, struct{}{})

		assert.NotNil(t, err)
	})

	t.Run("RespectsContextCancellation", func(t *testing.T) {
		var w bytes.Buffer

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := XML().RenderContext(ctx, &w, struct{}{})

		assert.ErrorIs(t, err, context.Canceled)
	})
}
