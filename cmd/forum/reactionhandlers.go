package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"git.01.alem.school/ggrks/forum.git/internal/forms"
)

func (app *application) leaveComment(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.layout.tmpl",
		"./ui/html/show.page.tmpl",
	}
	switch r.Method {
	case http.MethodGet:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			app.ErrorPage(w, 404)
			return

		}
		form := forms.New(nil)
		app.render(w, r, files, "show.page.tmpl", &templateData{
			Form: form,
		})
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			app.ErrorPage(w, 400)
			return
		}
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			app.ErrorPage(w, 404)
			return
		}
		if _, err := app.posts.Get(id); err != nil {
			app.ErrorPage(w, http.StatusBadRequest)
			return
		}
		form := forms.New(r.PostForm)
		form.Required("comment")
		if !form.Valid() {
			app.render(w, r, files, "./ui/html/show.page.tmpl", &templateData{
				Form: form,
			})
			return
		}
		c, err := r.Cookie("logged")
		if err != nil {
			app.ErrorPage(w, 404)
		}
		userid, username, err := app.session.GetUser(c.Value)

		err = app.comments.InsertComment(id, userid, username, form.Get("comment"))
		http.Redirect(w, r, fmt.Sprintf("/post?id=%d", id), http.StatusSeeOther)
	default:
		w.Header().Set("Allow", http.MethodPost)
		app.ErrorPage(w, 405)

	}
}

func (app *application) postReaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		app.ErrorPage(w, 405)
		return
	}
	err := r.ParseForm()
	if err != nil {
		app.ErrorPage(w, http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		app.ErrorPage(w, 404)
		return
	}
	if _, err := app.posts.Get(id); err != nil {
		app.ErrorPage(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	like, _ := strconv.Atoi(form.Get("like"))
	if like < 0 || like > 1 {
		app.ErrorPage(w, http.StatusBadRequest)
		return
	}
	dislike, _ := strconv.Atoi(form.Get("dislike"))
	if dislike < 0 || dislike > 1 {
		app.ErrorPage(w, http.StatusBadRequest)
		return
	}
	if like == 1 && dislike == 1 || like == 0 && dislike == 0 {
		app.ErrorPage(w, http.StatusBadRequest)
		return
	}

	c, err := r.Cookie("logged")
	if err != nil {
		app.ErrorPage(w, 500)
		return
	}
	userid, _, err := app.session.GetUser(c.Value)
	if err != nil {
		app.ErrorPage(w, 500)
		return
	}

	err = app.reaction.InsertReaction(userid, id, like, dislike)
	if err != nil {
		app.ErrorPage(w, 500)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", id), http.StatusSeeOther)
}

func (app *application) commentReaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		app.ErrorPage(w, 405)
		return
	}
	err := r.ParseForm()
	if err != nil {
		app.ErrorPage(w, 400)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.ErrorPage(w, 500)
		return
	}
	if err := app.comments.GetCommentById(id); err != nil {
		app.ErrorPage(w, http.StatusBadRequest)
	}
	if err != nil || id < 1 {
		app.ErrorPage(w, 404)
		return
	}
	form := forms.New(r.PostForm)
	info1 := strings.Split(form.Get("like"), " ")
	info2 := strings.Split(form.Get("dislike"), " ")
	like := 0
	dislike := 0
	postid := 0
	if len(info1) > 1 {
		like, err = strconv.Atoi(info1[0])
		if err != nil {
			app.ErrorPage(w, 500)
			return
		}
		postid, err = strconv.Atoi(info1[1])
		if err != nil {
			app.ErrorPage(w, 500)
			return
		}

	}
	if len(info2) > 1 {
		dislike, err = strconv.Atoi(info2[0])
		if err != nil {
			app.ErrorPage(w, 500)
			return
		}
		postid, err = strconv.Atoi(info2[1])
		if err != nil {
			app.ErrorPage(w, 500)
			return
		}

	}
	if like < 0 || like > 1 {
		app.ErrorPage(w, 405)
		return
	}
	if dislike < 0 || dislike > 1 {
		app.ErrorPage(w, 405)
		return
	}
	if like == 1 && dislike == 1 || like == 0 && dislike == 0 {
		app.ErrorPage(w, 405)
		return
	}

	c, err := r.Cookie("logged")
	userid, _, err := app.session.GetUser(c.Value)
	if err != nil {
		app.ErrorPage(w, 500)
		return
	}

	err = app.comment_reaction.InsertReaction(userid, id, like, dislike)
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postid), http.StatusSeeOther)
}
