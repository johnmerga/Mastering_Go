package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id              int
	Name            string
	Email           string
	Hashed_password []byte
	Created         time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Create(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stms := `INSERT INTO users(name,email,hashed_password,created)
  VALUES(?,?,?,UTC_TIMESTAMP())`
	_, err = m.DB.Exec(stms, name, email, string(hashedPassword))
	if err != nil {
		var mySQlErr *mysql.MySQLError
		if errors.As(err, &mySQlErr) {
			if mySQlErr.Number == 1062 && strings.Contains(mySQlErr.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) AuthenticateUser(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT id,hashed_password FROM users WHERE email=?"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}
	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	row := m.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id=?)", id)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
