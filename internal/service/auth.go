package service

import (
	"context"
	"errors"
	"wschat/internal/domain"
	"wschat/internal/dto"
	auth_token "wschat/internal/service/auth_token"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo domain.UserRepository
	tm   *auth_token.TokenManager
}

func New(repo domain.UserRepository, tm *auth_token.TokenManager) *AuthService {
	return &AuthService{
		repo: repo,
		tm:   tm,
	}
}

func (s *AuthService) SignUp(ctx context.Context, username string, password string) error {
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

func (s *AuthService) SignIn(ctx context.Context, username string, password string) (dto.LoginOutput, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)

	if err != nil {
		return dto.LoginOutput{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		return dto.LoginOutput{}, errors.New("incorrect username or password")
	}

	accessToken, refreshToken, err := s.GenerateTokens(user.ID)

	if err != nil {
		return dto.LoginOutput{}, err
	}

	return dto.LoginOutput{
		ID:           user.ID,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AccessExp:    s.tm.AccessExp,
		RefreshExp:   s.tm.RefreshExp,
	}, nil
}

func (s *AuthService) SignOut(ctx context.Context, refresh string) {
	s.tm.Redis.DeleteToken(refresh)
}

func (s *AuthService) GenerateTokens(id int64) (string, string, error) {
	accessToken, err := s.tm.GenerateAccess(id)

	if err != nil {
		return "", "", errors.New("failed to generate token")
	}

	refreshToken, err := s.tm.GenerateRefresh(id)

	if err != nil {
		return "", "", errors.New("failed to refresh generate token")
	}

	return accessToken, refreshToken, nil
}
