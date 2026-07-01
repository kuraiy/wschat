package middleware

import (
	"errors"
	"net/http"
	"time"
	auth_token "wschat/internal/service/auth_token"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(tm *auth_token.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("access_token")

		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		token, err := tm.ParseToken(tokenStr)

		if err != nil {
			if !errors.Is(err, jwt.ErrTokenExpired) {
				sendUnauthorized(c)
				return
			}

			refreshCookie, err := c.Cookie("refresh_token")

			if err != nil {
				sendUnauthorized(c)
				return
			}

			refreshToken, err := tm.TryRefresh(refreshCookie)

			if err != nil || !refreshToken.Valid {
				sendUnauthorized(c)
				return
			}
			refreshClaims := refreshToken.Claims.(jwt.MapClaims)
			userID := int64(refreshClaims["id"].(float64))

			token, err := tm.GenerateAccess(userID)

			if err != nil {
				sendUnauthorized(c)
				return
			}

			c.SetCookie(
				"access_token",
				token,
				int((time.Hour * time.Duration(tm.AccessExp)).Seconds()),
				"/",
				"",
				false,
				true,
			)
			c.Set("userID", userID)
		} else {
			if !token.Valid {
				sendUnauthorized(c)
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			userID := int64(claims["id"].(float64))
			c.Set("userID", userID)
		}
		c.Next()
	}
}

func sendUnauthorized(c *gin.Context) {
	c.Status(http.StatusUnauthorized)
	c.Abort()
}
