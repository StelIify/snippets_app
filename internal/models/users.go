package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int
	Name     string
	Email    string
	Password []byte
	Created  time.Time
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) Create(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	statement := `INSERT INTO users (name, email, hashed_password, created) VALUES (
		$1,
		$2,
		$3,
		CURRENT_TIMESTAMP
	)`
	_, err = m.DB.Exec(context.Background(), statement, name, email, hashedPassword)
	if err != nil {
		var pg_err *pgconn.PgError
		if errors.As(err, &pg_err) && pg_err.Code == pgerrcode.UniqueViolation {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 1, nil
}
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
