// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tmpl

import (
	"html/template"
	"testing"
	"time"

	"github.com/nanoninja/assert"
)

func TestDefaultFuncs(t *testing.T) {
	funcs := DefaultFuncs()

	t.Run("StringManipulation", func(t *testing.T) {
		lower := funcs["lower"].(func(string) string)
		assert.Equal(t, "hello", lower("HELLO"))

		upper := funcs["upper"].(func(string) string)
		assert.Equal(t, "HELLO", upper("hello"))

		trim := funcs["trim"].(func(string) string)
		assert.Equal(t, "hello", trim("  hello  "))

		replace := funcs["replace"].(func(string, string, string) string)
		assert.Equal(t, "baz bar", replace("foo bar", "foo", "baz"))

		contains := funcs["contains"].(func(string, string) bool)
		assert.True(t, contains("seafood", "foo"))
		assert.False(t, contains("seafood", "bar"))

		hasPrefix := funcs["hasPrefix"].(func(string, string) bool)
		assert.True(t, hasPrefix("hello", "he"))
		assert.False(t, hasPrefix("hello", "lo"))

		hasSuffix := funcs["hasSuffix"].(func(string, string) bool)
		assert.True(t, hasSuffix("hello", "lo"))
		assert.False(t, hasSuffix("hello", "he"))

		split := funcs["split"].(func(string, string) []string)
		assert.Equal(t, []string{"a", "b", "c"}, split("a,b,c", ","))

		join := funcs["join"].(func([]string, string) string)
		assert.Equal(t, "a,b,c", join([]string{"a", "b", "c"}, ","))

		truncate := funcs["truncate"].(func(string, int) string)
		assert.Equal(t, "Hello...", truncate("Hello World", 5))
		assert.Equal(t, "Hi", truncate("Hi", 5))
	})

	t.Run("LogicHelpers", func(t *testing.T) {
		def := funcs["default"].(func(any, any) any)
		assert.Equal(t, "Untitled", def("Untitled", ""))
		assert.Equal(t, "Untitled", def("Untitled", nil))
		assert.Equal(t, "My Title", def("Untitled", "My Title"))

		ternary := funcs["ternary"].(func(bool, any, any) any)
		assert.Equal(t, "yes", ternary(true, "yes", "no"))
		assert.Equal(t, "no", ternary(false, "yes", "no"))
	})

	t.Run("SliceHelpers", func(t *testing.T) {
		first := funcs["first"].(func([]any) any)
		assert.Equal(t, "a", first([]any{"a", "b", "c"}))
		assert.Nil(t, first([]any{}))

		last := funcs["last"].(func([]any) any)
		assert.Equal(t, "c", last([]any{"a", "b", "c"}))
		assert.Nil(t, last([]any{}))
	})

	t.Run("TypeConversion", func(t *testing.T) {
		toInt := funcs["int"].(func(float64) int)
		assert.Equal(t, 3, toInt(3.7))

		toFloat := funcs["float"].(func(int) float64)
		assert.Equal(t, 3.0, toFloat(3))
	})

	t.Run("DateAndTime", func(t *testing.T) {
		now := funcs["now"].(func() time.Time)
		assert.True(t, !now().IsZero())

		dateFn := funcs["date"].(func(time.Time, string) string)
		tests := []struct {
			name     string
			time     time.Time
			layout   string
			expected string
		}{
			{
				name:     "standard format",
				time:     time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
				layout:   "2006-01-02",
				expected: "2024-01-15",
			},
			{
				name:     "custom format",
				time:     time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
				layout:   "02/01/2006 15:04",
				expected: "15/01/2024 14:30",
			},
			{
				name:     "short format",
				time:     time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
				layout:   "Jan 02",
				expected: "Jan 15",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.expected, dateFn(tt.time, tt.layout))
			})
		}
	})

	t.Run("Arithmetic", func(t *testing.T) {
		add := funcs["add"].(func(float64, float64) float64)
		assert.Equal(t, 8.0, add(5.0, 3.0))

		sub := funcs["sub"].(func(float64, float64) float64)
		assert.Equal(t, 2.0, sub(5.0, 3.0))

		mul := funcs["mul"].(func(float64, float64) float64)
		assert.Equal(t, 15.0, mul(5.0, 3.0))

		div := funcs["div"].(func(float64, float64) (float64, error))
		result, err := div(6.0, 2.0)
		assert.NoError(t, err)
		assert.Equal(t, 3.0, result)
	})

	t.Run("DivisionByZero", func(t *testing.T) {
		div := funcs["div"].(func(float64, float64) (float64, error))

		_, err := div(1.0, 0.0)
		assert.Error(t, err)
		assert.Equal(t, "division by zero", err.Error())
	})

	t.Run("AggregateFunctions", func(t *testing.T) {
		numbers := []float64{1.0, 2.0, 3.0, 4.0, 5.0}

		sum := funcs["sum"].(func([]float64) float64)
		assert.Equal(t, 15.0, sum(numbers))

		avg := funcs["avg"].(func([]float64) float64)
		assert.Equal(t, 3.0, avg(numbers))
		assert.Equal(t, 0.0, avg([]float64{}))
	})

	t.Run("HTMLUtilities", func(t *testing.T) {
		nl2br := funcs["nl2br"].(func(string) template.HTML)

		tests := []struct {
			name     string
			input    string
			expected template.HTML
		}{
			{
				name:     "simple newline",
				input:    "line1\nline2",
				expected: "line1<br>line2",
			},
			{
				name:     "windows newline",
				input:    "line1\r\nline2",
				expected: "line1<br>line2",
			},
			{
				name:     "html escape",
				input:    "<script>\nAlert()</script>",
				expected: "&lt;script&gt;<br>Alert()&lt;/script&gt;",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.expected, nl2br(tt.input))
			})
		}
	})
}
