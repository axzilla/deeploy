package forms

import (
	"github.com/axzilla/deeploy/internal/app/utils"
	"github.com/charmbracelet/bubbles/textinput"
)

type RegisterErrors struct {
	Email           string
	Password        string
	PasswordConfirm string
}

type RegisterForm struct {
	Email           textinput.Model
	Password        textinput.Model
	PasswordConfirm textinput.Model
}

func (f *RegisterForm) Validate() RegisterErrors {
	var errors RegisterErrors
	if !utils.IsEmailValid(f.Email.Value()) {
		errors.Email = "Not a valid email"
	}
	if f.Email.Value() == "" {
		errors.Email = "Email is required"
	}
	if f.Password.Value() == "" {
		errors.Password = "Password is required"
	}
	if f.PasswordConfirm.Value() == "" {
		errors.PasswordConfirm = "Confirm your password"
	}
	if f.Password.Value() != f.PasswordConfirm.Value() {
		errors.PasswordConfirm = "Passwords do not match"
	}
	return errors
}

func (e *RegisterErrors) HasErrors() bool {
	return e.Email != "" || e.Password != "" || e.PasswordConfirm != ""
}

type LoginErrors struct {
	Email    string
	Password string
}

type LoginForm struct {
	Email    textinput.Model
	Password textinput.Model
}

func (f *LoginForm) Validate() LoginErrors {
	var errors LoginErrors
	if !utils.IsEmailValid(f.Email.Value()) {
		errors.Email = "Not a valid email"
	}
	if f.Email.Value() == "" {
		errors.Email = "Email is required"
	}
	if f.Password.Value() == "" {
		errors.Password = "Password is required"
	}
	return errors
}

func (e *LoginErrors) HasErrors() bool {
	return e.Email != "" || e.Password != ""
}
