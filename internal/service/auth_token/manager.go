package authtoken

import (
	"context"
	"errors"
	"time"
	"wschat/internal/repository"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenManager struct {
	AccessSecret  string
	RefreshSecret string
	AccessExp     int
	RefreshExp    int
	Redis         *repository.Redis
}

func NewManager(accSecret, refSecret string, accExp, refExp int, rd *repository.Redis) *TokenManager {
	return &TokenManager{
		AccessSecret:  accSecret,
		RefreshSecret: refSecret,
		AccessExp:     accExp,
		RefreshExp:    refExp,
		Redis:         rd,
	}
}

func (tm *TokenManager) GenerateAccess(id int64) (string, error) {
	generateAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Minute * time.Duration(tm.AccessExp)).Unix(),
	})

	accessToken, err := generateAccessToken.SignedString([]byte(tm.AccessSecret))

	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return accessToken, nil
}

func (tm *TokenManager) GenerateRefresh(ctx context.Context, id int64) (string, error) {

	refreshToken := uuid.New().String()

	if done := tm.Redis.InsertToken(ctx, refreshToken, id, tm.RefreshExp); !done {
		msg := "refresh token can't be inserted to redis"
		return "", errors.New(msg)
	}

	return refreshToken, nil
}

func (tm *TokenManager) ParseToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr,
		func(t *jwt.Token) (any, error) {
			return []byte(tm.AccessSecret), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
	)

	if err != nil {
		return token, err
	}

	return token, nil
}

func (tm *TokenManager) GetRefreshToken(ctx context.Context, refreshStr string) (int64, error) {
	id, err := tm.Redis.GetToken(ctx, refreshStr)

	if err != nil {
		return 0, err
	}

	return id, nil
}
