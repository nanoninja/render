// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forms

import (
	"html/template"
	"maps"
	"slices"
	"strings"
)

// FieldType identifies the kind of HTML widget a Field renders as.
type FieldType int

// blockName returns the template block name that renders this
// field type, per the Theme contract.
func (t FieldType) blockName() string {
	switch t {
	case TextareaType:
		return "field_textarea"
	case SelectType:
		return "field_select"
	case CheckboxType:
		return "field_checkbox"
	case CheckboxGroupType:
		return "field_checkboxes"
	case RadioType:
		return "field_radio"
	case HiddenType:
		return "field_hidden"
	case SubmitType:
		return "field_submit"
	default:
		return "field_input"
	}
}

// FieldType values, one per supported widget family.
const (
	InputType FieldType = iota
	TextareaType
	SelectType
	CheckboxType
	CheckboxGroupType
	RadioType
	HiddenType
	SubmitType
)

// Option represents a single choice in a Select or Radio field.
type Option struct {
	Value string
	Label string
}

// Field describes a single form field, independently of any theme.
// It carries the widget family (Type), the underlying data (Value),
// and the metadata (Label, constraints, errors) a renderer needs to
// produce HTML — without knowing how that HTML is styled.
//
// Type selects which template block renders the field (input,
// textarea, select...). HTMLType refines Type when it is Input,
// giving the literal HTML5 input type ("text", "email", "date",
// "number"...); it is ignored for every other Type.
//
// Options is only meaningful for Select and Radio. Errors is filled
// by validation performed elsewhere — Field itself does not validate.
type Field struct {
	Name        string
	Label       string
	Placeholder string
	HTMLType    string
	Group       string // arbitrary layout tag, e.g. "personal", "step1", "left"
	Type        FieldType
	Required    bool
	Multiple    bool
	Value       any
	Options     []Option
	Errors      []string
	Attrs       map[string]string
}

// WithLabel sets the field's display label.
func (f *Field) WithLabel(label string) *Field {
	f.Label = label
	return f
}

// WithValue sets the field's current value.
func (f *Field) WithValue(v any) *Field {
	f.Value = v
	return f
}

// WithPlaceholder sets the field's placeholder text.
func (f *Field) WithPlaceholder(p string) *Field {
	f.Placeholder = p
	return f
}

// WithRequired marks the field as required.
func (f *Field) WithRequired() *Field {
	f.Required = true
	return f
}

// WithAttr sets a single HTML attribute on the field.
func (f *Field) WithAttr(key, value string) *Field {
	if f.Attrs == nil {
		f.Attrs = make(map[string]string)
	}
	f.Attrs[key] = value
	return f
}

// WithClass adds one or more CSS classes to the field, appending to
// any classes already set rather than replacing them.
func (f *Field) WithClass(classes ...string) *Field {
	existing := f.Attrs["class"]
	for _, c := range classes {
		if existing == "" {
			existing = c
		} else {
			existing += " " + c
		}
	}
	return f.WithAttr("class", existing)
}

// WithStyle sets the field's inline CSS style, replacing any
// previous value.
func (f *Field) WithStyle(style string) *Field {
	return f.WithAttr("style", style)
}

// WithMultiple marks a select field as allowing multiple selections.
// Value is then expected to be a []string.
func (f *Field) WithMultiple() *Field {
	f.Multiple = true
	return f
}

// WithGroup tags the field with an arbitrary group name, e.g. for a
// two-column layout ("left"/"right") or a multi-step wizard
// ("step1"/"step2"). See Form.Group and Renderer.Group.
func (f *Field) WithGroup(name string) *Field {
	f.Group = name
	return f
}

// IsSelected reports whether value is currently selected for this
// field. It honors Multiple (Value expected as []string) as well
// as the single-value case (Value expected as string).
func (f *Field) IsSelected(value string) bool {
	if f.Multiple {
		values, ok := f.Value.([]string)
		if !ok {
			return false
		}
		return slices.Contains(values, value)
	}
	s, ok := f.Value.(string)
	return ok && s == value
}

// IsSubmit reports whether this field is a submit button. Used by
// theme templates to skip rendering a <label> for it — its Label
// is the button's own text, not a form label.
func (f *Field) IsSubmit() bool {
	return f.Type == SubmitType
}

// IsGroup reports whether this field renders as a group of
// individual inputs sharing one name (RadioType, CheckboxGroupType).
// Used by theme templates to skip the usual <label>: these widgets
// wrap themselves in a <fieldset><legend> instead, using Label as
// the legend text.
func (f *Field) IsGroup() bool {
	return f.Type == RadioType || f.Type == CheckboxGroupType
}

// withOverrides returns a shallow copy of f with the given
// attr/value pairs merged into its Attrs, without mutating f. Used
// to apply per-call rendering overrides (Symfony-style
// form_widget(field, {'class': '...'})) without touching the
// original field.
//
// Each pair replaces the corresponding Attrs entry outright — unlike
// WithClass, this does not accumulate onto an existing "class". Any
// base class baked into a theme's widget template (e.g. Bootstrap's
// "form-control") is unaffected either way, since themes never read
// it back from Attrs.
func withOverrides(f *Field, pairs ...string) *Field {
	if len(pairs) == 0 {
		return f
	}
	clone := *f
	clone.Attrs = make(map[string]string, len(f.Attrs)+len(pairs)/2)
	maps.Copy(clone.Attrs, f.Attrs)
	for i := 0; i+1 < len(pairs); i += 2 {
		clone.Attrs[pairs[i]] = pairs[i+1]
	}
	return &clone
}

// RenderAttrs returns Attrs rendered as a single, pre-escaped HTML
// attribute list (e.g. ` data-test="hello" aria-label="Close"`),
// skipping any key named in except (used by themes that merge
// "class" separately). Keys are sorted for deterministic output.
//
// This exists because html/template cannot safely escape a
// template action used as a dynamic attribute *name* — a theme
// template ranging over Attrs directly (`{{$k}}="{{$v}}"`) makes
// html/template poison the value with "ZgotmplZ" instead of
// rendering it, since it can't verify the escaping rules for an
// unknown attribute name at parse time. Doing the escaping here in
// Go instead sidesteps that: RenderAttrs only ever reads from
// Attrs, which stays server-controlled by design (see
// withOverrides and Form.Fill's doc comments), so there is no
// attacker-controlled attribute name to guard against in the first
// place — this is purely a template-engine limitation to work
// around, not an added security boundary.
func (f *Field) RenderAttrs(except ...string) template.HTMLAttr {
	if len(f.Attrs) == 0 {
		return ""
	}
	keys := make([]string, 0, len(f.Attrs))
	for k := range f.Attrs {
		if slices.Contains(except, k) {
			continue
		}
		keys = append(keys, k)
	}
	slices.Sort(keys)

	var b strings.Builder
	for _, k := range keys {
		b.WriteByte(' ')
		b.WriteString(template.HTMLEscapeString(k))
		b.WriteString(`="`)
		b.WriteString(template.HTMLEscapeString(f.Attrs[k]))
		b.WriteString(`"`)
	}
	return template.HTMLAttr(b.String())
}
