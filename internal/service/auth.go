package service

import (
	"context"
	"errors"
	"time"
	"wschat/internal/domain"

	"github.com/golang-jwt/jwt/v4"
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

func (s *AuthService) CreateUser(ctx context.Context, username string, password string) (domain.AuthOutput, error) {

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return domain.AuthOutput{
			Username: username,
		}, err
	}

	newUser, err := s.repo.CreateUser(ctx, username, string(passwordHash))

	if err != nil {
		return domain.AuthOutput{
			Username: username,
		}, err
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  newUser.ID,
		"exp": time.Now().Add(time.Hour * time.Duration(s.accessExp)).Unix(),
	})

	token, err := generateToken.SignedString([]byte(s.secret))

	if err != nil {
		return domain.AuthOutput{}, errors.New("failed to generate token")
	}

	return domain.AuthOutput{
		ID:           newUser.ID,
		Username:     newUser.Username,
		AccessToken:  token,
		RefreshToken: "test",
	}, nil
}
