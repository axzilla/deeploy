package models

import "time"

type UserDB struct {
	ID        string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserApp struct {
	ID    string
	Email string
}

func (u *UserDB) ToUserApp() *UserApp {
	return &UserApp{
		ID:    u.ID,
		Email: u.Email,
	}
}
