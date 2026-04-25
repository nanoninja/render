// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/nanoninja/assert"
	"go.yaml.in/yaml/v3"
)

type yamlErrorType struct{}

func (*yamlErrorType) MarshalYAML() (any, error) {
	return nil, errors.New("yaml marshal error")
}

var _ yaml.Marshaler = (*yamlErrorType)(nil)

var (
	_ Renderer = (*yamlRenderer)(nil)
	_ Renderer = YAML()
	_ Renderer = NewYAML(YAMLConfig{})
)

func TestYAMLRenderer(t *testing.T) {
	t.Run("RendersSimpleData", func(t *testing.T) {
		var w bytes.Buffer

		data := map[string]string{"message": "hello"}

		err := YAML().Render(context.Background(), &w, data, NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, w.String(), "message: hello")
	})

	t.Run("RendersStruct", func(t *testing.T) {
		var w bytes.Buffer

		data := struct {
			Name string `yaml:"name"`
			Age  int    `yaml:"age"`
		}{
			Name: "Gopher",
			Age:  10,
		}

		err := YAML().Render(context.Background(), &w, data, NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, w.String(), "name: Gopher")
		assert.StringContains(t, w.String(), "age: 10")
	})

	t.Run("RendersWithCustomIndent", func(t *testing.T) {
		var w bytes.Buffer

		data := map[string]any{
			"server": map[string]any{
				"host": "localhost",
				"port": 8080,
			},
		}

		err := NewYAML(YAMLConfig{Indent: 4}).Render(context.Background(), &w, data, NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, w.String(), "    host: localhost")
	})

	t.Run("RendersSlice", func(t *testing.T) {
		var w bytes.Buffer

		data := []string{"foo", "bar", "baz"}

		err := YAML().Render(context.Background(), &w, data, NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, w.String(), "- foo")
		assert.StringContains(t, w.String(), "- bar")
	})

	t.Run("SetsDefaultContentType", func(t *testing.T) {
		assert.Equal(t, "application/yaml; charset=utf-8", YAML().ContentType())
	})

	t.Run("OutputEndsWithNewline", func(t *testing.T) {
		var w bytes.Buffer

		err := YAML().Render(context.Background(), &w, map[string]string{"key": "value"}, NoOptions)

		assert.NoError(t, err)
		assert.True(t, strings.HasSuffix(w.String(), "\n"))
	})

	t.Run("RespectsContextCancellation", func(t *testing.T) {
		var w bytes.Buffer

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := YAML().Render(ctx, &w, map[string]string{"key": "value"}, NoOptions)

		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("RespectsTimeout", func(t *testing.T) {
		var w bytes.Buffer

		err := YAML().Render(context.Background(), &w, map[string]string{"key": "value"}, Options{Timeout: 1 * time.Hour})

		assert.NoError(t, err)
	})

	t.Run("ReturnsErrorOnEncodeFailure", func(t *testing.T) {
		data := &yamlErrorType{}
		var w bytes.Buffer

		err := YAML().Render(context.Background(), &w, data, NoOptions)

		assert.Error(t, err)
		assert.StringContains(t, err.Error(), "yaml:")
	})
}
