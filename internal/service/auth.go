package service

import (
	"context"
	"wschat/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      domain.UserRepository
	secret    string
	accessExp int
}

func New(repo domain.UserRepository, secret string, exp int) *AuthService {
	return &AuthService{
		repo:      repo,
		secret:    secret,
		accessExp: exp,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, username string, password string) error {

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	err = s.repo.CreateUser(ctx, username, string(passwordHash))

	if err != nil {
		return err
	}

	return nil
}
