package db

import (
	"context"

	"auth_service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct{ pool *pgxpool.Pool }

func New(ctx context.Context, dsn string) (*Storage, error) {
	pool, err := pgxpool.New(ctx, dsn)
	return &Storage{pool}, err
}

func (s *Storage) CreateUser(ctx context.Context, u *models.User) error {
	return s.pool.QueryRow(ctx,
		`INSERT INTO users(email,password_hash) VALUES($1,$2) RETURNING id,created_at`,
		u.Email, u.Password).Scan(&u.ID, &u.CreatedAt)
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var u models.User
	err := s.pool.QueryRow(ctx,
		`SELECT id,email,password_hash,created_at FROM users WHERE email=$1`, email).
		Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt)
	return u, err
}
