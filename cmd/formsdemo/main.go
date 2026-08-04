// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command formsdemo is a scratch net/http server demoing the forms
// sub-package end to end: field rendering, theme switching, layout
// groups, and redisplay after a failed submit via Form.Fill.
package main

import (
	"embed"
	"log"
	"net/http"
	"strings"

	"github.com/nanoninja/render"
	"github.com/nanoninja/render/forms"
	"github.com/nanoninja/render/tmpl"
	"github.com/nanoninja/render/tmpl/loader"
)

//go:embed templates/*.html
var templatesFS embed.FS

const maxUploadSize = 10 << 20 // 10 MB

// newRegisterForm builds a fresh registration form. It must be
// called once per request — Form/Field carry per-request state
// (Value, Errors) and are not meant to be shared or cached.
//
// Fields are tagged with WithGroup into "left"/"right" columns to
// exercise Form.Group/Renderer.Group/form_group end to end; the
// remaining fields (country, newsletter, avatar, submit) stay
// ungrouped and are rendered individually below the two columns.
func newRegisterForm() *forms.Form {
	return forms.New("register").
		WithAction("/register").
		Add(
			forms.Text("name").WithLabel("Nom complet").WithRequired().WithGroup("left"),
			forms.Email("email").WithLabel("Email").WithRequired().WithGroup("left"),
			forms.Password("password").WithLabel("Mot de passe").WithRequired().WithGroup("right"),
			forms.Password("password_confirm").WithLabel("Confirmer le mot de passe").WithRequired().WithGroup("right"),
			forms.Select("country",
				forms.Option{Value: "", Label: "-- Choisir --"},
				forms.Option{Value: "fr", Label: "France"},
				forms.Option{Value: "be", Label: "Belgique"},
				forms.Option{Value: "ch", Label: "Suisse"},
			).WithLabel("Pays").WithRequired(),
			forms.Checkbox("newsletter").WithLabel("Recevoir la newsletter"),
			forms.File("avatar").WithLabel("Photo de profil (optionnel)"),
			forms.Submit("submit").WithLabel("Créer le compte"),
		)
}

// validateRegister applies minimal, hand-written validation rules.
// This is exactly the kind of logic that would live in the future
// separate binding/validation module — done inline here only
// because this is a scratch demo, not a reusable library.
func validateRegister(r *http.Request) (fieldErrs map[string][]string, formErrs []string) {
	fieldErrs = make(map[string][]string)

	if r.PostForm.Get("name") == "" {
		fieldErrs["name"] = []string{"Le nom est requis."}
	}

	email := r.PostForm.Get("email")
	switch {
	case email == "":
		fieldErrs["email"] = []string{"L'email est requis."}
	case !strings.Contains(email, "@"):
		fieldErrs["email"] = []string{"L'email n'est pas valide."}
	}

	password := r.PostForm.Get("password")
	if len(password) < 8 {
		fieldErrs["password"] = []string{"8 caractères minimum."}
	}
	if password != r.PostForm.Get("password_confirm") {
		fieldErrs["password_confirm"] = []string{"Les mots de passe ne correspondent pas."}
	}

	if r.PostForm.Get("country") == "" {
		fieldErrs["country"] = []string{"Merci de choisir un pays."}
	}

	return fieldErrs, formErrs
}

func main() {
	r := forms.NewRenderer("bootstrap5", forms.ThemeBootstrap5())

	src := loader.NewEmbed(templatesFS, tmpl.LoaderConfig{
		Root:      "templates",
		Extension: ".html",
	})

	t := tmpl.HTML("formsweb", tmpl.SetFuncsHTML(forms.Funcs(r)))
	if err := t.Load(src); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, req *http.Request) {
		data := struct {
			Form    *forms.Form
			Success bool
		}{Form: newRegisterForm()}

		if err := t.Render(req.Context(), w, data, render.Options{Name: "register.html"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Stays on the registration page after submit, success or not —
	// re-fills the form from what was submitted so we can visually
	// check that Fill correctly round-trips every field type
	// (including the grouped columns and the select/checkbox/file
	// fields), instead of navigating away to a separate page.
	mux.HandleFunc("POST /register", func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseMultipartForm(maxUploadSize); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		fieldErrs, formErrs := validateRegister(req)

		form := newRegisterForm()
		form.Fill(req.PostForm, fieldErrs, formErrs...)

		data := struct {
			Form    *forms.Form
			Success bool
		}{
			Form:    form,
			Success: len(fieldErrs) == 0 && len(formErrs) == 0,
		}

		if err := t.Render(req.Context(), w, data, render.Options{Name: "register.html"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
