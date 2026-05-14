package postgres

import (
	"context"

	"github.com/enterprise-erp/core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, is_active, created_at, updated_at
	`
	return r.db.QueryRowxContext(ctx, query, user.Email, user.PasswordHash, user.FirstName, user.LastName).
		Scan(&user.ID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
