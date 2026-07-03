package auth_token

import (
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

func NewManager(accSecret, refSecret string, accExp, refExp int, rd repository.Redis) *TokenManager {
	return &TokenManager{
		AccessSecret:  accSecret,
		RefreshSecret: refSecret,
		AccessExp:     accExp,
		RefreshExp:    refExp,
		Redis:         &rd,
	}
}

func (t *TokenManager) GenerateAccess(id int64) (string, error) {
	generateAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Minute * time.Duration(t.AccessExp)).Unix(),
	})

	accessToken, err := generateAccessToken.SignedString([]byte(t.AccessSecret))

	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return accessToken, nil
}

func (t *TokenManager) GenerateRefresh(id int64) (string, error) {

	refreshToken := uuid.New().String()

	if done := t.Redis.InsertToken(refreshToken, t.RefreshExp); !done {
		msg := "refresh token can't be inserted to redis"
		return "", errors.New(msg)
	}

	return refreshToken, nil
}

func (tm *TokenManager) ParseToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return []byte(tm.AccessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (tm *TokenManager) CheckRefreshToken(refreshStr string) bool {
	if exist := tm.Redis.CheckToken(refreshStr); !exist {
		return false
	}

	return true
}
