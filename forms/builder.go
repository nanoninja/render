// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forms

// Text creates a single-line text input field.
func Text(name string) *Field {
	return &Field{Name: name, Type: InputType, HTMLType: "text"}
}

// Email creates an email input field.
func Email(name string) *Field {
	return &Field{Name: name, Type: InputType, HTMLType: "email"}
}

// File creates a file upload input field.
func File(name string) *Field {
	return &Field{Name: name, Type: InputType, HTMLType: "file"}
}

// Password creates a password input field.
func Password(name string) *Field {
	return &Field{Name: name, Type: InputType, HTMLType: "password"}
}

// Date creates a date input field.
func Date(name string) *Field {
	return &Field{Name: name, Type: InputType, HTMLType: "date"}
}

// Number creates a number input field.
func Number(name string) *Field {
	return &Field{Name: name, Type: InputType, HTMLType: "number"}
}

// Tel creates a telephone number input field.
func Tel(name string) *Field {
	return &Field{Name: name, Type: InputType, HTMLType: "tel"}
}

// URL creates a URL input field.
func URL(name string) *Field {
	return &Field{Name: name, Type: InputType, HTMLType: "url"}
}

// Color creates a color picker input field.
func Color(name string) *Field {
	return &Field{Name: name, Type: InputType, HTMLType: "color"}
}

// Range creates a range slider input field.
func Range(name string) *Field {
	return &Field{Name: name, Type: InputType, HTMLType: "range"}
}

// Time creates a time input field.
func Time(name string) *Field {
	return &Field{Name: name, Type: InputType, HTMLType: "time"}
}

// Search creates a search input field.
func Search(name string) *Field {
	return &Field{Name: name, Type: InputType, HTMLType: "search"}
}

// Textarea creates a multi-line text field.
func Textarea(name string) *Field {
	return &Field{Name: name, Type: TextareaType}
}

// Select creates a select field with the given options.
func Select(name string, options ...Option) *Field {
	return &Field{Name: name, Type: SelectType, Options: options}
}

// Checkbox creates a checkbox field.
func Checkbox(name string) *Field {
	return &Field{Name: name, Type: CheckboxType}
}

// CheckboxGroup creates a group of independent checkboxes sharing
// the same name — e.g. "select all that apply". Value is expected
// to be a []string (Multiple is set automatically, like a
// multi-select), and Field.IsSelected reports whether a given
// option is checked.
func CheckboxGroup(name string, options ...Option) *Field {
	return &Field{Name: name, Type: CheckboxGroupType, Options: options, Multiple: true}
}

// Hidden creates a hidden field.
func Hidden(name string) *Field {
	return &Field{Name: name, Type: HiddenType}
}

// Submit creates a submit button field. Label sets the button's
// visible text.
func Submit(name string) *Field {
	return &Field{Name: name, Type: SubmitType}
}
