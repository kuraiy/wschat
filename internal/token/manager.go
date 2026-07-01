package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenManager struct {
	AccessSecret  string
	RefreshSecret string
	AccessExp     int
	RefreshExp    int
}

func NewManager(accSecret, refSecret string, accExp, refExp int) *TokenManager {
	return &TokenManager{
		AccessSecret:  accSecret,
		RefreshSecret: refSecret,
		AccessExp:     accExp,
		RefreshExp:    refExp,
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
	generateRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   id,
		"type": "refresh",
		"exp":  time.Now().Add(time.Hour * time.Duration(t.RefreshExp)).Unix(),
	})

	refreshToken, err := generateRefreshToken.SignedString([]byte(t.RefreshSecret))

	if err != nil {
		return "", errors.New("failed to generate token")
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

func (tm *TokenManager) TryRefresh(refreshStr string) (*jwt.Token, error) {
	refresh, err := jwt.Parse(refreshStr, func(t *jwt.Token) (any, error) {
		return []byte(tm.RefreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return refresh, nil
}
