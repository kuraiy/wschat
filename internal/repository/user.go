package repository

import (
	"context"
	database "wschat/gen"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	queries *database.Queries
}

func New(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		queries: database.New(db),
	}
}

func (r *AuthRepository) CreateUser(ctx context.Context, username string, password string) error {
	_, err := r.queries.CreateUser(ctx, database.CreateUserParams{
		Username:     username,
		PasswordHash: password,
	})

	if err != nil {
		return err
	}

	return nil
}
