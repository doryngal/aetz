package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	UserInfo(id int) (*User, error)
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

func (m *UserModel) UserInfo(id int) (*User, error) {
	var user User

	stmt := `SELECT id,name,email FROM users WHERE id = $1`
	row := m.DB.QueryRow(stmt, id)

	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err == sql.ErrNoRows {
		fmt.Println("err user", err)
		return &User{}, err
	} else if err != nil {
		return &User{}, err
	}

	return &user, nil
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES ($1, $2, $3, NOW())`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// var mySQLError *mysql.MySQLError
		// if errors.As(err, &mySQLError) {
		// 	if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
		// 		return ErrDuplicateEmail
		// 	}
		// }
		return fmt.Errorf("Insert into users error: %s", err)
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = $1"

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
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

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = $1)"

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
