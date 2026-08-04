// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forms

import (
	"html/template"
	"testing"

	"github.com/nanoninja/assert"
)

// incompleteTheme defines only a subset of the blocks required by
// the Theme contract, to exercise validateTheme's failure path.
type incompleteTheme struct{}

func (incompleteTheme) Templates() *template.Template {
	return template.Must(template.New("incomplete").Parse(`{{define "field_input"}}<input>{{end}}`))
}

func TestValidateTheme(t *testing.T) {
	t.Run("PassesForACompleteTheme", func(t *testing.T) {
		err := validateTheme(ThemeDefault().Templates())
		assert.NoError(t, err)
	})

	t.Run("PassesForBootstrap5", func(t *testing.T) {
		err := validateTheme(ThemeBootstrap5().Templates())
		assert.NoError(t, err)
	})

	t.Run("FailsAndNamesEveryMissingBlock", func(t *testing.T) {
		err := validateTheme(incompleteTheme{}.Templates())

		assert.Error(t, err)
		assert.StringContains(t, err.Error(), "field_select")
		assert.StringContains(t, err.Error(), "row")
		assert.StringContains(t, err.Error(), "form_start")
	})
}

func TestRenderer_RegisterTheme(t *testing.T) {
	t.Run("RejectsAnIncompleteTheme", func(t *testing.T) {
		r := NewRenderer("default", ThemeDefault())

		err := r.RegisterTheme("broken", incompleteTheme{})

		assert.Error(t, err)
		assert.StringContains(t, err.Error(), `"broken"`)
	})

	t.Run("DoesNotRegisterAnInvalidTheme", func(t *testing.T) {
		r := NewRenderer("default", ThemeDefault())

		_ = r.RegisterTheme("broken", incompleteTheme{})
		_, err := r.theme("broken")

		assert.Error(t, err)
	})
}

func TestNewRenderer_PanicsOnIncompleteTheme(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r)
	}()

	NewRenderer("broken", incompleteTheme{})

	t.Fatal("expected NewRenderer to panic on an incomplete theme")
}
