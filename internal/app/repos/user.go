package repos

import (
	"database/sql"

	"github.com/axzilla/deeploy/internal/app/models"
)

type UserRepoInterface interface {
	CountUsers() (int, error)
	CreateUser(user *models.UserDB) error
	GetUserByEmail(email string) (*models.UserDB, error)
	GetUserByID(id string) (*models.UserDB, error)
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CountUsers() (int, error) {
	var count int

	query := `
		SELECT count(*)
		from users`

	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
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
	if err == sql.ErrNoRows {
		return nil, nil // INFO: Like user not found
	}
	if err != nil {
		return nil, err // INFO: real db errors
	}
	return user, nil
}

func (r *UserRepo) GetUserByID(id string) (*models.UserDB, error) {
	user := &models.UserDB{}

	query := `
		SELECT id, email, password, created_at, updated_at 
		FROM users
		WHERE id = ?`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // INFO: Like user not found
	}
	if err != nil {
		return nil, err // INFO: real db errors
	}
	return user, nil
}
