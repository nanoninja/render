// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forms

// Form is an ordered collection of fields, plus the attributes of
// the HTML <form> element itself (action, method, encoding).
//
// Fields preserves declaration order for iteration in templates.
// Field looks a field up by name for direct access, e.g. when a
// template renders fields individually rather than looping over
// Fields.
type Form struct {
	Name   string
	Action string
	Method string
	Fields []*Field
	Errors []string
	byName map[string]*Field
}

// New creates an empty form with the given name and the default
// method "POST".
func New(name string) *Form {
	return &Form{
		Name:   name,
		Method: "POST",
		byName: make(map[string]*Field),
	}
}

// Add appends one or more fields to the form, in order.
func (f *Form) Add(fields ...*Field) *Form {
	for _, field := range fields {
		f.Fields = append(f.Fields, field)
		f.byName[field.Name] = field
	}
	return f
}

// Field returns the field registered under name, or nil if none
// exists.
func (f *Form) Field(name string) *Field {
	return f.byName[name]
}

// Group returns, in declaration order, the fields tagged with the
// given group name via WithGroup — for two-column layouts or
// multi-step wizards, where the caller controls the surrounding
// markup (see Renderer.Group).
func (f *Form) Group(name string) []*Field {
	var fields []*Field
	for _, field := range f.Fields {
		if field.Group == name {
			fields = append(fields, field)
		}
	}
	return fields
}

// WithAction sets the form's target URL.
func (f *Form) WithAction(action string) *Form {
	f.Action = action
	return f
}

// WithMethod sets the form's HTTP method.
func (f *Form) WithMethod(method string) *Form {
	f.Method = method
	return f
}

// Enctype returns the encoding the <form> tag must declare. It is
// "multipart/form-data" as soon as the form contains a file field,
// since HTML requires it for file uploads to work; otherwise it is
// empty, letting the browser default apply.
func (f *Form) Enctype() string {
	for _, field := range f.Fields {
		if field.Type == InputType && field.HTMLType == "file" {
			return "multipart/form-data"
		}
	}
	return ""
}

// Fill re-populates Value and Errors on each field from values and
// fieldErrs (both keyed by field name), and Form.Errors from the
// optional formErrs — used to redisplay a form with the user's own
// input and validation errors after a failed submit. Fields absent
// from values or fieldErrs are left untouched.
//
// Fill deliberately performs no sanitization: values are echoed
// back verbatim, and safety relies entirely on html/template's
// contextual auto-escaping at render time, not on filtering here.
// For the same reason, Fill never writes into Attrs — attribute
// names are a template context where dynamic, unescaped-by-design
// data is far riskier than in Value or Errors; Attrs must stay
// server-controlled only.
func (f *Form) Fill(values, errs map[string][]string, formErrs ...string) {
	for _, field := range f.Fields {
		if vs, ok := values[field.Name]; ok {
			switch {
			case field.Multiple:
				field.Value = vs
			case len(vs) > 0:
				field.Value = vs[0]
			default:
				field.Value = ""
			}
		}
		if msgs, ok := errs[field.Name]; ok {
			field.Errors = msgs
		}
	}
	f.Errors = formErrs
}
