package handler

import (
	"log"
	"net/http"

	"github.com/axzilla/deeploy/internal/app/auth"
	"github.com/axzilla/deeploy/internal/app/cookie"
	"github.com/axzilla/deeploy/internal/app/forms"
	"github.com/axzilla/deeploy/internal/app/jwt"
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

	foundUser, err := h.service.GetUserByEmail(form.Email)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		errs.General = "Something went wrong. Please try again."
		pages.Login(errs, form).Render(r.Context(), w)
		return
	}
	if foundUser == nil {
		errs.Email = "Email or password incorrect"
		errs.Password = "Email or password incorrect"
		pages.Login(errs, form).Render(r.Context(), w)
		return
	}

	matched := auth.ComparePassword(foundUser.Password, form.Password)
	if !matched {
		errs.Email = "Email or password incorrect"
		errs.Password = "Email or password incorrect"
		pages.Login(errs, form).Render(r.Context(), w)
		return
	}

	token, err := jwt.CreateToken(foundUser.ID)
	if err != nil {
		log.Printf("Failed to create token: %v", err)
		errs.General = "Something went wrong. Please try again."
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

	errs := form.Validate()
	if errs.HasErrors() {
		pages.Register(errs, form).Render(r.Context(), w)
		return
	}

	foundUser, err := h.service.GetUserByEmail(form.Email)
	if foundUser != nil {
		errs.Email = "Email address is already in use"
		pages.Register(errs, form).Render(r.Context(), w)
		return
	}

	user, err := h.service.CreateUser(form)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		errs.General = "Something went wrong. Please try again."
		pages.Register(errs, form).Render(r.Context(), w)
		return
	}

	token, err := jwt.CreateToken(user.ID)
	if err != nil {
		log.Printf("Failed to create token: %v", err)
		errs.General = "Something went wrong. Please try again."
		pages.Register(errs, form).Render(r.Context(), w)
		return
	}

	cookie.SetCookie(w, token)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
