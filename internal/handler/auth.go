package handler

import (
	"net/http"
	"time"
	"wschat/internal/domain"
	"wschat/internal/dto"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc domain.UserService
}

func NewAuth(s domain.UserService) *AuthHandler {
	return &AuthHandler{
		svc: s,
	}
}

func (h *AuthHandler) AuthRoutes(g *gin.Engine) {
	g.POST("/sign-up", h.SignUp)
	g.POST("/sign-in", h.SignIn)
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var json dto.SignUpDTO

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.svc.SignUp(c.Request.Context(), json.Username, json.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.Status(http.StatusOK)
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var json dto.SignInDTO

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	creds, err := h.svc.SignIn(c.Request.Context(), json.Username, json.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	setCookie(c, "access_token", creds.AccessToken, creds.AccessExp)
	setCookie(c, "refresh_token", creds.RefreshToken, creds.RefreshExp)

	c.JSON(http.StatusOK, creds.ID)
}

func setCookie(c *gin.Context, name, value string, exp int) {
	c.SetCookie(
		name,
		value,
		int((time.Hour * time.Duration(exp)).Seconds()),
		"/",
		"",
		false,
		true,
	)
}
