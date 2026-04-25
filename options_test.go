// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"testing"
	"time"

	"github.com/nanoninja/assert"
)

func TestOptions(t *testing.T) {
	t.Run("ZeroValueIsValid", func(t *testing.T) {
		var opts Options

		assert.Equal(t, "", opts.Name)
		assert.Equal(t, time.Duration(0), opts.Timeout)
		assert.False(t, opts.Pretty)
		assert.Equal(t, "", opts.Indent)
		assert.Equal(t, "", opts.Prefix)
		assert.Equal(t, "", opts.LineEnding)
		assert.Equal(t, "", opts.Padding)
		assert.Equal(t, rune(0), opts.Separator)
		assert.Nil(t, opts.Args)
		assert.Nil(t, opts.Headers)
	})

	t.Run("NoOptionsIsZeroValue", func(t *testing.T) {
		assert.Equal(t, Options{}, NoOptions)
	})

	t.Run("FieldsAreSetDirectly", func(t *testing.T) {
		opts := Options{
			Name:       "template.tmpl",
			Timeout:    5 * time.Second,
			Pretty:     true,
			Indent:     "  ",
			Prefix:     ">>",
			LineEnding: "\r\n",
			Padding:    "callback",
			Separator:  ';',
			Args:       []any{"arg1", "arg2"},
		}

		assert.Equal(t, "template.tmpl", opts.Name)
		assert.Equal(t, 5*time.Second, opts.Timeout)
		assert.True(t, opts.Pretty)
		assert.Equal(t, "  ", opts.Indent)
		assert.Equal(t, ">>", opts.Prefix)
		assert.Equal(t, "\r\n", opts.LineEnding)
		assert.Equal(t, "callback", opts.Padding)
		assert.Equal(t, ';', opts.Separator)
		assert.Len(t, opts.Args, 2)
	})

	t.Run("PassedByValueDoesNotMutateOriginal", func(t *testing.T) {
		original := Options{Name: "original", Timeout: time.Second}

		mutate := func(o Options) {
			o.Name = "mutated"
			o.Timeout = 0
		}
		mutate(original)

		assert.Equal(t, "original", original.Name)
		assert.Equal(t, time.Second, original.Timeout)
	})

	t.Run("HeadersMapIsNilByDefault", func(t *testing.T) {
		var opts Options
		assert.Nil(t, opts.Headers)
	})

	t.Run("HeadersMapCanBeSet", func(t *testing.T) {
		opts := Options{
			Headers: map[string][]string{
				"Content-Type": {"text/plain"},
			},
		}
		assert.Equal(t, "text/plain", opts.Headers["Content-Type"][0])
	})
}
