package main

import (
	"errors"
	"net/http"
	"time"

	"git.01.alem.school/ggrks/forum.git/internal/forms"
	"git.01.alem.school/ggrks/forum.git/internal/models"
)

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/signup.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	switch r.Method {
	case http.MethodGet:
		form := forms.New(nil)
		app.render(w, r, files, "signup.page.tmpl", &templateData{
			Form: form,
		})
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			app.ErrorPage(w, 400)
			return
		}
		form := forms.New(r.PostForm)
		form.Required("name", "email", "password")
		form.MaxLength("name", 255)
		form.MaxLength("email", 255)
		form.MatchesPattern("email", forms.EmailRX)
		form.MinLength("password", 10)
		if !form.Valid() {
			app.render(w, r, files, "signup.page.tmpl", &templateData{
				Form: form,
			})
			return
		}
		err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
		if err != nil {
			if errors.Is(err, models.ErrDuplicateEmail) {
				form.Errors.Add("email", "Address is already in use")
				app.render(w, r, files, "signup.page.tmpl", &templateData{
					Form: form,
				})
				return
			} else {
				app.ErrorPage(w, 500)
			}
			return
		}
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	default:
		w.Header().Set("Allow", http.MethodPost)
		app.ErrorPage(w, 405)
		return

	}
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/login.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	switch r.Method {
	case http.MethodGet:
		form := forms.New(nil)
		app.render(w, r, files, "login.page.tmpl", &templateData{
			Form: form,
		})
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			app.ErrorPage(w, 405)
			return
		}
		form := forms.New(r.PostForm)
		form.Required("email", "password")
		form.MaxLength("email", 255)
		form.MatchesPattern("email", forms.EmailRX)
		form.MinLength("password", 10)
		if !form.Valid() {
			w.WriteHeader(400)
			app.render(w, r, files, "login.page.tmpl", &templateData{
				Form: form,
			})
			return
		}
		id, name, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				w.WriteHeader(400)
				form.Errors.Add("generic", "Email or Password is incorrect")
				app.render(w, r, files, "login.page.tmpl", &templateData{
					Form: form,
				})
				return
			} else {
				app.ErrorPage(w, 500)
			}
			return
		}

		err = app.session.DeleteSession(id, "")
		if err != nil {
			app.ErrorPage(w, 500)
			return
		}

		x := app.session.CreateToken()
		err = app.session.Insert(id, name, x)
		if err != nil {
			app.ErrorPage(w, 500)
			return
		}

		cookie := &http.Cookie{
			Name:   "logged",
			Value:  x,
			MaxAge: 2 * int(time.Hour),
			Path:   "/",
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/post/create", http.StatusSeeOther)
	default:
		w.Header().Set("Allow", http.MethodPost)
		app.ErrorPage(w, 405)
		return

	}
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("logged")
	if err != nil {
		app.ErrorPage(w, 500)
		return
	}
	if r.Method != "POST" {
		app.ErrorPage(w, 405)
		return
	}
	err = app.session.DeleteSession(0, c.Value)
	if err != nil {
		app.ErrorPage(w, 500)
		return
	}
	cookie := &http.Cookie{
		Name:   "logged",
		Value:  "0",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
