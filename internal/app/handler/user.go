package handler

import (
	"log"
	"net/http"

	"github.com/axzilla/deeploy/internal/app/cookie"
	"github.com/axzilla/deeploy/internal/app/errs"
	"github.com/axzilla/deeploy/internal/app/forms"
	"github.com/axzilla/deeploy/internal/app/services"
	"github.com/axzilla/deeploy/internal/app/ui/pages"
)

type UserHandler struct {
	service services.UserServiceInterface
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) LoginView(w http.ResponseWriter, r *http.Request) {
	pages.Login(forms.LoginErrors{}, forms.LoginForm{}).Render(r.Context(), w)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	form := forms.LoginForm{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	errs := form.Validate()
	if errs.HasErrors() {
		pages.Login(errs, form).Render(r.Context(), w)
		return
	}

	token, err := h.service.Login(form.Email, form.Password)
	if err != nil {
		log.Printf("Login failed: %v", err)
		errs.Email = "Email or password incorrect"
		errs.Password = "Email or password incorrect"
		pages.Login(errs, form).Render(r.Context(), w)
		return
	}

	cookie.SetCookie(w, token)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *UserHandler) RegisterView(w http.ResponseWriter, r *http.Request) {
	pages.Register(forms.RegisterErrors{}, forms.RegisterForm{}).Render(r.Context(), w)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	form := forms.RegisterForm{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		PasswordConfirm: r.FormValue("passwordConfirm"),
	}
	formeErrs := form.Validate()
	if formeErrs.HasErrors() {
		pages.Register(formeErrs, form).Render(r.Context(), w)
		return
	}

	token, err := h.service.Register(form)
	if err == errs.ErrDuplicateEmail {
		formeErrs.Email = "Email address is already in use"
		pages.Register(formeErrs, form).Render(r.Context(), w)
		return
	}
	if err != nil {
		log.Printf("User creation failed: %v", err)
		formeErrs.General = "Something went wrong. Please try again."
		pages.Register(formeErrs, form).Render(r.Context(), w)
		return
	}

	cookie.SetCookie(w, token)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie.ClearCookie(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
