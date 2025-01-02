package routes

import (
	"github.com/axzilla/deeploy/internal/app/deeploy"
	"github.com/axzilla/deeploy/internal/app/handler"
)

func Base(app deeploy.App) {
	baseHandler := handler.NewBaseHandler()

	app.Router.HandleFunc("GET /", baseHandler.GetLanding)
}
