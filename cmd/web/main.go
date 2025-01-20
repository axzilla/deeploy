package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/axzilla/deeploy/internal/app/install"
	"github.com/axzilla/deeploy/internal/web/assets"
	"github.com/axzilla/deeploy/internal/web/config"
	"github.com/axzilla/deeploy/internal/web/ui/pages"
)

func main() {
	config.LoadConfig()
	mux := http.NewServeMux()
	SetupAssetsRoutes(mux)
	mux.Handle("GET /", templ.Handler(pages.Landing()))
	mux.Handle("GET /install.sh", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := install.InstallScript.ReadFile("install.sh")
		if err != nil {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Disposition", "attachment; filename=install.sh")
		w.Write(file)
	}))
	fmt.Println("Server is running on http://localhost:8090")
	http.ListenAndServe(":8090", mux)
}

func SetupAssetsRoutes(mux *http.ServeMux) {
	var isDev = os.Getenv("GO_ENV") != "prod"
	assetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var fs http.Handler
		if isDev {
			w.Header().Set("Cache-Control", "no-store")
			fs = http.FileServer(http.Dir("./internal/web/assets"))
		} else {
			fs = http.FileServer(http.FS(assets.Assets))
		}
		fs.ServeHTTP(w, r)
	})
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", assetHandler))
}
