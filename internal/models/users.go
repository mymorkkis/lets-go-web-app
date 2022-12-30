package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	// can be between 4 and 31, the higher the cost the harder to crack but more intensive
	// TODO Test what would be an acceptable cost
	cost := 12
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(
		`INSERT INTO users (name, email, hashed_password, created)
		 VALUES(?, ?, ?, UTC_TIMESTAMP())`,
		name,
		email,
		string(hashedPassword),
	)
	if err != nil {
		// We could also use the Exists method to check before calling DB.Exec and error if true
		// That would introduce a race condition though (although very unlikely to occur)
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var (
		id             int
		hashedPassword []byte
	)

	err := m.DB.QueryRow(
		"SELECT id, hashed_password FROM users WHERE email = ?",
		email,
	).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (exists bool, err error) {
	err = m.DB.QueryRow(
		"SELECT EXISTS(SELECT true FROM users WHERE id = ?)",
		id,
	).Scan(&exists)

	return
}
