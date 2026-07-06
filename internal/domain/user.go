package domain

import (
	"context"
	"wschat/internal/dto"
)

type User struct {
	ID           int64
	Username     string
	PasswordHash string
}

type UserRepository interface {
	CreateUser(ctx context.Context, username string, password string) error
	GetUserByUsername(ctx context.Context, username string) (User, error)
}

type UserService interface {
	SignUp(ctx context.Context, username string, password string) error
	SignIn(ctx context.Context, username string, password string) (dto.LoginOutput, error)
	SignOut(ctx context.Context, refresh string)
}
