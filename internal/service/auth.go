package service

import (
	"context"
	"fmt"
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
		return dto.LoginOutput{}, domain.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		return dto.LoginOutput{}, domain.ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.GenerateTokens(ctx, user.ID)

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
	s.tm.Redis.DeleteToken(ctx, refresh)
}

func (s *AuthService) ChangeUsername(ctx context.Context, id int64, newUsername string) error {
	err := s.repo.ChangeUsername(ctx, id, newUsername)

	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ChangePassword(ctx context.Context, id int64, passJson dto.ChangePasswordDTO) error {
	user, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(passJson.OldPassword))

	if err != nil {
		return domain.ErrInvalidCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passJson.NewPassword), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	err = s.repo.ChangePassword(ctx, id, string(hashedPassword))

	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GetUser(ctx context.Context, id int64) (dto.GetMeDTO, error) {
	user, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return dto.GetMeDTO{}, err
	}

	return dto.GetMeDTO{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

func (s *AuthService) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *AuthService) GenerateTokens(ctx context.Context, id int64) (string, string, error) {
	accessToken, err := s.tm.GenerateAccess(id)

	if err != nil {
		return "", "", fmt.Errorf("generate access token %w", err)
	}

	refreshToken, err := s.tm.GenerateRefresh(ctx, id)

	if err != nil {
		return "", "", fmt.Errorf("generate refresh token %w", err)
	}

	return accessToken, refreshToken, nil
}
