package middleware

import (
	"errors"
	"net/http"
	"time"
	"wschat/internal/helpers"
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

			idFromRedis, err := tm.GetRefreshToken(c.Request.Context(), refreshCookie)

			if err != nil {
				sendUnauthorized(c)
				return
			}
			token, err := tm.GenerateAccess(idFromRedis)

			if err != nil {
				sendUnauthorized(c)
				return
			}

			helpers.SetCookie(c, "access_token", token, tm.AccessExp, time.Minute)
			c.Set("userID", idFromRedis)
		} else {
			if !token.Valid {
				sendUnauthorized(c)
				return
			}

			userID := getClaims(token)
			c.Set("userID", userID)
		}
		c.Next()
	}
}

func sendUnauthorized(c *gin.Context) {
	c.Status(http.StatusUnauthorized)
	c.Abort()
}

func getClaims(token *jwt.Token) int64 {
	claims := token.Claims.(jwt.MapClaims)

	return int64(claims["id"].(float64))
}
