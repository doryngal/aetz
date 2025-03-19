package repository

import (
	"database/sql"
	"log"
)

type InitRepository struct {
	UserRepo UserRepository
	LotRepo  LotRepository
	AuthRepo AuthRepository
}

func NewRepository(db *sql.DB) *InitRepository {
	return &InitRepository{
		UserRepo: NewUserRepository(db),
		LotRepo:  NewLotRepository(db),
		AuthRepo: NewPgUserRepository(db),
	}
}

func InitRepositories(db *sql.DB) *InitRepository {
	repos := NewRepository(db)
	log.Println("Repositories have been initialized")
	return repos
}
