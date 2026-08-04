// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This example runs a net/http server serving a 3-step registration
// wizard, using the dependency-free default theme (login and quiz
// both use bootstrap5). Each step is its own page; there is no
// server-side session — data already entered is carried forward as
// Hidden fields in the next step's form, and finally read back from
// the last POST to render the confirmation page. This is the "wizard"
// half of Field.WithGroup/Form.Group/Renderer.Group; the login
// example already covers its other half, two-column layouts.
//
// Run it with:
//
//	go run ./_example/wizard
package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/nanoninja/render"
	"github.com/nanoninja/render/forms"
	"github.com/nanoninja/render/tmpl"
	"github.com/nanoninja/render/tmpl/loader"
)

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed static/*.css
var staticFS embed.FS

var interestOptions = []forms.Option{
	{Value: "workshops", Label: "Workshops"},
	{Value: "networking", Label: "Networking"},
	{Value: "talks", Label: "Talks"},
	{Value: "expo", Label: "Sponsor expo"},
}

// step1Form asks for the attendee's identity.
func step1Form() *forms.Form {
	return forms.New("step1").
		WithAction("/wizard/step2").
		Add(
			forms.Text("name").WithLabel("Full name").WithRequired(),
			forms.Email("email").WithLabel("Email").WithRequired(),
			forms.Submit("next").WithLabel("Next"),
		)
}

// step2Form asks for preferences, carrying step 1's answers forward
// as hidden fields so nothing is lost between requests. The "back"
// button uses the HTML5 formaction attribute to submit this same
// form's hidden fields to a different route than the form's own
// Action — no JavaScript, no separate form needed to go back.
func step2Form(name, email string) *forms.Form {
	return forms.New("step2").
		WithAction("/wizard/step3").
		Add(
			forms.Hidden("name").WithValue(name),
			forms.Hidden("email").WithValue(email),
			forms.CheckboxGroup("interests", interestOptions...).WithLabel("What are you interested in?"),
			forms.Textarea("comment").WithLabel("Anything else you'd like us to know?"),
			forms.Submit("back").WithLabel("Previous").WithClass("secondary").
				WithAttr("formaction", "/wizard/back-to-step1").
				WithAttr("formnovalidate", "formnovalidate"),
			forms.Submit("next").WithLabel("Next"),
		)
}

// step3Form carries every answer so far forward as hidden fields,
// ready for the final confirmation.
func step3Form(name, email, comment string, interests []string) *forms.Form {
	f := forms.New("step3").
		WithAction("/wizard/complete").
		Add(
			forms.Hidden("name").WithValue(name),
			forms.Hidden("email").WithValue(email),
			forms.Hidden("comment").WithValue(comment),
		)
	for _, interest := range interests {
		f.Add(forms.Hidden("interests").WithValue(interest))
	}
	f.Add(
		forms.Submit("back").WithLabel("Previous").WithClass("secondary").
			WithAttr("formaction", "/wizard/back-to-step2").
			WithAttr("formnovalidate", "formnovalidate"),
		forms.Submit("confirm").WithLabel("Confirm registration"),
	)
	return f
}

// interestLabels turns submitted option values back into their
// display labels, for the confirmation page.
func interestLabels(values []string) []string {
	labels := make([]string, 0, len(values))
	for _, v := range values {
		for _, opt := range interestOptions {
			if opt.Value == v {
				labels = append(labels, opt.Label)
			}
		}
	}
	return labels
}

func main() {
	r := forms.NewRenderer("default", forms.ThemeDefault())

	src := loader.NewEmbed(templatesFS, tmpl.LoaderConfig{
		Root:      "templates",
		Extension: ".html",
	})

	t := tmpl.HTML("wizard", tmpl.SetFuncsHTML(forms.Funcs(r)))
	if err := t.Load(src); err != nil {
		log.Fatal(err)
	}

	renderPage := func(w http.ResponseWriter, req *http.Request, name string, data any) {
		if err := t.Render(req.Context(), w, data, render.Options{Name: name}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	staticDir, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}

	type stepData struct {
		Step     int
		Total    int
		Progress int
		Form     *forms.Form
	}

	newStepData := func(step, total int, form *forms.Form) stepData {
		return stepData{Step: step, Total: total, Progress: step * 100 / total, Form: form}
	}

	mux := http.NewServeMux()

	mux.Handle("GET /wizard/static/", http.StripPrefix("/wizard/static/", http.FileServerFS(staticDir)))

	mux.HandleFunc("GET /wizard/{$}", func(w http.ResponseWriter, req *http.Request) {
		renderPage(w, req, "step.html", newStepData(1, 3, step1Form()))
	})

	mux.HandleFunc("POST /wizard/back-to-step1", func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		form := step1Form()
		form.Fill(req.PostForm, nil)
		renderPage(w, req, "step.html", newStepData(1, 3, form))
	})

	mux.HandleFunc("POST /wizard/back-to-step2", func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		name, email := req.PostForm.Get("name"), req.PostForm.Get("email")

		form := step2Form(name, email)
		form.Fill(req.PostForm, nil)
		renderPage(w, req, "step.html", newStepData(2, 3, form))
	})

	mux.HandleFunc("POST /wizard/step2", func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		name, email := req.PostForm.Get("name"), req.PostForm.Get("email")

		fieldErrs := map[string][]string{}
		if name == "" {
			fieldErrs["name"] = []string{"Please enter your name."}
		}
		if email == "" {
			fieldErrs["email"] = []string{"Please enter your email."}
		}

		if len(fieldErrs) > 0 {
			form := step1Form()
			form.Fill(req.PostForm, fieldErrs)
			renderPage(w, req, "step.html", newStepData(1, 3, form))
			return
		}

		renderPage(w, req, "step.html", newStepData(2, 3, step2Form(name, email)))
	})

	mux.HandleFunc("POST /wizard/step3", func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		name := req.PostForm.Get("name")
		email := req.PostForm.Get("email")
		comment := req.PostForm.Get("comment")
		interests := req.PostForm["interests"]

		renderPage(w, req, "review.html", struct {
			Step      int
			Total     int
			Progress  int
			Form      *forms.Form
			Name      string
			Email     string
			Comment   string
			Interests []string
		}{
			Step: 3, Total: 3, Progress: 100,
			Form:      step3Form(name, email, comment, interests),
			Name:      name,
			Email:     email,
			Comment:   comment,
			Interests: interestLabels(interests),
		})
	})

	mux.HandleFunc("POST /wizard/complete", func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		renderPage(w, req, "complete.html", struct {
			Name      string
			Email     string
			Comment   string
			Interests []string
		}{
			Name:      req.PostForm.Get("name"),
			Email:     req.PostForm.Get("email"),
			Comment:   req.PostForm.Get("comment"),
			Interests: interestLabels(req.PostForm["interests"]),
		})
	})

	log.Println("listening on http://localhost:8080/wizard/")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
