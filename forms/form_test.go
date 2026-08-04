// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forms

import (
	"testing"

	"github.com/nanoninja/assert"
)

func TestNew(t *testing.T) {
	f := New("login")

	assert.Equal(t, "login", f.Name)
	assert.Equal(t, "POST", f.Method)
	assert.Equal(t, 0, len(f.Fields))
	assert.Nil(t, f.Field("missing"))
}

func TestForm_Add(t *testing.T) {
	t.Run("AppendsFieldsInOrder", func(t *testing.T) {
		username := Text("username")
		email := Email("email")

		f := New("login").Add(username, email)

		assert.Equal(t, 2, len(f.Fields))
		assert.Equal(t, username, f.Fields[0])
		assert.Equal(t, email, f.Fields[1])
	})

	t.Run("IndexesFieldsByTheirOwnName", func(t *testing.T) {
		username := Text("username")
		email := Email("email")

		f := New("login").Add(username, email)

		assert.Equal(t, username, f.Field("username"))
		assert.Equal(t, email, f.Field("email"))
	})

	t.Run("SupportsMultipleCalls", func(t *testing.T) {
		f := New("login").Add(Text("username"))
		f.Add(Email("email"))

		assert.Equal(t, 2, len(f.Fields))
	})

	t.Run("ReturnsSameFormForChaining", func(t *testing.T) {
		f := New("login")
		got := f.Add(Text("username"))
		assert.Equal(t, f, got)
	})
}

func TestForm_Field(t *testing.T) {
	f := New("login").Add(Text("username"))

	assert.NotNil(t, f.Field("username"))
	assert.Nil(t, f.Field("nonexistent"))
}

func TestForm_Group(t *testing.T) {
	t.Run("ReturnsFieldsInDeclarationOrder", func(t *testing.T) {
		firstName := Text("first_name").WithGroup("left")
		lastName := Text("last_name").WithGroup("left")
		email := Email("email").WithGroup("right")

		f := New("register").Add(firstName, email, lastName)

		got := f.Group("left")

		assert.Equal(t, 2, len(got))
		assert.Equal(t, firstName, got[0])
		assert.Equal(t, lastName, got[1])
	})

	t.Run("EmptyWhenNoFieldMatches", func(t *testing.T) {
		f := New("register").Add(Text("first_name"))

		assert.Equal(t, 0, len(f.Group("left")))
	})

	t.Run("UngroupedFieldsHaveEmptyGroup", func(t *testing.T) {
		f := New("register").Add(Text("first_name"))

		assert.Equal(t, 1, len(f.Group("")))
	})
}

func TestForm_With(t *testing.T) {
	t.Run("WithAction", func(t *testing.T) {
		f := New("login").WithAction("/login")
		assert.Equal(t, "/login", f.Action)
	})

	t.Run("WithMethod", func(t *testing.T) {
		f := New("login").WithMethod("GET")
		assert.Equal(t, "GET", f.Method)
	})
}

func TestForm_Enctype(t *testing.T) {
	t.Run("EmptyWithoutFileField", func(t *testing.T) {
		f := New("login").Add(Text("username"), Email("email"))
		assert.Equal(t, "", f.Enctype())
	})

	t.Run("MultipartWhenFileFieldPresent", func(t *testing.T) {
		f := New("upload").Add(Text("title"), File("avatar"))
		assert.Equal(t, "multipart/form-data", f.Enctype())
	})

	t.Run("MultipartRegardlessOfFieldPosition", func(t *testing.T) {
		f := New("upload").Add(File("avatar"), Text("title"))
		assert.Equal(t, "multipart/form-data", f.Enctype())
	})
}

func TestForm_Fill(t *testing.T) {
	t.Run("PopulatesFieldValue", func(t *testing.T) {
		f := New("login").Add(Text("username"))

		f.Fill(map[string][]string{"username": {"jdoe"}}, nil)

		assert.Equal(t, "jdoe", f.Field("username").Value)
	})

	t.Run("PopulatesMultipleValueAsSlice", func(t *testing.T) {
		langs := Select("langs").WithMultiple()
		f := New("prefs").Add(langs)

		f.Fill(map[string][]string{"langs": {"fr", "de"}}, nil)

		assert.Equal(t, []string{"fr", "de"}, langs.Value.([]string))
	})

	t.Run("PopulatesCheckboxGroupValueAsSlice", func(t *testing.T) {
		interests := CheckboxGroup("interests", Option{Value: "go", Label: "Go"}, Option{Value: "rust", Label: "Rust"})
		f := New("profile").Add(interests)

		f.Fill(map[string][]string{"interests": {"go", "rust"}}, nil)

		assert.Equal(t, []string{"go", "rust"}, interests.Value.([]string))
	})

	t.Run("LeavesFieldsAbsentFromValuesUntouched", func(t *testing.T) {
		username := Text("username").WithValue("previous")
		f := New("login").Add(username)

		f.Fill(map[string][]string{}, nil)

		assert.Equal(t, "previous", username.Value)
	})

	t.Run("EmptySubmittedValueClearsField", func(t *testing.T) {
		username := Text("username").WithValue("previous")
		f := New("login").Add(username)

		f.Fill(map[string][]string{"username": {}}, nil)

		assert.Equal(t, "", username.Value)
	})

	t.Run("PopulatesFieldErrors", func(t *testing.T) {
		f := New("login").Add(Text("username"))

		f.Fill(nil, map[string][]string{"username": {"required"}})

		assert.Equal(t, []string{"required"}, f.Field("username").Errors)
	})

	t.Run("PopulatesFormLevelErrors", func(t *testing.T) {
		f := New("login").Add(Text("username"))

		f.Fill(nil, nil, "too many attempts")

		assert.Equal(t, []string{"too many attempts"}, f.Errors)
	})

	t.Run("NeverWritesIntoAttrs", func(t *testing.T) {
		username := Text("username").WithAttr("data-test", "1")
		f := New("login").Add(username)

		f.Fill(map[string][]string{"username": {"jdoe"}}, nil)

		assert.Equal(t, "1", username.Attrs["data-test"])
		assert.Equal(t, 1, len(username.Attrs))
	})
}
