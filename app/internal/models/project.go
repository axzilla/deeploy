package models

import "time"

type ProjectDB struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProjectApp struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (u *ProjectDB) ToProjectApp() *ProjectApp {
	return &ProjectApp{
		ID:          u.ID,
		UserID:      u.UserID,
		Title:       u.Title,
		Description: u.Description,
	}
}
