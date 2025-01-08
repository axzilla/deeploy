package routes

import (
	"github.com/axzilla/deeploy/internal/app/deeploy"
	"github.com/axzilla/deeploy/internal/app/handler"
	"github.com/axzilla/deeploy/internal/app/middleware"
	"github.com/axzilla/deeploy/internal/app/repos"
	"github.com/axzilla/deeploy/internal/app/services"

	mw "github.com/axzilla/deeploy/internal/app/middleware"
)

func Base(app deeploy.App) {
	baseHandler := handler.NewBaseHandler()
	userRepo := repos.NewUserRepo(app.DB)
	userService := services.NewUserService(userRepo)

	auth := middleware.NewAuthMiddleware(userService)

	app.Router.HandleFunc("GET /dashboard", mw.RequireAuth(auth.Auth(baseHandler.DashboardView)))
	app.Router.HandleFunc("GET /", baseHandler.LandingView)
}
