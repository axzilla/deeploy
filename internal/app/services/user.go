package services

import (
	"github.com/axzilla/deeploy/internal/app/auth"
	"github.com/axzilla/deeploy/internal/app/forms"
	"github.com/axzilla/deeploy/internal/app/models"
	"github.com/axzilla/deeploy/internal/app/repos"
	"github.com/google/uuid"
)

type UserServiceInterface interface {
	CreateUser(form forms.RegisterForm) (*models.UserDB, error)
	GetUserByEmail(email string) (*models.UserDB, error)
}

type UserService struct {
	repo repos.UserRepoInterface
}

func NewUserService(repo *repos.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(form forms.RegisterForm) (*models.UserDB, error) {
	hashedPwd, err := auth.HashPassword(form.Password)
	if err != nil {
		return nil, err
	}
	user := &models.UserDB{
		ID:       uuid.New().String(),
		Email:    form.Email,
		Password: hashedPwd,
	}
	err = s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.UserDB, error) {
	return s.repo.GetUserByEmail(email)
}
