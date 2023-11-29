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
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}