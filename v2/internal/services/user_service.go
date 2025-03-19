package services

import (
	"binai.net/v2/internal/models"
	"binai.net/v2/internal/repository"
)

type UserService interface {
	GetUserInfo(userID int) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) GetUserInfo(userID int) (*models.User, error) {
	return s.repo.UserInfo(userID)
}
