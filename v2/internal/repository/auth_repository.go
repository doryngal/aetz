package repository

import (
	"binai.net/v2/internal/models"
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

type AuthRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	IsEmailRegistered(email string) (bool, error)
	UpdateResetCode(ctx context.Context, userID int64, code string, expiresAt, lastSentAt time.Time) error
	ClearResetCode(ctx context.Context, userID int64) error
}

type PqAuthRepository struct {
	DB *sql.DB
	//Redis *redis.Client
	//Ctx   context.Context
}

func NewPgUserRepository(db *sql.DB) *PqAuthRepository {
	return &PqAuthRepository{DB: db}
}

func (r *PqAuthRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (id, name, email, phone, password_hash, role, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.DB.Exec(query, user.ID, user.Name, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *PqAuthRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, first_name, last_name, email, phone, password_hash, role, reset_code_expires_at, reset_code, created_at, updated_at, company_id 
	          FROM users WHERE email = $1`
	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.Name,
		&user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PqAuthRepository) IsEmailRegistered(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.DB.QueryRow(query, email).Scan(&exists)

	return exists, err
}

func (r *PqAuthRepository) IsPhoneRegistered(phone string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)`
	err := r.DB.QueryRow(query, phone).Scan(&exists)

	return exists, err
}

func (r *PqAuthRepository) UpdateResetCode(ctx context.Context, userID int64, code string, expiresAt, lastSentAt time.Time) error {
	query := `
		UPDATE users
		SET reset_code = $1, reset_code_expires_at = $2, last_reset_sent_at = $3
		WHERE id = $4
	`
	_, err := r.DB.ExecContext(ctx, query, code, expiresAt, lastSentAt, userID)
	return err
}

func (r *PqAuthRepository) ClearResetCode(ctx context.Context, userID int64) error {
	query := `
		UPDATE users
		SET reset_code = NULL, reset_code_expires_at = NULL
		WHERE id = $1
	`
	_, err := r.DB.ExecContext(ctx, query, userID)
	return err
}
