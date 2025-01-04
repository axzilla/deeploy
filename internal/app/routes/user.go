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

	app.Router.HandleFunc("GET /login", mw.RequireGuest(userHandler.LoginView))
	app.Router.HandleFunc("POST /login", userHandler.Login)
	app.Router.HandleFunc("GET /register", mw.RequireGuest(userHandler.RegisterView))
	app.Router.HandleFunc("POST /register", userHandler.Register)
}
