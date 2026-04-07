package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lonaas1337-sudo/vinylvault/user-service/internal/config"
	"github.com/lonaas1337-sudo/vinylvault/user-service/internal/model"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(cfg config.Config) (*UserRepository, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &UserRepository{
		pool: pool,
	}, nil
}

func (r *UserRepository) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, u *model.User) (int64, error) {
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`

	var userID int64

	err := r.pool.QueryRow(ctx, query, u.Email(), u.PasswordHash()).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}
