package routes

import (
	"github.com/axzilla/deeploy/internal/app/deeploy"
	"github.com/axzilla/deeploy/internal/app/handler"

	mw "github.com/axzilla/deeploy/internal/app/middleware"
)

func Base(app deeploy.App) {
	baseHandler := handler.NewBaseHandler()

	app.Router.HandleFunc("GET /", baseHandler.LandingView)
	app.Router.HandleFunc("GET /dashboard", mw.RequireAuth(mw.AuthMiddleware(baseHandler.DashboardView)))
}
