package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("access_token")

		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		valid, exp, err := checkToken(tokenStr, secret)

		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		if exp {
			refresh, err := c.Cookie("refresh_token")
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "token expired",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := int64(claims["id"].(float64))
		c.Set("userID", userID)
		c.Next()
	}
}

func checkToken(tokenStr string, secret string) (valid bool, expired bool, err error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return false, true, nil
		}
		return false, false, errors.New("invalid token")
	}

	if !token.Valid {
		return false, false, errors.New("invalid token")
	}

	return true, false, nil
}
