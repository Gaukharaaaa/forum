package main

import (
	"database/sql"
	"fmt"

	"git.01.alem.school/ggrks/forum.git/pkg/migrator"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		panic(err)
	}
	err = migrator.Migrate("cmd/migrator/up.sql", db)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
