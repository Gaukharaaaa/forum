package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	sqlite "git.01.alem.school/ggrks/forum.git/internal/repository"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	errorLog         *log.Logger
	infoLog          *log.Logger
	posts            *sqlite.PostModel
	users            *sqlite.UserModel
	session          *sqlite.SessionModel
	comments         *sqlite.CommentModel
	reaction         *sqlite.ReactionModel
	comment_reaction *sqlite.CommentReactionModel
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	db, err := openDB("forum.db")
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	app := application{
		errorLog:         errorLog,
		infoLog:          infoLog,
		posts:            &sqlite.PostModel{DB: db},
		users:            &sqlite.UserModel{DB: db},
		session:          &sqlite.SessionModel{DB: db},
		comments:         &sqlite.CommentModel{DB: db},
		reaction:         &sqlite.ReactionModel{DB: db},
		comment_reaction: &sqlite.CommentReactionModel{DB: db},
	}
	srv := &http.Server{
		Addr:         ":8080",
		ErrorLog:     errorLog,
		Handler:      app.router(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Println("http://localhost:8080")
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
