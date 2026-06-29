package domain

import "context"

type User struct {
	ID           int64
	Username     string
	PasswordHash string
}

type AuthOutput struct {
	ID           int64
	Username     string
	AccessToken  string
	RefreshToken string
}

type UserRepository interface {
	CreateUser(ctx context.Context, username string, password string) (User, error)
}

type UserService interface {
	CreateUser(ctx context.Context, username string, password string) (AuthOutput, error)
}
