package repository

import (
	models2 "binai.net/v2/internal/models"
	"database/sql"
	"errors"
)

type UserRepository interface {
	UserInfo(id int) (*models2.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) UserInfo(id int) (*models2.User, error) {
	stmt := `SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRow(stmt, id)

	var user models2.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models2.ErrNoRecord
		}
		return nil, err
	}

	return &user, nil
}
