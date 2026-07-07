package domain

import (
	"context"
	"errors"
	"wschat/internal/dto"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

type User struct {
	ID           int64
	Username     string
	PasswordHash string
}

type UserRepository interface {
	CreateUser(ctx context.Context, username string, password string) error
	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetByID(ctx context.Context, id int64) (User, error)
	ChangeUsername(ctx context.Context, id int64, newUsername string) error
	ChangePassword(ctx context.Context, id int64, newPass string) error
	DeleteUser(ctx context.Context, id int64) error
}

type UserService interface {
	SignUp(ctx context.Context, username string, password string) error
	SignIn(ctx context.Context, username string, password string) (dto.LoginOutput, error)
	SignOut(ctx context.Context, refresh string)
	ChangeUsername(ctx context.Context, id int64, newUsername string) error
	GetUser(ctx context.Context, id int64) (dto.GetMeDTO, error)
	ChangePassword(ctx context.Context, id int64, json dto.ChangePasswordDTO) error
	DeleteUser(ctx context.Context, id int64) error
}
