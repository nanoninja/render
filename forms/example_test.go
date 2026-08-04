// Copyright 2026 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forms_test

import (
	"fmt"

	"github.com/nanoninja/render"
	"github.com/nanoninja/render/forms"
)

func ExampleRenderer_Field() {
	r := forms.NewRenderer("default", forms.ThemeDefault())

	f := forms.Text("username").WithRequired()

	html, _ := r.Field(f, render.NoOptions)
	fmt.Println(html)

	// Output:
	// <input type="text" name="username" id="username" required>
}

func ExampleRenderer_Row() {
	r := forms.NewRenderer("default", forms.ThemeDefault())

	f := forms.Email("email").WithLabel("Email")
	f.Errors = []string{"required"}

	html, _ := r.Row(f, render.NoOptions)
	fmt.Println(html)

	// Output:
	// <div class="form-row">
	// <label for="email">Email</label>
	// <input type="email" name="email" id="email" aria-invalid="true" aria-describedby="email-errors">
	// <div id="email-errors"><div class="form-error">required</div></div>
	// </div>
}

func ExampleRenderer_Form() {
	r := forms.NewRenderer("default", forms.ThemeDefault())

	login := forms.New("login").
		WithAction("/login").
		Add(
			forms.Text("username").WithLabel("Username").WithRequired(),
			forms.Submit("submit").WithLabel("Sign in"),
		)

	html, _ := r.Form(login, render.NoOptions)
	fmt.Println(html)

	// Output:
	// <form name="login" method="POST" action="/login"><div class="form-row">
	// <label for="username">Username *</label>
	// <input type="text" name="username" id="username" required>
	//
	// </div><div class="form-row">
	//
	// <button type="submit" name="submit">Sign in</button>
	//
	// </div></form>
}

func ExampleForm_Fill() {
	f := forms.New("login").Add(forms.Text("username"))

	f.Fill(map[string][]string{"username": {"jdoe"}}, nil)

	fmt.Println(f.Field("username").Value)

	// Output:
	// jdoe
}
