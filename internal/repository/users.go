package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"html"

	"git.01.alem.school/ggrks/forum.git/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	name = html.EscapeString(name)
	email = html.EscapeString(email)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
VALUES(?, ?, ?, datetime('now'))`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		return models.ErrDuplicateEmail
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, string, error) {
	var id int
	var name string
	var hashedPassword []byte
	stmt := "SELECT id, name, hashed_password FROM users WHERE email = ? AND active = TRUE"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &name, &hashedPassword)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", models.ErrInvalidCredentials
		} else {
			return 0, "", err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, "", models.ErrInvalidCredentials
		} else {
			return 0, "", err
		}
	}
	// Otherwise, the password is correct. Return the user ID.
	return id, name, nil
}

func (m *UserModel) Get(id int) (string, error) {
	var name string
	stmt := `SELECT name FROM users where id=?`
	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}
