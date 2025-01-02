package routes

import (
	"github.com/axzilla/deeploy/internal/app/deeploy"
	"github.com/axzilla/deeploy/internal/app/handler"
	"github.com/axzilla/deeploy/internal/app/models"
	"github.com/axzilla/deeploy/internal/app/services"
)

func Auth(app deeploy.App) {
	authModel := models.NewAuthModel(app.DB)
	authService := services.NewAuthService(authModel)
	authHandler := handler.NewAuthHandler(authService)

	app.Router.HandleFunc("GET /login", authHandler.GetLogin)
	app.Router.HandleFunc("POST /login", authHandler.Login)
	app.Router.HandleFunc("GET /register", authHandler.GetLogin)
	app.Router.HandleFunc("POST /register", authHandler.Register)
}
