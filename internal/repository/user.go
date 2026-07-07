package repository

import (
	"context"
	"errors"
	"fmt"
	database "wschat/gen"
	"wschat/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pgUniqueViolation = "23505"

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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return domain.ErrUsernameTaken
		}
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *AuthRepository) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	res, err := r.queries.GetUser(ctx, username)

	if err != nil {
		return domain.User{}, fmt.Errorf("get user by username %w", err)
	}

	return domain.User{
		ID:           res.ID,
		Username:     res.Username,
		PasswordHash: res.PasswordHash,
	}, nil
}

func (r *AuthRepository) GetByID(ctx context.Context, id int64) (domain.User, error) {
	res, err := r.queries.GetById(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("get user by id: %w", err)
	}

	return domain.User{
		ID:           res.ID,
		Username:     res.Username,
		PasswordHash: res.PasswordHash,
	}, nil
}

func (r *AuthRepository) ChangeUsername(ctx context.Context, id int64, username string) error {
	_, err := r.queries.UpdateUsername(ctx, database.UpdateUsernameParams{
		ID:       id,
		Username: username,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return domain.ErrUsernameTaken
		}
		return fmt.Errorf("update username: %w", err)
	}

	return nil
}

func (r *AuthRepository) DeleteUser(ctx context.Context, id int64) error {
	_, err := r.queries.DeleteUser(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}

func (r *AuthRepository) ChangePassword(ctx context.Context, id int64, hashedPassword string) error {
	_, err := r.queries.UpdatePassword(ctx, database.UpdatePasswordParams{
		ID:           id,
		PasswordHash: hashedPassword,
	})

	if err != nil {
		return fmt.Errorf("change password: %w", err)
	}

	return nil
}
