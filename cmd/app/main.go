package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/axzilla/deeploy/internal/app/db"
	"github.com/axzilla/deeploy/internal/app/handler"
	"github.com/axzilla/deeploy/internal/web/assets"
	"github.com/axzilla/deeploy/internal/web/config"
)

func main() {
	config.LoadConfig()

	db, err := db.Init()
	fmt.Println(db)
	if err != nil {
		fmt.Println("DB error")
	}

	mux := http.NewServeMux()

	SetupAssetsRoutes(mux)

	mux.HandleFunc("GET /", handler.GetLanding)

	mux.HandleFunc("GET /login", handler.GetLogin)
	mux.HandleFunc("POST /login", handler.Login)

	mux.HandleFunc("GET /register", handler.GetLogin)
	mux.HandleFunc("POST /register", handler.Register)

	fmt.Println("Server is running on http://localhost:8090")
	http.ListenAndServe(":8090", mux)
}

func SetupAssetsRoutes(mux *http.ServeMux) {
	var isDevelopment = os.Getenv("GO_ENV") != "production"

	assetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isDevelopment {
			w.Header().Set("Cache-Control", "no-store")
		}

		var fs http.Handler
		if isDevelopment {
			fs = http.FileServer(http.Dir("./internal/app/assets"))
		} else {
			fs = http.FileServer(http.FS(assets.Assets))
		}

		fs.ServeHTTP(w, r)
	})

	mux.Handle("GET /assets/", http.StripPrefix("/assets/", assetHandler))
}
