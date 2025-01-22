package routes

import (
	"fmt"
	"net/http"

	"github.com/axzilla/deeploy/internal/app/deeploy"
	"github.com/axzilla/deeploy/internal/app/handlers"
	"github.com/axzilla/deeploy/internal/shared/repos"
	"github.com/axzilla/deeploy/internal/shared/services"

	mw "github.com/axzilla/deeploy/internal/app/middleware"
)

func Base(app deeploy.App) {
	baseHandler := handlers.NewBaseHandler()
	userRepo := repos.NewUserRepo(app.DB)
	userService := services.NewUserService(userRepo)

	auth := mw.NewAuthMiddleware(userService)

	app.Router.HandleFunc("GET /dashboard", mw.RequireAuth(auth.Auth(baseHandler.DashboardView)))
	app.Router.HandleFunc("POST /dashboard", auth.Auth(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello from server!")
	}))
}
