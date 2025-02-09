// Copyright 2025 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package render

import (
	"testing"

	"github.com/nanoninja/assert"
)

func TestFormatOptions(t *testing.T) {
	t.Run("IndentReturnsConfiguraredIndentation", func(t *testing.T) {
		format := FormatOptions{indent: "    "}
		assert.Equals(t, format.Indent(), "    ")
	})

	t.Run("LineEndingReturnsDefaultWhenEmpty", func(t *testing.T) {
		format := FormatOptions{}
		assert.Equals(t, format.LineEnding(), "\n")
	})

	t.Run("LineEndingReturnsConfiguredValue", func(t *testing.T) {
		format := FormatOptions{lineEnding: "\r\n"}
		assert.Equals(t, format.LineEnding(), "\r\n")
	})

	t.Run("PrefixReturnsConfiguredPrefix", func(t *testing.T) {
		format := FormatOptions{prefix: "-->"}
		assert.Equals(t, format.Prefix(), "-->")
	})

	t.Run("PrettyReturnsConfifuredState", func(t *testing.T) {
		format := FormatOptions{pretty: true}
		assert.True(t, format.Pretty())
	})

	t.Run("CloneCreatesSeparateInstance", func(t *testing.T) {
		original := FormatOptions{
			prefix:     "-->",
			lineEnding: "\r\n",
			indent:     "    ",
			pretty:     true,
			args:       []any{"test", 123},
		}

		clone := original.Clone()

		assert.Equals(t, clone.prefix, original.prefix)
		assert.Equals(t, clone.lineEnding, original.lineEnding)
		assert.Equals(t, clone.indent, original.indent)
		assert.Equals(t, clone.pretty, original.pretty)
		assert.Equals(t, len(clone.args), len(original.args))
		assert.Equals(t, clone.args[0], original.args[0])
		assert.Equals(t, clone.args[1], original.args[1])

		clone.args[0] = "modified"
		assert.NotEquals(t, clone.args[0], original.args[0])
	})

	t.Run("ArgsSetArgumentInFormatOptions", func(t *testing.T) {
		format := &FormatOptions{}

		Args("test", 123)(format)

		assert.Equals(t, len(format.args), 2)
		assert.Equals(t, format.args[0], "test")
		assert.Equals(t, format.args[1], 123)
	})

	t.Run("FormatAppliesMultipleFormatters", func(t *testing.T) {
		opts := NewOptions()

		Format(
			Indent("  "),
			Pretty(),
			LineEnding("\n"),
		)(opts)

		assert.Equals(t, opts.format.Indent(), "  ")
		assert.True(t, opts.format.Pretty())
		assert.Equals(t, opts.format.LineEnding(), "\n")
	})

	t.Run("CommentAddsMarkerWithSpace", func(t *testing.T) {
		format := &FormatOptions{}

		Comment("//")(format)

		assert.Equals(t, format.Prefix(), "// ")
	})

	t.Run("PrettyPreservesExistingIndent", func(t *testing.T) {
		format := &FormatOptions{
			indent: "    ",
		}

		Pretty()(format)

		assert.True(t, format.Pretty())
		assert.Equals(t, format.Indent(), "    ")
	})

	t.Run("IndentSetsIndentationsString", func(t *testing.T) {
		format := &FormatOptions{}

		Indent("    ")(format)

		assert.Equals(t, format.Indent(), "    ")
	})

	t.Run("PrefixSetsPrefixString", func(t *testing.T) {
		format := &FormatOptions{}

		Prefix("-->")(format)

		assert.Equals(t, format.Prefix(), "-->")
	})

	t.Run("TextCombinesFormatAndArgs", func(t *testing.T) {
		opts := NewOptions()

		Textf("test", 123)(opts)

		assert.Len(t, opts.format.args, 2)
		assert.Equals(t, opts.format.args[0], "test")
		assert.Equals(t, opts.format.args[1], 123)
	})

	t.Run("SetsCRLFLineEndingInOptions", func(t *testing.T) {
		opts := NewOptions().
			Use(UseCRLF())

		assert.Equals(t, opts.format.LineEnding(), "\r\n")
	})

	t.Run("EquivalentToFormatWithLineEnding", func(t *testing.T) {
		opts1 := NewOptions()
		opts2 := NewOptions()

		opts1.Use(UseCRLF())
		opts2.Use(Format(LineEnding("\r\n")))

		// Verify they produce the same result
		assert.Equals(t, opts1.format.LineEnding(), opts2.format.LineEnding())
	})
}
