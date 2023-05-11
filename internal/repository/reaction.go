package repository

import (
	"database/sql"
	"errors"
	"log"

	"git.01.alem.school/ggrks/forum.git/internal/models"
)

type ReactionModel struct {
	DB *sql.DB
}

func (m *ReactionModel) CheckExist(userid, postid int) bool {
	stmt := `SELECT like, dislike FROM reaction WHERE postid = ? and userid=?`
	row := m.DB.QueryRow(stmt, postid, userid)
	s := &models.Reaction{}

	err := row.Scan(&s.Like, &s.Dislike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		} else {
			log.Println(err)
			return false
		}
	}
	return true
}

func (m *ReactionModel) InsertReaction(userid, postid, like, dislike int) error {
	if !m.CheckExist(userid, postid) {
		stmt := `INSERT INTO reaction (userid, postid, like, dislike) VALUES(?, ?, ?, ?);`
		_, err := m.DB.Exec(stmt, userid, postid, like, dislike)
		if err != nil {
			return err
		}
		return nil
	}

	if like == 1 && dislike == 0 {
		stmt := `update reaction set like = CASE WHEN like=1 then 0 else 1 end, dislike = 0  where postid=? and userid=?;`
		_, err := m.DB.Exec(stmt, postid, userid)
		if err != nil {
			return err
		}
	} else if dislike == 1 && like == 0 {
		stmt := `update reaction set dislike = CASE WHEN dislike=1 then 0 else 1 end, like = 0  where postid=? and userid=?;`
		_, err := m.DB.Exec(stmt, postid, userid)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *ReactionModel) GetReaction(postid int) (*models.Reaction, error) {
	stmt := `SELECT SUM(like), SUM(dislike) FROM reaction GROUP BY postid HAVING postid=?;`
	row := m.DB.QueryRow(stmt, postid)
	s := &models.Reaction{}

	err := row.Scan(&s.Like, &s.Dislike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}
