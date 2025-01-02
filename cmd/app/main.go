package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/axzilla/deeploy/internal/app/ui/pages"
	"github.com/axzilla/deeploy/internal/web/assets"
	"github.com/axzilla/deeploy/internal/web/config"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	err := map[string]string{}
	if email == "" {
		err["email"] = "Email is required"
	}
	if password == "" {
		err["password"] = "Password is required"
	}

	formData := map[string]string{
		"email": email,
	}

	pages.Login(err, formData).Render(r.Context(), w)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("passwordConfirm")

	err := map[string]string{}
	if email == "" {
		err["email"] = "Email is required"
	}
	if password == "" {
		err["password"] = "Password is required"
	}
	if passwordConfirm == "" {
		err["passwordConfirm"] = "Confirm your password"
	}
	if password != passwordConfirm {
		err["password"] = "Passwords do not match"
		err["passwordConfirm"] = "Passwords do not match"
	}

	formData := map[string]string{
		"email": email,
	}

	pages.Register(err, formData).Render(r.Context(), w)
}

func main() {
	config.LoadConfig()

	mux := http.NewServeMux()

	SetupAssetsRoutes(mux)

	mux.Handle("GET /", templ.Handler(pages.Dashboard()))

	mux.Handle("GET /login", templ.Handler(pages.Login(nil, nil)))
	mux.HandleFunc("POST /login", loginHandler)

	mux.Handle("GET /register", templ.Handler(pages.Register(nil, nil)))
	mux.HandleFunc("POST /register", registerHandler)

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
