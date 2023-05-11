package main

import (
	"text/template"
	"time"

	"git.01.alem.school/ggrks/forum.git/internal/forms"
	"git.01.alem.school/ggrks/forum.git/internal/models"
)

type templateData struct {
	IsAuthenticated  bool
	Form             *forms.Form
	Union            models.Union
	Posts            []models.Posts
	Comments         []*models.Comment
	Reaction         *models.Reaction
	Comment_Reaction *models.Comment_Reaction
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
