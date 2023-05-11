package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func (app *application) isAuthentifacated(r *http.Request) bool {
	c, err := r.Cookie("logged")
	if err != nil {
		return false
	}

	return app.session.IsSessionExist(c.Value)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.IsAuthenticated = app.isAuthentifacated(r)
	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, files []string, name string, td *templateData) {
	ts, err := template.New(name).Funcs(functions).ParseFiles(files...)
	if err != nil {
		fmt.Println(err)
		app.ErrorPage(w, 500)
		return
	}
	err = ts.Execute(w, app.addDefaultData(td, r))
	if err != nil {
		fmt.Println(err)
		app.ErrorPage(w, 500)
		return
	}
}

func (app *application) ErrorPage(w http.ResponseWriter, code int) {
	var errors string
	switch code {
	case 400:
		errors = "Bad request"
	case 405:
		errors = "Method not allowed"
	case 500:
		errors = "Internal server error"
	case 404:
		errors = "Page not found"
	}
	w.WriteHeader(code)
	ero, err := template.ParseFiles("./ui/html/error.page.html")
	if err != nil {
		fmt.Println(err)

		http.Error(w, "500"+"\n"+"Internal Server Error", 500)
		return
	}
	ero.Execute(w, errors)
}
