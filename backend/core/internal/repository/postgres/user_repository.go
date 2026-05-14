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
		INSERT INTO users (tenant_id, email, password_hash, first_name, last_name)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, is_active, created_at, updated_at
	`
	return r.db.QueryRowxContext(ctx, query, user.TenantID, user.Email, user.PasswordHash, user.FirstName, user.LastName).
		Scan(&user.ID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByEmailAndTenant(ctx context.Context, email string, tenantID string) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE email = $1 AND tenant_id = $2`
	err := r.db.GetContext(ctx, &user, query, email, tenantID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
