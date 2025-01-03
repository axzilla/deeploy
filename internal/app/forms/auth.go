package forms

type RegisterForm struct {
	Email           string
	Password        string
	PasswordConfirm string
}

type RegisterErrors struct {
	Email           string
	Password        string
	PasswordConfirm string
}

func (f *RegisterForm) Validate() RegisterErrors {
	var errors RegisterErrors
	if f.Email == "" {
		errors.Email = "Email is required"
	}
	if f.Password == "" {
		errors.Password = "Password is required"
	}
	if f.PasswordConfirm == "" {
		errors.PasswordConfirm = "Password confirmation is required"
	}
	if f.Password != f.PasswordConfirm {
		errors.Password = "Passwords do not match"
		errors.PasswordConfirm = "Passwords do not match"
	}
	return errors
}

func (e *RegisterErrors) HasErrors() bool {
	return e.Email != "" || e.Password != "" || e.PasswordConfirm != ""
}

type LoginForm struct {
	Email    string
	Password string
}

type LoginErrors struct {
	Email    string
	Password string
}

func (f *LoginForm) Validate() LoginErrors {
	var errors LoginErrors
	if f.Email == "" {
		errors.Email = "Email is required"
	}
	if f.Password == "" {
		errors.Password = "Password is required"
	}
	return errors
}

func (e *LoginErrors) HasErrors() bool {
	return e.Email != "" || e.Password != ""
}
