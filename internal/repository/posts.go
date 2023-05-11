package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"html"

	"git.01.alem.school/ggrks/forum.git/internal/models"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(u models.Union) (int, error) {
	description := html.EscapeString(u.Post.Description)
	title := html.EscapeString(u.Post.Title)
	stmt := `INSERT INTO post (title, description, userid,username, created_at)
	VALUES(?, ?,?,?, datetime('now'))`
	result, err := m.DB.Exec(stmt, title, description, u.Post.UserId, u.Post.UserName)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	stmt2 := `INSERT INTO relation (post_id,cat_id) VALUES(?,?)`
	for _, v := range u.Categories {
		_, err := m.DB.Exec(stmt2, id, v.ID)
		if err != nil {
			return 0, err
		}
	}

	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *PostModel) Get(id int) (models.Union, error) {
	var u models.Union
	stmt := `SELECT post.id, title, userid, username, description, created_at FROM post WHERE post.id = ?`
	row := m.DB.QueryRow(stmt, id)
	s := models.Posts{}

	err := row.Scan(&s.ID, &s.Title, &s.UserId, &s.UserName, &s.Description, &s.Created_At)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, models.ErrNoRecord
		} else {
			return u, err
		}
	}

	stmt2 := `SELECT category.id, type from category INNER JOIN relation on relation.cat_id = category.id WHERE relation.post_id = ?`
	row2, err := m.DB.Query(stmt2, s.ID)
	if err != nil {
		return u, err
	}

	defer row2.Close()
	var catArr []models.Category
	for row2.Next() {
		var cat models.Category
		err := row2.Scan(&cat.ID, &cat.Type)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return u, models.ErrNoRecord
			} else {
				return u, err
			}
		}
		catArr = append(catArr, cat)

	}

	u = models.Union{
		Post:       s,
		Categories: catArr,
	}

	return u, nil
}

func (m *PostModel) GetCategory(cat_id int) ([]models.Posts, error) {
	stmt := `SELECT post.id, post.title, post.userid, post.username, post.description, post.created_at FROM post
	INNER JOIN relation ON relation.post_id = post.id WHERE relation.cat_id = ?`
	posts := []models.Posts{}
	row, err := m.DB.Query(stmt, cat_id)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		s := models.Posts{}
		err := row.Scan(&s.ID, &s.Title, &s.UserId, &s.UserName, &s.Description, &s.Created_At)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, models.ErrNoRecord
			} else {
				return nil, err
			}
		}
		posts = append(posts, s)

	}

	if err = row.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *PostModel) Latest() ([]models.Posts, error) {
	stmt := `select * from post order by created_at DESC LIMIT 10`
	posts := []models.Posts{}
	row, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		s := models.Posts{}
		err := row.Scan(&s.ID, &s.Title, &s.UserId, &s.UserName, &s.Description, &s.Created_At)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, models.ErrNoRecord
			} else {
				return nil, err
			}
		}
		posts = append(posts, s)
	}
	fmt.Println(posts)
	if err = row.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *PostModel) SeeAll() ([]*models.Posts, error) {
	stmt := `select * from post order by created_at DESC`
	posts := []*models.Posts{}
	row, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		s := &models.Posts{}
		err := row.Scan(&s.ID, &s.Title, &s.Description, &s.Created_At)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, models.ErrNoRecord
			} else {
				return nil, err
			}
		}
		posts = append(posts, s)
	}
	if err = row.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *PostModel) MyPosts(userid int) ([]models.Posts, error) {
	stmt := `select * from post where userid=?`
	posts := []models.Posts{}
	row, err := m.DB.Query(stmt, userid)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		s := models.Posts{}
		err := row.Scan(&s.ID, &s.Title, &s.UserId, &s.UserName, &s.Description, &s.Created_At)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, models.ErrNoRecord
			} else {
				return nil, err
			}
		}
		posts = append(posts, s)
	}
	if err = row.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *PostModel) LikedPosts(userid int) ([]models.Posts, error) {
	stmt := `select post.id, title, post.userid,post.username, post.description, created_at from post inner join reaction on post.id=reaction.postid where reaction.userid=? and reaction.like=1`
	posts := []models.Posts{}
	row, err := m.DB.Query(stmt, userid)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		s := models.Posts{}
		err := row.Scan(&s.ID, &s.Title, &s.UserId, &s.UserName, &s.Description, &s.Created_At)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, models.ErrNoRecord
			} else {
				return nil, err
			}
		}
		posts = append(posts, s)
	}
	if err = row.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *PostModel) PostAuthor(id int) (int, error) {
	stmt := `select userid from post where id=?`
	row := m.DB.QueryRow(stmt, id)
	var userid int
	err := row.Scan(&userid)
	if err != nil {
		return 0, err
	}
	return userid, nil
}
