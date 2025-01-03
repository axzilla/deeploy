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
	ID         string
	Email      string
	CreatedAt  string
	UpdateddAt string
}
