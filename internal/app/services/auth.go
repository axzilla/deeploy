package services

import "github.com/axzilla/deeploy/internal/app/models"

type AuthServiceInterface interface{}

type AuthService struct {
	model models.AuthModelInterface
}

func NewAuthService(model *models.AuthModel) *AuthService {
	return &AuthService{model: model}
}
