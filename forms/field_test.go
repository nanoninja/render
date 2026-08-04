// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forms

import (
	"html/template"
	"strings"
	"testing"

	"github.com/nanoninja/assert"

	"github.com/nanoninja/render"
)

func TestFieldType(t *testing.T) {
	t.Run("BlockName", func(t *testing.T) {
		cases := []struct {
			typ  FieldType
			want string
		}{
			{InputType, "field_input"},
			{TextareaType, "field_textarea"},
			{SelectType, "field_select"},
			{CheckboxType, "field_checkbox"},
			{RadioType, "field_radio"},
			{HiddenType, "field_hidden"},
		}
		for _, c := range cases {
			assert.Equal(t, c.want, c.typ.blockName())
		}
	})
}

func TestBuilders(t *testing.T) {
	t.Run("Text", func(t *testing.T) {
		f := Text("username")
		assert.Equal(t, "username", f.Name)
		assert.Equal(t, InputType, f.Type)
		assert.Equal(t, "text", f.HTMLType)
	})

	t.Run("Email", func(t *testing.T) {
		f := Email("email")
		assert.Equal(t, InputType, f.Type)
		assert.Equal(t, "email", f.HTMLType)
	})

	t.Run("Password", func(t *testing.T) {
		f := Password("password")
		assert.Equal(t, InputType, f.Type)
		assert.Equal(t, "password", f.HTMLType)
	})

	t.Run("Date", func(t *testing.T) {
		f := Date("birthday")
		assert.Equal(t, InputType, f.Type)
		assert.Equal(t, "date", f.HTMLType)
	})

	t.Run("Number", func(t *testing.T) {
		f := Number("age")
		assert.Equal(t, InputType, f.Type)
		assert.Equal(t, "number", f.HTMLType)
	})

	t.Run("Tel", func(t *testing.T) {
		f := Tel("phone")
		assert.Equal(t, InputType, f.Type)
		assert.Equal(t, "tel", f.HTMLType)
	})

	t.Run("URL", func(t *testing.T) {
		f := URL("website")
		assert.Equal(t, InputType, f.Type)
		assert.Equal(t, "url", f.HTMLType)
	})

	t.Run("Color", func(t *testing.T) {
		f := Color("theme_color")
		assert.Equal(t, InputType, f.Type)
		assert.Equal(t, "color", f.HTMLType)
	})

	t.Run("Range", func(t *testing.T) {
		f := Range("volume")
		assert.Equal(t, InputType, f.Type)
		assert.Equal(t, "range", f.HTMLType)
	})

	t.Run("Time", func(t *testing.T) {
		f := Time("appointment")
		assert.Equal(t, InputType, f.Type)
		assert.Equal(t, "time", f.HTMLType)
	})

	t.Run("Search", func(t *testing.T) {
		f := Search("query")
		assert.Equal(t, InputType, f.Type)
		assert.Equal(t, "search", f.HTMLType)
	})

	t.Run("File", func(t *testing.T) {
		f := File("avatar")
		assert.Equal(t, InputType, f.Type)
		assert.Equal(t, "file", f.HTMLType)
	})

	t.Run("Textarea", func(t *testing.T) {
		f := Textarea("bio")
		assert.Equal(t, TextareaType, f.Type)
	})

	t.Run("Select", func(t *testing.T) {
		opts := []Option{{Value: "fr", Label: "France"}, {Value: "be", Label: "Belgique"}}
		f := Select("country", opts...)
		assert.Equal(t, SelectType, f.Type)
		assert.Equal(t, 2, len(f.Options))
		assert.Equal(t, "fr", f.Options[0].Value)
	})

	t.Run("Radio", func(t *testing.T) {
		opts := []Option{{Value: "yes", Label: "Yes"}, {Value: "no", Label: "No"}}
		f := Radio("newsletter_opt_in", opts...)
		assert.Equal(t, RadioType, f.Type)
		assert.Equal(t, 2, len(f.Options))
	})

	t.Run("Checkbox", func(t *testing.T) {
		f := Checkbox("newsletter")
		assert.Equal(t, CheckboxType, f.Type)
	})

	t.Run("CheckboxGroup", func(t *testing.T) {
		opts := []Option{{Value: "go", Label: "Go"}, {Value: "rust", Label: "Rust"}}
		f := CheckboxGroup("interests", opts...)
		assert.Equal(t, CheckboxGroupType, f.Type)
		assert.True(t, f.Multiple)
		assert.Equal(t, 2, len(f.Options))
	})

	t.Run("Hidden", func(t *testing.T) {
		f := Hidden("csrf")
		assert.Equal(t, HiddenType, f.Type)
	})
}

func TestField_With(t *testing.T) {
	t.Run("WithLabel", func(t *testing.T) {
		f := Text("username").WithLabel("Nom d'utilisateur")
		assert.Equal(t, "Nom d'utilisateur", f.Label)
	})

	t.Run("WithValue", func(t *testing.T) {
		f := Text("username").WithValue("jdoe")
		assert.Equal(t, "jdoe", f.Value)
	})

	t.Run("WithPlaceholder", func(t *testing.T) {
		f := Text("username").WithPlaceholder("ex: jdoe")
		assert.Equal(t, "ex: jdoe", f.Placeholder)
	})

	t.Run("WithRequired", func(t *testing.T) {
		f := Text("username").WithRequired()
		assert.True(t, f.Required)
	})

	t.Run("WithMultiple", func(t *testing.T) {
		f := Select("langs").WithMultiple()
		assert.True(t, f.Multiple)
	})

	t.Run("WithGroup", func(t *testing.T) {
		f := Text("first_name").WithGroup("left")
		assert.Equal(t, "left", f.Group)
	})

	t.Run("WithAttrOnNilMap", func(t *testing.T) {
		f := Text("username")
		f.WithAttr("data-test", "1")
		assert.Equal(t, "1", f.Attrs["data-test"])
	})

	t.Run("WithAttrOverwritesExistingKey", func(t *testing.T) {
		f := Text("username").WithAttr("data-test", "1").WithAttr("data-test", "2")
		assert.Equal(t, "2", f.Attrs["data-test"])
		assert.Equal(t, 1, len(f.Attrs))
	})

	t.Run("WithClassAccumulates", func(t *testing.T) {
		f := Text("username").WithClass("form-control").WithClass("is-invalid")
		assert.Equal(t, "form-control is-invalid", f.Attrs["class"])
	})

	t.Run("WithClassAcceptsMultipleArgumentsAtOnce", func(t *testing.T) {
		f := Text("username").WithClass("a", "b")
		assert.Equal(t, "a b", f.Attrs["class"])
	})

	t.Run("WithStyleReplacesPreviousValue", func(t *testing.T) {
		f := Text("username").WithStyle("color: red").WithStyle("color: blue")
		assert.Equal(t, "color: blue", f.Attrs["style"])
	})

	t.Run("MethodsReturnSamePointerForChaining", func(t *testing.T) {
		f := Text("username")
		got := f.WithLabel("x").WithRequired()
		assert.Equal(t, f, got)
	})
}

func TestField_IsSelected(t *testing.T) {
	t.Run("SingleValueMatches", func(t *testing.T) {
		f := Select("country").WithValue("be")
		assert.True(t, f.IsSelected("be"))
		assert.False(t, f.IsSelected("fr"))
	})

	t.Run("SingleValueWrongType", func(t *testing.T) {
		f := Select("country").WithValue(42)
		assert.False(t, f.IsSelected("42"))
	})

	t.Run("MultipleMatchesAnyValueInSlice", func(t *testing.T) {
		f := Select("langs").WithMultiple()
		f.Value = []string{"fr", "de"}
		assert.True(t, f.IsSelected("fr"))
		assert.True(t, f.IsSelected("de"))
		assert.False(t, f.IsSelected("en"))
	})

	t.Run("MultipleWithWrongValueType", func(t *testing.T) {
		f := Select("langs").WithMultiple()
		f.Value = "fr"
		assert.False(t, f.IsSelected("fr"))
	})
}

func TestField_IsGroup(t *testing.T) {
	cases := []struct {
		name string
		f    *Field
		want bool
	}{
		{"Radio", Radio("plan", Option{Value: "free", Label: "Free"}), true},
		{"CheckboxGroup", CheckboxGroup("interests", Option{Value: "go", Label: "Go"}), true},
		{"Text", Text("username"), false},
		{"Checkbox", Checkbox("newsletter"), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, c.f.IsGroup())
		})
	}
}

func TestField_RenderAttrs(t *testing.T) {
	t.Run("EmptyWhenNoAttrs", func(t *testing.T) {
		f := Text("username")
		assert.Equal(t, template.HTMLAttr(""), f.RenderAttrs())
	})

	t.Run("RendersASingleAttribute", func(t *testing.T) {
		f := Text("username").WithAttr("data-test", "hello")
		assert.Equal(t, template.HTMLAttr(` data-test="hello"`), f.RenderAttrs())
	})

	t.Run("SortsKeysForDeterministicOutput", func(t *testing.T) {
		f := Text("username").WithAttr("zebra", "1").WithAttr("alpha", "2")
		assert.Equal(t, template.HTMLAttr(` alpha="2" zebra="1"`), f.RenderAttrs())
	})

	t.Run("SkipsKeysNamedInExcept", func(t *testing.T) {
		f := Text("username").WithAttr("class", "form-control").WithAttr("data-test", "1")
		assert.Equal(t, template.HTMLAttr(` data-test="1"`), f.RenderAttrs("class"))
	})

	t.Run("EscapesKeyAndValue", func(t *testing.T) {
		f := Text("username").WithAttr(`x"y`, `a"b`)
		assert.Equal(t, template.HTMLAttr(` x&#34;y="a&#34;b"`), f.RenderAttrs())
	})

	t.Run("NeverProducesZgotmplZWhenRenderedInATemplate", func(t *testing.T) {
		r := NewRenderer("default", ThemeDefault())
		f := Text("username").WithAttr("formaction", "/elsewhere")

		html, err := r.Field(f, render.NoOptions)

		assert.NoError(t, err)
		assert.False(t, strings.Contains(string(html), "ZgotmplZ"))
		assert.StringContains(t, string(html), `formaction="/elsewhere"`)
	})
}

func TestWithOverrides(t *testing.T) {
	t.Run("NoPairsReturnsSameField", func(t *testing.T) {
		f := Text("username")
		got := withOverrides(f)
		assert.Equal(t, f, got)
	})

	t.Run("DoesNotMutateOriginal", func(t *testing.T) {
		f := Text("username").WithClass("form-control")

		clone := withOverrides(f, "class", "email-input")

		assert.Equal(t, "form-control", f.Attrs["class"])
		assert.Equal(t, "email-input", clone.Attrs["class"])
	})

	t.Run("PreservesUnrelatedAttrs", func(t *testing.T) {
		f := Text("username").WithAttr("data-test", "1")

		clone := withOverrides(f, "class", "email-input")

		assert.Equal(t, "1", clone.Attrs["data-test"])
		assert.Equal(t, "email-input", clone.Attrs["class"])
	})

	t.Run("OverwritesRatherThanAccumulatesClass", func(t *testing.T) {
		f := Text("username").WithClass("form-control")

		clone := withOverrides(f, "class", "email-input")

		assert.Equal(t, "email-input", clone.Attrs["class"])
	})
}
