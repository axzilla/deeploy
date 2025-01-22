package routes

import (
	"github.com/axzilla/deeploy/internal/deeploy"
	handlers "github.com/axzilla/deeploy/internal/handlers/web"
	"github.com/axzilla/deeploy/internal/repos"
	"github.com/axzilla/deeploy/internal/services"

	mw "github.com/axzilla/deeploy/internal/middleware"
)

func Base(app deeploy.App) {
	dashboardHandler := handlers.NewDashboardHandler()
	userRepo := repos.NewUserRepo(app.DB)
	userService := services.NewUserService(userRepo)

	auth := mw.NewAuthMiddleware(userService)

	app.Router.HandleFunc("GET /dashboard", mw.RequireAuth(auth.Auth(dashboardHandler.DashboardView)))
}
