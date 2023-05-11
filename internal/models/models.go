package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")

	ErrInvalidCredentials = errors.New("models: invalid credentials")

	ErrDuplicateEmail = errors.New("models: duplicate email")
)

type Posts struct {
	ID          int
	Title       string
	UserId      int
	UserName    string
	Description string
	Created_At  time.Time
}

type Category struct {
	ID   int
	Type string
}

type Relation struct {
	ID     int
	PostID int
	CatID  int
}

type Union struct {
	Post       Posts
	Categories []Category
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}

type Session struct {
	ID       int
	UserID   int
	UserName string
	Token    string
}

type Comment struct {
	ID      int
	Author  string
	Userid  int
	Postid  int
	Content string
	Like    int
	Dislike int
}

type Reaction struct {
	ID      int
	UserID  int
	PostID  int
	Like    int
	Dislike int
}

type Comment_Reaction struct {
	ID        int
	UserID    int
	CommentID int
	Like      int
	Dislike   int
}
