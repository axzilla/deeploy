package services

import (
	"github.com/axzilla/deeploy/internal/app/auth"
	"github.com/axzilla/deeploy/internal/app/errs"
	"github.com/axzilla/deeploy/internal/app/forms"
	"github.com/axzilla/deeploy/internal/app/jwt"
	"github.com/axzilla/deeploy/internal/app/models"
	"github.com/axzilla/deeploy/internal/app/repos"
	"github.com/google/uuid"
)

type UserServiceInterface interface {
	Register(form forms.RegisterForm) (string, error)
	Login(email, password string) (string, error)
	GetUserByID(id string) (*models.UserApp, error)
	HasUser() (bool, error)
}

type UserService struct {
	repo repos.UserRepoInterface
}

func NewUserService(repo *repos.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) HasUser() (bool, error) {
	count, err := s.repo.CountUsers()
	if err != nil {
		return false, err
	}
	hasUser := count > 0
	return hasUser, nil
}

func (s *UserService) Register(form forms.RegisterForm) (string, error) {
	foundUser, err := s.repo.GetUserByEmail(form.Email)
	if err != nil {
		return "", err
	}
	if foundUser != nil {
		return "", errs.ErrDuplicateEmail
	}
	hashedPwd, err := auth.HashPassword(form.Password)
	if err != nil {
		return "", err
	}
	userDB := &models.UserDB{
		ID:       uuid.New().String(),
		Email:    form.Email,
		Password: hashedPwd,
	}
	err = s.repo.CreateUser(userDB)
	if err != nil {
		return "", err
	}
	token, err := jwt.CreateToken(userDB.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *UserService) Login(email, password string) (string, error) {
	userDB, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if userDB == nil {
		return "", errs.ErrInvalidCredentials
	}
	if !auth.ComparePassword(userDB.Password, password) {
		return "", errs.ErrInvalidCredentials
	}
	token, err := jwt.CreateToken(userDB.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *UserService) GetUserByID(id string) (*models.UserApp, error) {
	userDB, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return userDB.ToUserApp(), nil
}
