// Copyright 2025 The Nanoninja Authors. All rights reserved.
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

	t.Run("ToHTML", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected template.HTML
		}{
			{
				name:     "simple text",
				input:    "Hello World",
				expected: "Hello World",
			},
			{
				name:     "html content",
				input:    "<p>Hello</p>",
				expected: "<p>Hello</p>",
			},
			{
				name:     "empty string",
				input:    "",
				expected: "",
			},
		}

		toHTML := funcs["toHTML"].(func(string) template.HTML)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := toHTML(tt.input)

				assert.Equals(t, result, tt.expected)
				if result != tt.expected {
					t.Errorf("toHTML(%q) = %q; want %q", tt.input, result, tt.expected)
				}
			})
		}
	})

	t.Run("nl2br", func(t *testing.T) {
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

		nl2br := funcs["nl2br"].(func(string) template.HTML)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := nl2br(tt.input)

				if result != tt.expected {
					t.Errorf("nl2br(%q) = %q; want %q", tt.input, result, tt.expected)
				}
			})
		}
	})

	t.Run("Date", func(t *testing.T) {
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

		dateFn := funcs["date"].(func(time.Time, string) string)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := dateFn(tt.time, tt.layout)

				if result != tt.expected {
					t.Errorf("date(%v, %q) = %q; want %q", tt.time, tt.layout, result, tt.expected)
				}
			})
		}
	})

	t.Run("ArithmeticOperation", func(t *testing.T) {
		tests := []struct {
			name     string
			fn       string
			a, b     float64
			expected float64
		}{
			{"addition", "add", 5.0, 3.0, 8.0},
			{"subtraction", "sub", 5.0, 3.0, 2.0},
			{"multiplication", "mul", 5.0, 3.0, 15.0},
			{"division", "div", 6.0, 2.0, 3.0},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				fn, exists := funcs[tt.fn]
				if !exists {
					t.Fatalf("function %s not found", tt.fn)
				}

				result := fn.(func(float64, float64) float64)(tt.a, tt.b)
				if result != tt.expected {
					t.Errorf("%s(%v, %v) = %v; want %v", tt.fn, tt.a, tt.b, result, tt.expected)
				}
			})
		}
	})

	t.Run("DivisionByZero", func(t *testing.T) {
		div := funcs["div"].(func(float64, float64) float64)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("division by zero did not panic")
			}
		}()

		div(1.0, 0.0)
	})

	t.Run("AggregateFunctions", func(t *testing.T) {
		numbers := []float64{1.0, 2.0, 3.0, 4.0, 5.0}

		// Test sum
		sum := funcs["sum"].(func([]float64) float64)(numbers)
		if sum != 15.0 {
			t.Errorf("sum(%v) = %v; want 15.0", numbers, sum)
		}

		// Test average
		avg := funcs["avg"].(func([]float64) float64)(numbers)
		if avg != 3.0 {
			t.Errorf("avg(%v) = %v; want 3.0", numbers, avg)
		}

		// Test average empty values
		avg = funcs["avg"].(func([]float64) float64)([]float64{})
		if avg != 0 {
			t.Errorf("avg(%v) = %v; want 0.0", []float64{}, avg)
		}
	})
}
