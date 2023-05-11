package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"git.01.alem.school/ggrks/forum.git/internal/forms"

	"git.01.alem.school/ggrks/forum.git/internal/models"
)

var seeall = []string{
	"./ui/html/seeall.page.tmpl",
	"./ui/html/base.layout.tmpl",
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.ErrorPage(w, 404)
		return
	}
	if r.Method != "GET" {
		app.ErrorPage(w, 405)
		return
	}
	x, err := app.posts.Latest()
	if err != nil {
		log.Println(err.Error())
		app.ErrorPage(w, 500)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	app.render(w, r, files, "home.page.tmpl", &templateData{
		Posts: x,
	})
}

func (app *application) showPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if r.Method != "GET" {
		app.ErrorPage(w, 405)
		return
	}
	if err != nil || id < 1 {
		app.ErrorPage(w, 404)
		return
	}
	s, err := app.posts.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.ErrorPage(w, 404)
		} else {
			app.ErrorPage(w, 500)
		}
		return
	}

	c, err := app.comments.GetComment(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.ErrorPage(w, 404)
		} else {
			app.ErrorPage(w, 500)
		}
		return
	}

	reaction, err := app.reaction.GetReaction(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			reaction = &models.Reaction{
				Like:    0,
				Dislike: 0,
			}
		} else {
			app.ErrorPage(w, 500)
			return
		}
	}

	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	app.render(w, r, files, "show.page.tmpl", &templateData{
		Union:    s,
		Comments: c,
		Reaction: reaction,
	})
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/create.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	switch r.Method {
	case http.MethodGet:
		app.render(w, r, files, "create.page.tmpl", &templateData{
			Form: forms.New(nil),
		})
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			app.ErrorPage(w, 500)
			return
		}
		c, err := r.Cookie("logged")
		if err != nil {
			app.ErrorPage(w, 500)
			return
		}
		userid, name, err := app.session.GetUser(c.Value)
		if err != nil {
			app.ErrorPage(w, 500)
			return
		}

		form := forms.New(r.PostForm)
		form.Required("title", "description", "category")
		form.MaxLength("title", 100)
		if !form.Valid() {
			app.render(w, r, files, "create.page.tmpl", &templateData{
				Form: form,
			})
			return
		}

		cat := r.Form["category"]
		if len(cat) == 0 {
			app.ErrorPage(w, 400)
			return
		}

		p := models.Posts{
			Title:       form.Get("title"),
			UserId:      userid,
			UserName:    name,
			Description: form.Get("description"),
		}
		categories := DefineCategories(cat)
		u := models.Union{
			Post:       p,
			Categories: categories,
		}

		id, err := app.posts.Insert(u)
		if err != nil {
			log.Println(err.Error())
			app.ErrorPage(w, 500)
		}

		http.Redirect(w, r, fmt.Sprintf("/post?id=%d", id), http.StatusSeeOther)
	default:
		w.Header().Set("Allow", http.MethodPost)
		app.ErrorPage(w, 405)
	}
}

func (app *application) myposts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/myposts/" {
		app.ErrorPage(w, 404)
		return

	}
	if r.Method != "GET" {
		app.ErrorPage(w, 405)
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
	x, err := app.posts.MyPosts(userid)
	if err != nil {
		app.ErrorPage(w, 500)
		return
	}

	app.render(w, r, seeall, "seeall.page.tmpl", &templateData{
		Posts: x,
	})
}

func (app *application) likedPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/likedposts/" {
		app.ErrorPage(w, 404)
		return

	}
	if r.Method != "GET" {
		app.ErrorPage(w, 405)
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
	x, err := app.posts.LikedPosts(userid)
	if err != nil {
		app.ErrorPage(w, 500)
		return
	}

	app.render(w, r, seeall, "seeall.page.tmpl", &templateData{
		Posts: x,
	})
}

func (app *application) categoryfilter(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		app.ErrorPage(w, 405)
		return
	}
	category := r.URL.Query().Get("c")
	cat_id, err := strconv.Atoi(category)

	s, err := app.posts.GetCategory(cat_id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.ErrorPage(w, 404)
		} else {
			app.ErrorPage(w, 500)
		}
		return
	}
	fmt.Println(s)
	app.render(w, r, seeall, "seeall.page.tmpl", &templateData{
		Posts: s,
	})
}

func (app *application) isInCategories(category string) bool {
	return category == "General topic" || category == "Life style" || category == "Food" || category == "Sport" || category == "Fassion"
}

func DefineCategories(arr []string) []models.Category {
	var c models.Category
	var res []models.Category
	for _, v := range arr {
		if v == "General topic" {
			c = models.Category{
				ID:   1,
				Type: "General topic",
			}
		}
		if v == "Food" {
			c = models.Category{
				ID:   2,
				Type: "Food",
			}
		}
		if v == "Life style" {
			c = models.Category{
				ID:   3,
				Type: "Life style",
			}
		}
		if v == "Sport" {
			c = models.Category{
				ID:   4,
				Type: "Sport",
			}
		}
		if v == "Fashion" {
			c = models.Category{
				ID:   5,
				Type: "Fashion",
			}
		}
		res = append(res, c)

	}
	return res
}
