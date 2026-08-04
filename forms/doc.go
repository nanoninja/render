// Package forms provides a Symfony-inspired form rendering system on top
// of the render.Renderer interface and the tmpl sub-package.
//
// Fields are described in plain Go, independently of any markup:
//
//	f := forms.Text("username").WithLabel("Username").WithRequired()
//	form := forms.New("login").WithAction("/login").Add(f, forms.Submit("submit"))
//
// A [Renderer] turns fields and forms into HTML through a swappable
// [Theme] — [ThemeDefault] ships dependency-free, [ThemeBootstrap5] targets
// Bootstrap 5. Both themes implement the same named-block contract, so
// switching theme never requires touching calling code:
//
//	r := forms.NewRenderer("default", forms.ThemeDefault())
//	r.RegisterTheme("bootstrap5", forms.ThemeBootstrap5())
//
//	html, err := r.Form(form, render.NoOptions)
//	html, err = r.Form(form, render.Options{Name: "bootstrap5"})
//
// [Renderer] also implements [render.Renderer], so it composes with
// render.Buffer, render.Cache, render.Gzip, and render.Multi, and covers
// the case where a form is the entire HTTP response (e.g. an HTMX partial
// swap). For embedding inside a larger page template, use [Funcs] to
// expose form, form_row, form_widget, form_group, and related helpers to
// a tmpl.HTML template's function map:
//
//	t := tmpl.HTML("myapp", tmpl.SetFuncsHTML(forms.Funcs(r)))
//
// Then, in the template:
//
//	{{ form .Form }}
//
// Scope is rendering only. Binding an *http.Request into submitted
// values, and validating them, are deliberately out of scope — see
// [Form.Fill] for the one narrow, mechanical bridge this package
// provides for redisplaying a form with the user's own input and
// validation errors after a failed submit.
package forms
