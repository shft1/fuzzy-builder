// Package repositories contains the repositories for the application.
package repositories

import (
	"context"
	"errors"

	"github.com/alexm/fuzzy-builder/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	row := r.pool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, role, full_name) VALUES ($1,$2,$3,$4) RETURNING id, created_at`,
		user.Email, user.PasswordHash, user.Role, user.FullName,
	)
	return row.Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, role, full_name, created_at FROM users WHERE email=$1`, email,
	)
	var u models.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.FullName, &u.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, role, full_name, created_at FROM users WHERE id=$1`, id,
	)
	var u models.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.FullName, &u.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}
