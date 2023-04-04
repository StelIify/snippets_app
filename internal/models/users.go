package models

import (
	"context"
	"database/sql"
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
	var id int
	var hashed_password []byte
	statement := `SELECT id, hashed_password from users WHERE email=$1`
	err := m.DB.QueryRow(context.Background(), statement, email).Scan(&id, &hashed_password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashed_password, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}
func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	statement := `SELECT EXISTS (SELECT name from users where id=$1)`
	err := m.DB.QueryRow(context.Background(), statement, id).Scan(&exists)

	return exists, err
}
