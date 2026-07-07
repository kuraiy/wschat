package helpers

import (
	"errors"
	"log"
	"net/http"
	"wschat/internal/domain"

	"github.com/gin-gonic/gin"
)

func WriteError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrUsernameTaken):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		log.Printf("unexpected: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}
