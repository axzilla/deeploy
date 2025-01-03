package main

import (
	"fmt"
	"net/http"

	"github.com/axzilla/deeploy/internal/app/db"
	"github.com/axzilla/deeploy/internal/app/deeploy"
	"github.com/axzilla/deeploy/internal/app/routes"
	"github.com/axzilla/deeploy/internal/web/config"
)

func main() {
	config.LoadConfig()

	db, err := db.Init()
	if err != nil {
		fmt.Printf("DB Error: %s", err)
	}

	mux := http.NewServeMux()
	app := deeploy.New(mux, db)

	routes.Base(app)
	routes.Assets(app)
	routes.User(app)

	fmt.Println("Server is running on http://localhost:8090")
	http.ListenAndServe(":8090", mux)
}
