// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tmpl

import (
	"errors"
	htmltpl "html/template"
	"strings"
	"text/template"
	"time"
)

// DefaultFuncs returns a set of commonly used template functions.
// These functions are available via WithDefaultFuncs() and WithDefaultFuncsHTML().
// They cover string manipulation, type conversion, logic helpers,
// date formatting, arithmetic, and HTML utilities.
func DefaultFuncs() template.FuncMap {
	return template.FuncMap{

		// -------------------------
		// String manipulation
		// -------------------------

		// lower converts a string to lowercase.
		// Example: {{ "Hello" | lower }} → "hello"
		"lower": strings.ToLower,

		// upper converts a string to uppercase.
		// Example: {{ "hello" | upper }} → "HELLO"
		"upper": strings.ToUpper,

		// trim removes leading and trailing whitespace.
		// Example: {{ "  hello  " | trim }} → "hello"
		"trim": strings.TrimSpace,

		// replace replaces all occurrences of old with new in s.
		// Example: {{ replace "foo bar" "foo" "baz" }} → "baz bar"
		"replace": strings.ReplaceAll,

		// contains reports whether substr is within s.
		// Example: {{ if contains "seafood" "foo" }}yes{{ end }}
		"contains": strings.Contains,

		// hasPrefix reports whether s begins with prefix.
		// Example: {{ if hasPrefix "hello" "he" }}yes{{ end }}
		"hasPrefix": strings.HasPrefix,

		// hasSuffix reports whether s ends with suffix.
		// Example: {{ if hasSuffix "hello" "lo" }}yes{{ end }}
		"hasSuffix": strings.HasSuffix,

		// split splits s into substrings separated by sep.
		// Example: {{ split "a,b,c" "," }} → ["a", "b", "c"]
		"split": strings.Split,

		// join concatenates elements of a slice with sep.
		// Example: {{ join .Items ", " }} → "a, b, c"
		"join": strings.Join,

		// truncate shortens s to n characters, appending "..." if truncated.
		// Example: {{ truncate "Hello World" 5 }} → "Hello..."
		"truncate": func(s string, n int) string {
			if len(s) <= n {
				return s
			}
			return s[:n] + "..."
		},

		// -------------------------
		// Logic helpers
		// -------------------------

		// default returns val if non-zero, otherwise returns def.
		// Example: {{ default "Untitled" .Title }}
		"default": func(def, val any) any {
			if val == nil || val == "" || val == 0 || val == false {
				return def
			}
			return val
		},

		// ternary returns a if condition is true, otherwise b.
		// Example: {{ ternary .IsAdmin "Admin" "User" }}
		"ternary": func(condition bool, a, b any) any {
			if condition {
				return a
			}
			return b
		},

		// -------------------------
		// Slice helpers
		// -------------------------

		// first returns the first element of a slice, or nil if empty.
		// Example: {{ first .Items }}
		"first": func(slice []any) any {
			if len(slice) == 0 {
				return nil
			}
			return slice[0]
		},

		// last returns the last element of a slice, or nil if empty.
		// Example: {{ last .Items }}
		"last": func(slice []any) any {
			if len(slice) == 0 {
				return nil
			}
			return slice[len(slice)-1]
		},

		// dict builds a map[string]any from alternating key/value
		// arguments. Keys must be strings.
		// Example: {{ dict "class" "email-input" "autofocus" true }}
		"dict": func(pairs ...any) (map[string]any, error) {
			if len(pairs)%2 != 0 {
				return nil, errors.New("dict: odd number of arguments")
			}
			d := make(map[string]any, len(pairs)/2)
			for i := 0; i < len(pairs); i += 2 {
				key, ok := pairs[i].(string)
				if !ok {
					return nil, errors.New("dict: keys must be string")
				}
				d[key] = pairs[i+1]
			}
			return d, nil
		},

		// -------------------------
		// Type conversion
		// -------------------------

		// int converts a float64 to int.
		// Example: {{ int 3.7 }} → 3
		"int": func(f float64) int { return int(f) },

		// float converts an int to float64.
		// Example: {{ float 3 }} → 3.0
		"float": func(i int) float64 { return float64(i) },

		// -------------------------
		// Date and time
		// -------------------------

		// now returns the current local time.
		// Example: {{ now | date "2006-01-02" }}
		"now": time.Now,

		// date formats t using the given layout string.
		// Example: {{ date .CreatedAt "02/01/2006" }}
		"date": func(t time.Time, layout string) string {
			return t.Format(layout)
		},

		// -------------------------
		// Arithmetic
		// -------------------------

		// add returns a + b.
		"add": func(a, b float64) float64 { return a + b },

		// sub returns a - b.
		"sub": func(a, b float64) float64 { return a - b },

		// mul returns a * b.
		"mul": func(a, b float64) float64 { return a * b },

		// div returns a / b. Returns an error on division by zero.
		"div": func(a, b float64) (float64, error) {
			if b == 0 {
				return 0, errors.New("division by zero")
			}
			return a / b, nil
		},

		// sum returns the total of a slice of float64.
		// Example: {{ sum .Prices }}
		"sum": func(numbers []float64) float64 {
			var total float64
			for _, n := range numbers {
				total += n
			}
			return total
		},

		// avg returns the arithmetic mean of a slice of float64.
		// Returns 0 for an empty slice.
		// Example: {{ avg .Scores }}
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

		// -------------------------
		// HTML utilities
		// -------------------------

		// nl2br converts newline characters to HTML <br> tags.
		// HTML special characters are escaped to prevent XSS.
		// Example: {{ nl2br .Comment }}
		"nl2br": func(text string) htmltpl.HTML {
			text = strings.ReplaceAll(text, "\r\n", "\n")
			text = htmltpl.HTMLEscapeString(text)
			text = strings.ReplaceAll(text, "\n", "<br>")
			return htmltpl.HTML(text)
		},
	}
}
