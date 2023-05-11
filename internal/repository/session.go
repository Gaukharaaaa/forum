package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type SessionModel struct {
	DB *sql.DB
}

func (m *SessionModel) Insert(userid int, username, token string) error {
	stmt := `INSERT INTO session (user_id, user_name, token) VALUES(?, ?,?);`
	_, err := m.DB.Exec(stmt, userid, username, token)
	fmt.Println(err)
	if err != nil {
		return err
	}
	return nil
}

func (m *SessionModel) CreateToken() string {
	u2, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}
	return u2.String()
}

func (m *SessionModel) DeleteSession(id int, token string) error {
	if token == "" {
		stmt := `DELETE FROM session WHERE user_id=?;`
		_, err := m.DB.Exec(stmt, id)
		if err != nil {
			return err
		}
	}
	if id == 0 {
		stmt := `DELETE FROM session WHERE token=?;`
		_, err := m.DB.Exec(stmt, token)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *SessionModel) IsSessionExist(token string) bool {
	stmt := `SELECT user_id FROM session where token=?`
	row := m.DB.QueryRow(stmt, token)
	var id int
	row.Scan(&id)
	if id < 1 {
		return false
	}
	return true
}

func (m *SessionModel) GetUser(token string) (int, string, error) {
	stmt := `SELECT user_id, user_name FROM session where token=?`
	row := m.DB.QueryRow(stmt, token)
	var id int
	var name string
	err := row.Scan(&id, &name)
	if err != nil {
		return 0, "", err
	}
	return id, name, nil
}
