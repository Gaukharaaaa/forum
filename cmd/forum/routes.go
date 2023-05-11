package main

import (
	"net/http"
)

func (app *application) router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/post", app.showPost)
	mux.HandleFunc("/myposts/", app.requireAuthentication(app.myposts))
	mux.HandleFunc("/post/create/", app.requireAuthentication(app.createPost))
	mux.HandleFunc("/post/comment", app.requireAuthentication(app.leaveComment))
	mux.HandleFunc("/post/reaction", app.requireAuthentication(app.postReaction))
	mux.HandleFunc("/post/commentreaction", app.requireAuthentication(app.commentReaction))
	mux.HandleFunc("/user/signup", app.signupUserForm)
	mux.HandleFunc("/category", app.categoryfilter)
	mux.HandleFunc("/user/login", app.loginUserForm)
	mux.HandleFunc("/user/logout", app.requireAuthentication(app.logoutUser))
	mux.HandleFunc("/likedposts/", app.requireAuthentication(app.likedPosts))
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
