package repos

import (
	"database/sql"

	"github.com/axzilla/deeploy/internal/app/models"
)

type UserRepoInterface interface {
	CreateUser(user *models.UserDB) error
	GetUserByEmail(email string) (*models.UserDB, error)
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (m *UserRepo) CreateUser(user *models.UserDB) error {
	query := `
		INSERT INTO users (id, email, password)
		VALUES(?, ?, ?)`

	_, err := m.db.Exec(query, user.ID, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) GetUserByEmail(email string) (*models.UserDB, error) {
	user := &models.UserDB{}

	query := `
		SELECT id, email, password, created_at, updated_at 
		FROM users
		WHERE email = ?`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
