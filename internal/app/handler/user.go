package handler

import (
	"fmt"
	"net/http"

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
	}

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

	err := h.service.CreateUser(&form)
	if err != nil {
		fmt.Printf("RegisterUser Error: %s", err)
	}
}
