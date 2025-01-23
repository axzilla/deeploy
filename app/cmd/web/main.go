package main

import (
	"fmt"
	"net/http"

	"github.com/axzilla/deeploy/internal/config"
	"github.com/axzilla/deeploy/internal/db"
	"github.com/axzilla/deeploy/internal/deeploy"
	"github.com/axzilla/deeploy/internal/routes"
)

func main() {
	config.LoadConfig()

	db, err := db.Init()
	if err != nil {
		fmt.Printf("DB Error: %s", err)
	}

	mux := http.NewServeMux()
	app := deeploy.New(mux, db)

	routes.Assets(app)
	routes.Base(app)
	routes.User(app)
	routes.Dashboard(app)
	routes.Project(app)
	routes.Pod(app)

	fmt.Println("Server is running on http://localhost:8090")
	http.ListenAndServe(":8090", mux)
}
