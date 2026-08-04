// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This example runs a net/http server serving a 5-question quiz on
// a single page: each question is a Radio field, the whole form is
// rendered with one {{ form .Form }} call (Submit is just another
// field), and the score is computed and displayed on the very same
// page after submit — no redirect, no separate result page.
//
// Run it with:
//
//	go run ./_example/quiz
package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/nanoninja/render"
	"github.com/nanoninja/render/forms"
	"github.com/nanoninja/render/tmpl"
	"github.com/nanoninja/render/tmpl/loader"
)

//go:embed templates/*.html
var templatesFS embed.FS

// question pairs a Radio field's options with its correct answer,
// kept server-side only — never sent to the client ahead of time.
type question struct {
	Name    string
	Text    string
	Options []forms.Option
	Correct string
}

var questions = []question{
	{
		Name: "q1",
		Text: "Which Go keyword declares a variable with an inferred type?",
		Options: []forms.Option{
			{Value: "let", Label: "let"},
			{Value: "var", Label: "var"},
			{Value: "auto", Label: "auto"},
			{Value: "def", Label: "def"},
		},
		Correct: "var",
	},
	{
		Name: "q2",
		Text: "Which keyword starts a goroutine?",
		Options: []forms.Option{
			{Value: "async", Label: "async"},
			{Value: "spawn", Label: "spawn"},
			{Value: "go", Label: "go"},
			{Value: "fork", Label: "fork"},
		},
		Correct: "go",
	},
	{
		Name: "q3",
		Text: "Which standard library package handles HTML with automatic escaping?",
		Options: []forms.Option{
			{Value: "text/template", Label: "text/template"},
			{Value: "html/template", Label: "html/template"},
			{Value: "net/html", Label: "net/html"},
			{Value: "fmt", Label: "fmt"},
		},
		Correct: "html/template",
	},
	{
		Name: "q4",
		Text: "What's the communication primitive between goroutines called?",
		Options: []forms.Option{
			{Value: "pipe", Label: "pipe"},
			{Value: "queue", Label: "queue"},
			{Value: "channel", Label: "channel"},
			{Value: "socket", Label: "socket"},
		},
		Correct: "channel",
	},
	{
		Name: "q5",
		Text: "Which interface must a type satisfy to customize its fmt.Println output?",
		Options: []forms.Option{
			{Value: "Printer", Label: "Printer"},
			{Value: "Stringer", Label: "Stringer"},
			{Value: "Displayer", Label: "Displayer"},
			{Value: "Formatter", Label: "Formatter"},
		},
		Correct: "Stringer",
	},
}

// newQuizForm builds a fresh quiz form. Called once per request —
// Form/Field carry per-request state (Value, Errors).
func newQuizForm() *forms.Form {
	f := forms.New("quiz").WithAction("/quiz")
	for _, q := range questions {
		f.Add(forms.Radio(q.Name, q.Options...).WithLabel(q.Text).WithRequired())
	}
	f.Add(forms.Submit("submit").WithLabel("Submit my answers"))
	return f
}

// score counts correct answers, and reports any unanswered question
// by name so the form can be redisplayed with a clear error instead
// of silently scoring a blank as wrong.
func score(r *http.Request) (correct int, unanswered []string) {
	for _, q := range questions {
		answer := r.PostForm.Get(q.Name)
		if answer == "" {
			unanswered = append(unanswered, q.Name)
			continue
		}
		if answer == q.Correct {
			correct++
		}
	}
	return correct, unanswered
}

func main() {
	r := forms.NewRenderer("bootstrap5", forms.ThemeBootstrap5())

	src := loader.NewEmbed(templatesFS, tmpl.LoaderConfig{
		Root:      "templates",
		Extension: ".html",
	})

	t := tmpl.HTML("quiz", tmpl.SetFuncsHTML(forms.Funcs(r)))
	if err := t.Load(src); err != nil {
		log.Fatal(err)
	}

	type pageData struct {
		Form      *forms.Form
		Submitted bool
		Score     int
		Total     int
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, req *http.Request) {
		data := pageData{Form: newQuizForm(), Total: len(questions)}
		if err := t.Render(req.Context(), w, data, render.Options{Name: "quiz.html"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("POST /quiz", func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		correct, unanswered := score(req)

		form := newQuizForm()

		data := pageData{Form: form, Total: len(questions)}

		if len(unanswered) > 0 {
			fieldErrs := make(map[string][]string, len(unanswered))
			for _, name := range unanswered {
				fieldErrs[name] = []string{"Please answer this question."}
			}
			form.Fill(req.PostForm, fieldErrs, "Please answer every question.")
		} else {
			form.Fill(req.PostForm, nil)
			data.Submitted = true
			data.Score = correct
		}

		if err := t.Render(req.Context(), w, data, render.Options{Name: "quiz.html"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
