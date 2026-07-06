package helpers

import (
	"time"

	"github.com/gin-gonic/gin"
)

func SetCookie(c *gin.Context, name, value string, exp int, ttl time.Duration) {

	c.SetCookie(
		name,
		value,
		int((ttl * time.Duration(exp)).Seconds()),
		"/",
		"",
		false,
		true,
	)
}
