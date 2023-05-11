package repository

import (
	"database/sql"
	"errors"
	"log"

	"git.01.alem.school/ggrks/forum.git/internal/models"
)

type CommentReactionModel struct {
	DB *sql.DB
}

func (m *CommentReactionModel) CheckExist(userid, commentid int) bool {
	stmt := `SELECT like, dislike FROM commentreaction WHERE commentid = ? and userid=?`
	row := m.DB.QueryRow(stmt, commentid, userid)
	s := &models.Comment_Reaction{}

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

func (m *CommentReactionModel) InsertReaction(userid, commentid, like, dislike int) error {
	if !m.CheckExist(userid, commentid) {
		stmt := `INSERT INTO commentreaction (userid, commentid, like, dislike) VALUES(?, ?, ?, ?);`
		_, err := m.DB.Exec(stmt, userid, commentid, like, dislike)
		if err != nil {
			return err
		}
		return nil
	}

	if like == 1 && dislike == 0 {
		stmt := `update commentreaction set like = CASE WHEN like=1 then 0 else 1 end, dislike = 0  where commentid=? and userid=?;`
		_, err := m.DB.Exec(stmt, commentid, userid)
		if err != nil {
			return err
		}
	} else if dislike == 1 && like == 0 {
		stmt := `update commentreaction set dislike = CASE WHEN dislike=1 then 0 else 1 end, like = 0  where commentid=? and userid=?;`
		_, err := m.DB.Exec(stmt, commentid, userid)
		if err != nil {
			return err
		}
	}

	return nil
}
