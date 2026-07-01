package domain

import "context"

type User struct {
	ID           int64
	Username     string
	PasswordHash string
}

type UserRepository interface {
	CreateUser(ctx context.Context, username string, password string) error
}

type UserService interface {
	CreateUser(ctx context.Context, username string, password string) error
}
