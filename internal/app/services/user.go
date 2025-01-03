package services

import (
	"github.com/axzilla/deeploy/internal/app/forms"
	"github.com/axzilla/deeploy/internal/app/models"
	"github.com/axzilla/deeploy/internal/app/repos"
	"github.com/google/uuid"
)

type UserServiceInterface interface {
	CreateUser(form *forms.RegisterForm) error
}

type UserService struct {
	repo repos.UserRepoInterface
}

func NewUserService(repo *repos.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(form *forms.RegisterForm) error {
	userDB := models.UserDB{
		ID:       uuid.New().String(),
		Email:    form.Email,
		Password: form.Password, // TODO: Here we hash the password hash(password)
	}
	return s.repo.CreateUser(userDB)
}
