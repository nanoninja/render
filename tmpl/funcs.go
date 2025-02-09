// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package tmpl

import (
	"strings"
	"time"

	"html/template"
)

// DefaultFuncs returns a set of commonly used template functions.
// These functions are automatically available in all templates
// unless overridden by custom functions.
func DefaultFuncs() template.FuncMap {
	return template.FuncMap{
		// String manipulations
		"lower":    strings.ToLower,
		"upper":    strings.ToUpper,
		"trim":     strings.TrimSpace,
		"contains": strings.Contains,

		// Type conversion
		"toHTML": func(s string) template.HTML {
			return template.HTML(s)
		},

		// nl2br converts newlines to <br> tags in text.
		// It handles both \n and \r\n line endings and escapes HTML special characters
		// to prevent XSS attacks while preserving line breaks in HTML output.
		//
		// Example:
		//   Input:  "Line 1\nLine 2"
		//   Output: "Line 1<br>Line 2"
		"nl2br": func(text string) template.HTML {
			// Normalize all line endings to \n
			text = strings.ReplaceAll(text, "\r\n", "\n")

			// Escape HTML special characters for security
			text = template.HTMLEscapeString(text)

			// Replace newlines with <br> tags
			text = strings.ReplaceAll(text, "\n", "<br>")

			// Return as template.HTML to mark content as safely escaped
			return template.HTML(text)
		},

		// Date formating
		"now": time.Now,
		"date": func(t time.Time, layout string) string {
			return t.Format(layout)
		},

		// Basic arithmetic operations using our generic calc function
		"add": func(a, b float64) float64 { return a + b },
		"sub": func(a, b float64) float64 { return a - b },
		"mul": func(a, b float64) float64 { return a * b },
		"div": func(a, b float64) float64 {
			if b == 0 {
				panic("division by zero")
			}
			return a / b
		},

		// Sum calculates the total of a slice of numbers
		"sum": func(numbers []float64) float64 {
			var total float64
			for _, n := range numbers {
				total += n
			}
			return total
		},

		// Average calculates the arithmetic mean of a slice of numbers
		"avg": func(numbers []float64) float64 {
			if len(numbers) == 0 {
				return 0
			}
			var sum float64
			for _, n := range numbers {
				sum += n
			}
			return sum / float64(len(numbers))
		},
	}
}
