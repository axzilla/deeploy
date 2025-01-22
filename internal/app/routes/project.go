package routes

import (
	"github.com/axzilla/deeploy/internal/app/deeploy"
	mw "github.com/axzilla/deeploy/internal/app/middleware"
	apihandler "github.com/axzilla/deeploy/internal/cli/handler"
	"github.com/axzilla/deeploy/internal/shared/repos"
	"github.com/axzilla/deeploy/internal/shared/services"
)

func Project(app deeploy.App) {
	projectRepo := repos.NewProjectRepo(app.DB)
	projectService := services.NewProjectService(projectRepo)
	apiProjectHandler := apihandler.NewProjectHandler(projectService)

	userRepo := repos.NewUserRepo(app.DB)
	userService := services.NewUserService(userRepo)
	auth := mw.NewAuthMiddleware(userService)

	// CLI
	app.Router.HandleFunc("POST /api/projects", auth.Auth(apiProjectHandler.Create))
	app.Router.HandleFunc("GET /api/projects/{id}", auth.Auth(apiProjectHandler.Project))
	app.Router.HandleFunc("GET /api/projects", auth.Auth(apiProjectHandler.ProjectsByUser))
	app.Router.HandleFunc("PATCH /api/projects/{id}", auth.Auth(apiProjectHandler.Update))
	app.Router.HandleFunc("DELETE /api/projects/{id}", auth.Auth(apiProjectHandler.Delete))

	// Web coming soon
}
