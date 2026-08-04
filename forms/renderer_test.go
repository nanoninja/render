// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forms

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/nanoninja/assert"
	"github.com/nanoninja/render"
)

var _ render.Renderer = (*Renderer)(nil)

func TestRenderer_ContentType(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())
	assert.Equal(t, "text/html; charset=utf-8", r.ContentType())
}

func TestRenderer_Field(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())

	t.Run("RendersInput", func(t *testing.T) {
		f := Text("username").WithPlaceholder("ex: jdoe").WithRequired()

		html, err := r.Field(f, render.NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, string(html), `type="text"`)
		assert.StringContains(t, string(html), `name="username"`)
		assert.StringContains(t, string(html), `placeholder="ex: jdoe"`)
		assert.StringContains(t, string(html), `required`)
	})

	t.Run("EscapesValues", func(t *testing.T) {
		f := Text("bio").WithValue(`<script>alert(1)</script>`)

		html, err := r.Field(f, render.NoOptions)

		assert.NoError(t, err)
		assert.False(t, strings.Contains(string(html), "<script>"))
	})

	t.Run("UnknownThemeReturnsError", func(t *testing.T) {
		f := Text("username")

		_, err := r.Field(f, render.Options{Name: "bogus"})

		assert.Error(t, err)
	})
}

func TestRenderer_Label(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())

	t.Run("RendersLabelWhenSet", func(t *testing.T) {
		f := Text("username").WithLabel("Nom d'utilisateur")

		html, err := r.Label(f, render.NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, string(html), "Nom d&#39;utilisateur")
		assert.StringContains(t, string(html), `for="username"`)
	})

	t.Run("EmptyWhenNoLabel", func(t *testing.T) {
		f := Text("username")

		html, err := r.Label(f, render.NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "", string(html))
	})
}

func TestRenderer_Errors(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())

	t.Run("RendersEachErrorMessage", func(t *testing.T) {
		f := Text("username")
		f.Errors = []string{"required", "too short"}

		html, err := r.Errors(f, render.NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, string(html), "required")
		assert.StringContains(t, string(html), "too short")
	})

	t.Run("EmptyWhenNoErrors", func(t *testing.T) {
		f := Text("username")

		html, err := r.Errors(f, render.NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "", string(html))
	})
}

func TestRenderer_FormErrors(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())

	t.Run("RendersFormLevelErrors", func(t *testing.T) {
		f := New("login")
		f.Errors = []string{"too many attempts"}

		html, err := r.FormErrors(f, render.NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, string(html), "too many attempts")
	})

	t.Run("EmptyWhenNoFormErrors", func(t *testing.T) {
		f := New("login")

		html, err := r.FormErrors(f, render.NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "", string(html))
	})

	t.Run("IncludedAutomaticallyInForm", func(t *testing.T) {
		f := New("login")
		f.Errors = []string{"too many attempts"}

		html, err := r.Form(f, render.NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, string(html), "too many attempts")
	})
}

func TestRenderer_FormStartAndEnd(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())
	f := New("login").WithAction("/login")

	t.Run("FormStart", func(t *testing.T) {
		html, err := r.FormStart(f, render.NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, string(html), `<form`)
		assert.StringContains(t, string(html), `action="/login"`)
		assert.False(t, strings.Contains(string(html), "</form>"))
	})

	t.Run("FormEnd", func(t *testing.T) {
		html, err := r.FormEnd(f, render.NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "</form>", string(html))
	})

	t.Run("AllowManualAssemblyWithASubmitButton", func(t *testing.T) {
		start, err := r.FormStart(f, render.NoOptions)
		assert.NoError(t, err)

		button := `<button type="submit">Envoyer</button>`

		end, err := r.FormEnd(f, render.NoOptions)
		assert.NoError(t, err)

		got := string(start) + button + string(end)

		assert.StringContains(t, got, `<form`)
		assert.StringContains(t, got, button)
		assert.StringContains(t, got, "</form>")
	})
}

func TestRenderer_Row(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())

	f := Text("username").WithLabel("Nom d'utilisateur")
	f.Errors = []string{"Ce nom d'utilisateur est déjà pris."}

	html, err := r.Row(f, render.NoOptions)

	assert.NoError(t, err)
	assert.StringContains(t, string(html), "Nom d&#39;utilisateur")
	assert.StringContains(t, string(html), `type="text"`)
	assert.StringContains(t, string(html), "déjà pris")
}

func TestRenderer_Form(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())

	t.Run("RendersFormTagAndFields", func(t *testing.T) {
		f := New("login").
			WithAction("/login").
			Add(Text("username"), Email("email"))

		html, err := r.Form(f, render.NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, string(html), `<form`)
		assert.StringContains(t, string(html), `action="/login"`)
		assert.StringContains(t, string(html), `name="username"`)
		assert.StringContains(t, string(html), `name="email"`)
		assert.StringContains(t, string(html), `</form>`)
	})

	t.Run("DeclaresMultipartEnctypeWithFileField", func(t *testing.T) {
		f := New("upload").Add(File("avatar"))

		html, err := r.Form(f, render.NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, string(html), `enctype="multipart/form-data"`)
	})
}

func TestRenderer_Group(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())

	t.Run("RendersOnlyFieldsInGroup", func(t *testing.T) {
		f := New("register").Add(
			Text("first_name").WithGroup("left"),
			Email("email").WithGroup("right"),
			Text("last_name").WithGroup("left"),
		)

		html, err := r.Group(f, "left", render.NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, string(html), `name="first_name"`)
		assert.StringContains(t, string(html), `name="last_name"`)
		assert.False(t, strings.Contains(string(html), `name="email"`))
	})

	t.Run("EmptyWhenGroupHasNoFields", func(t *testing.T) {
		f := New("register").Add(Text("first_name").WithGroup("left"))

		html, err := r.Group(f, "right", render.NoOptions)

		assert.NoError(t, err)
		assert.Equal(t, "", string(html))
	})

	t.Run("DoesNotIncludeFormTagsOrGlobalErrors", func(t *testing.T) {
		f := New("register")
		f.Errors = []string{"should not appear"}
		f.Add(Text("first_name").WithGroup("left"))

		html, err := r.Group(f, "left", render.NoOptions)

		assert.NoError(t, err)
		assert.False(t, strings.Contains(string(html), "<form"))
		assert.False(t, strings.Contains(string(html), "should not appear"))
	})
}

func TestRenderer_ThemeSelection(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())
	assert.NoError(t, r.RegisterTheme("bootstrap5", ThemeBootstrap5()))

	f := Email("email").WithLabel("Email")

	defaultHTML, err := r.Row(f, render.NoOptions)
	assert.NoError(t, err)

	bootstrapHTML, err := r.Row(f, render.Options{Name: "bootstrap5"})
	assert.NoError(t, err)

	assert.StringContains(t, string(defaultHTML), `class="form-row"`)
	assert.StringContains(t, string(bootstrapHTML), `class="form-label"`)
}

func TestRenderer_FieldTypes(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())
	assert.NoError(t, r.RegisterTheme("bootstrap5", ThemeBootstrap5()))

	themes := []struct {
		name string
		opts render.Options
	}{
		{"default", render.NoOptions},
		{"bootstrap5", render.Options{Name: "bootstrap5"}},
	}

	for _, th := range themes {
		t.Run(th.name, func(t *testing.T) {
			t.Run("Textarea", func(t *testing.T) {
				f := Textarea("bio").WithValue("hello")

				html, err := r.Field(f, th.opts)

				assert.NoError(t, err)
				assert.StringContains(t, string(html), `<textarea`)
				assert.StringContains(t, string(html), `name="bio"`)
				assert.StringContains(t, string(html), ">hello</textarea>")
			})

			t.Run("Select", func(t *testing.T) {
				f := Select("country",
					Option{Value: "fr", Label: "France"},
					Option{Value: "be", Label: "Belgique"},
				).WithValue("be")

				html, err := r.Field(f, th.opts)

				assert.NoError(t, err)
				assert.StringContains(t, string(html), `<select`)
				assert.StringContains(t, string(html), `<option value="fr">France</option>`)
				assert.StringContains(t, string(html), `<option value="be" selected>Belgique</option>`)
			})

			t.Run("SelectMultiple", func(t *testing.T) {
				f := Select("langs",
					Option{Value: "fr", Label: "Français"},
					Option{Value: "en", Label: "Anglais"},
				).WithMultiple()
				f.Value = []string{"fr"}

				html, err := r.Field(f, th.opts)

				assert.NoError(t, err)
				assert.StringContains(t, string(html), "multiple")
				assert.StringContains(t, string(html), `<option value="fr" selected>Français</option>`)
				assert.StringContains(t, string(html), `<option value="en">Anglais</option>`)
			})

			t.Run("Checkbox", func(t *testing.T) {
				f := Checkbox("newsletter")
				f.Value = true

				html, err := r.Field(f, th.opts)

				assert.NoError(t, err)
				assert.StringContains(t, string(html), `type="checkbox"`)
				assert.StringContains(t, string(html), "checked")
			})

			t.Run("CheckboxGroup", func(t *testing.T) {
				f := CheckboxGroup("interests",
					Option{Value: "go", Label: "Go"},
					Option{Value: "rust", Label: "Rust"},
				)
				f.Value = []string{"go"}

				html, err := r.Field(f, th.opts)

				assert.NoError(t, err)
				assert.Equal(t, 2, strings.Count(string(html), `type="checkbox"`))
				assert.Equal(t, 1, strings.Count(string(html), "checked"))
				assert.StringContains(t, string(html), `name="interests"`)
			})

			t.Run("CheckboxUnchecked", func(t *testing.T) {
				f := Checkbox("newsletter")

				html, err := r.Field(f, th.opts)

				assert.NoError(t, err)
				assert.False(t, strings.Contains(string(html), "checked"))
			})

			t.Run("Radio", func(t *testing.T) {
				f := Select("plan",
					Option{Value: "free", Label: "Free"},
					Option{Value: "pro", Label: "Pro"},
				)
				f.Type = RadioType
				f.Value = "pro"

				html, err := r.Field(f, th.opts)

				assert.NoError(t, err)
				assert.StringContains(t, string(html), `type="radio"`)
				assert.StringContains(t, string(html), `id="plan-free"`)
				assert.StringContains(t, string(html), `id="plan-pro"`)
				assert.StringContains(t, string(html), "checked")
			})

			t.Run("Hidden", func(t *testing.T) {
				f := Hidden("csrf").WithValue("token-123")

				html, err := r.Field(f, th.opts)

				assert.NoError(t, err)
				assert.StringContains(t, string(html), `type="hidden"`)
				assert.StringContains(t, string(html), `value="token-123"`)
			})

			t.Run("Submit", func(t *testing.T) {
				f := Submit("submit").WithLabel("Créer le compte")

				html, err := r.Field(f, th.opts)

				assert.NoError(t, err)
				assert.StringContains(t, string(html), `type="submit"`)
				assert.StringContains(t, string(html), "Créer le compte")
			})

			t.Run("SubmitRowHasNoLabel", func(t *testing.T) {
				f := Submit("submit").WithLabel("Créer le compte")

				html, err := r.Row(f, th.opts)

				assert.NoError(t, err)
				assert.False(t, strings.Contains(string(html), "<label"))
				assert.StringContains(t, string(html), `type="submit"`)
			})

			t.Run("InputWithErrorsHasAriaAttrs", func(t *testing.T) {
				f := Text("username")
				f.Errors = []string{"required"}

				html, err := r.Row(f, th.opts)

				assert.NoError(t, err)
				assert.StringContains(t, string(html), `aria-invalid="true"`)
				assert.StringContains(t, string(html), `aria-describedby="username-errors"`)
				assert.StringContains(t, string(html), `id="username-errors"`)
			})

			t.Run("InputWithoutErrorsHasNoAriaAttrs", func(t *testing.T) {
				f := Text("username")

				html, err := r.Field(f, th.opts)

				assert.NoError(t, err)
				assert.False(t, strings.Contains(string(html), "aria-invalid"))
				assert.False(t, strings.Contains(string(html), "aria-describedby"))
			})

			t.Run("RadioRendersAsFieldsetWithLegendAndNoSeparateLabel", func(t *testing.T) {
				f := Select("plan", Option{Value: "free", Label: "Free"}, Option{Value: "pro", Label: "Pro"}).
					WithLabel("Choisissez un plan")
				f.Type = RadioType

				html, err := r.Row(f, th.opts)

				assert.NoError(t, err)
				assert.StringContains(t, string(html), "<fieldset")
				assert.StringContains(t, string(html), "<legend")
				assert.StringContains(t, string(html), "Choisissez un plan")
				assert.False(t, strings.Contains(string(html), `<label for="plan">`))
			})

			t.Run("CheckboxGroupRendersAsFieldsetWithLegend", func(t *testing.T) {
				f := CheckboxGroup("interests",
					Option{Value: "go", Label: "Go"},
					Option{Value: "rust", Label: "Rust"},
				).WithLabel("Centres d'intérêt")

				html, err := r.Row(f, th.opts)

				assert.NoError(t, err)
				assert.StringContains(t, string(html), "<fieldset")
				assert.StringContains(t, string(html), "<legend")
				assert.False(t, strings.Contains(string(html), `<label for="interests">`))
			})
		})
	}

	t.Run("Bootstrap5UsesFormSelectClass", func(t *testing.T) {
		f := Select("country", Option{Value: "fr", Label: "France"})

		html, err := r.Field(f, render.Options{Name: "bootstrap5"})

		assert.NoError(t, err)
		assert.StringContains(t, string(html), `class="form-select"`)
	})

	t.Run("Bootstrap5UsesFormCheckInputClassForCheckboxAndRadio", func(t *testing.T) {
		checkbox := Checkbox("newsletter")
		html, err := r.Field(checkbox, render.Options{Name: "bootstrap5"})
		assert.NoError(t, err)
		assert.StringContains(t, string(html), `class="form-check-input"`)

		radio := Select("plan", Option{Value: "free", Label: "Free"})
		radio.Type = RadioType
		html, err = r.Field(radio, render.Options{Name: "bootstrap5"})
		assert.NoError(t, err)
		assert.StringContains(t, string(html), `class="form-check-input"`)
	})
}

func TestRenderer_Render(t *testing.T) {
	r := NewRenderer("default", ThemeDefault())

	t.Run("RendersField", func(t *testing.T) {
		var w bytes.Buffer
		f := Text("username")

		err := r.Render(context.Background(), &w, f, render.NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, w.String(), `name="username"`)
	})

	t.Run("UnknownThemeReturnsError", func(t *testing.T) {
		var w bytes.Buffer

		err := r.Render(context.Background(), &w, Text("username"), render.Options{Name: "bogus"})

		assert.Error(t, err)
	})

	t.Run("RendersForm", func(t *testing.T) {
		var w bytes.Buffer
		form := New("login").Add(Text("username"))

		err := r.Render(context.Background(), &w, form, render.NoOptions)

		assert.NoError(t, err)
		assert.StringContains(t, w.String(), `<form`)
	})

	t.Run("UnsupportedDataTypeReturnsError", func(t *testing.T) {
		var w bytes.Buffer

		err := r.Render(context.Background(), &w, "not a field", render.NoOptions)

		assert.Error(t, err)
	})

	t.Run("RespectsContextCancellation", func(t *testing.T) {
		var w bytes.Buffer
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := r.Render(ctx, &w, Text("username"), render.NoOptions)

		assert.ErrorIs(t, err, context.Canceled)
	})
}
