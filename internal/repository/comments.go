package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"html"

	"git.01.alem.school/ggrks/forum.git/internal/models"
)

type CommentModel struct {
	DB *sql.DB
}

func (m *CommentModel) InsertComment(postid, userid int, username, content string) error {
	username = html.EscapeString(username)
	content = html.EscapeString(content)
	stmt := `INSERT INTO comments (author, userid, postid, content)
	VALUES(?, ?, ?, ?)`
	_, err := m.DB.Exec(stmt, username, userid, postid, content)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (m *CommentModel) GetComment(postid int) ([]*models.Comment, error) {
	stmt := `SELECT id, author, userid, content, postid FROM comments where postid=?`
	row, err := m.DB.Query(stmt, postid)
	comments := []*models.Comment{}
	defer row.Close()
	for row.Next() {
		s := &models.Comment{}
		err := row.Scan(&s.ID, &s.Author, &s.Userid, &s.Content, &s.Postid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, models.ErrNoRecord
			} else {
				return nil, err
			}
		}
		stmt := `SELECT SUM(like), SUM(dislike) FROM commentreaction GROUP BY commentid HAVING commentid=?`
		row := m.DB.QueryRow(stmt, s.ID)
		err = row.Scan(&s.Like, &s.Dislike)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				s = &models.Comment{
					ID:      s.ID,
					Author:  s.Author,
					Userid:  s.Userid,
					Content: s.Content,
					Postid:  s.Postid,
					Like:    s.Like,
					Dislike: s.Dislike,
				}
			} else {
				return nil, err
			}
		}

		comments = append(comments, s)
	}
	if err = row.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (m *CommentModel) GetCommentById(commentId int) error {
	stmt := `SELECT id FROM comments WHERE id = ?`
	row := m.DB.QueryRow(stmt, commentId)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return err
	}
	return nil
}
