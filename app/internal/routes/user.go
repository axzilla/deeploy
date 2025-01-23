package routes

import (
	"github.com/axzilla/deeploy/internal/deeploy"
	handlers "github.com/axzilla/deeploy/internal/handlers/web"
	mw "github.com/axzilla/deeploy/internal/middleware"
	"github.com/axzilla/deeploy/internal/repos"
	"github.com/axzilla/deeploy/internal/services"
)

func User(app deeploy.App) {
	userRepo := repos.NewUserRepo(app.DB)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Views
	app.Router.HandleFunc("GET /", mw.RequireGuest(userHandler.AuthView))

	// APIs
	app.Router.HandleFunc("POST /login", userHandler.Login)
	app.Router.HandleFunc("POST /register", userHandler.Register)
	app.Router.HandleFunc("GET /logout", userHandler.Logout)
}
