package models

import "database/sql"

type AuthModelInterface interface{}

type AuthModel struct {
	db *sql.DB
}

func NewAuthModel(db *sql.DB) *AuthModel {
	return &AuthModel{db: db}
}
