package routes

import (
	"github.com/axzilla/deeploy/internal/app/deeploy"
	"github.com/axzilla/deeploy/internal/app/handler"
	mw "github.com/axzilla/deeploy/internal/app/middleware"
	"github.com/axzilla/deeploy/internal/app/repos"
	"github.com/axzilla/deeploy/internal/app/services"
)

func User(app deeploy.App) {
	userRepo := repos.NewUserRepo(app.DB)
	userService := services.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Views
	app.Router.HandleFunc("GET /", mw.RequireGuest(userHandler.AuthView))

	// APIs
	app.Router.HandleFunc("POST /login", userHandler.Login)
	app.Router.HandleFunc("POST /register", userHandler.Register)
	app.Router.HandleFunc("GET /logout", userHandler.Logout)
}
