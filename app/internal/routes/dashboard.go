package routes

import (
	"fmt"
	"net/http"

	"github.com/axzilla/deeploy/internal/data"
	"github.com/axzilla/deeploy/internal/deeploy"
	"github.com/axzilla/deeploy/internal/services"

	mw "github.com/axzilla/deeploy/internal/middleware"
)

func Dashboard(app deeploy.App) {
	userRepo := data.NewUserRepo(app.DB)
	userService := services.NewUserService(userRepo)

	auth := mw.NewAuthMiddleware(userService)

	app.Router.HandleFunc("GET /api/dashboard", auth.Auth(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello from server!")
	}))
}
