package handler

import (
	"net/http"

	"github.com/axzilla/deeploy/internal/app/ui/pages"
)

func GetLogin(w http.ResponseWriter, r *http.Request) {
	pages.Login(nil, nil).Render(r.Context(), w)
}

func Login(w http.ResponseWriter, r *http.Request) {
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

func GetRegister(w http.ResponseWriter, r *http.Request) {
	pages.Register(nil, nil).Render(r.Context(), w)
}

func Register(w http.ResponseWriter, r *http.Request) {
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
