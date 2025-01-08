package auth

import (
	"context"

	"github.com/axzilla/deeploy/internal/app/models"
)

func GetUser(ctx context.Context) *models.UserApp {
	user, ok := ctx.Value("user").(*models.UserApp)
	if ok {
		return user
	}
	return nil
}

func IsAuthenticated(ctx context.Context) bool {
	return GetUser(ctx) != nil
}
