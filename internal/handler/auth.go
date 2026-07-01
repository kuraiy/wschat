package handler

import (
	"net/http"
	"time"
	"wschat/internal/domain"
	"wschat/internal/dto"
	"wschat/internal/middleware"
	"wschat/internal/service/auth_token"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc domain.UserService
	tm  *auth_token.TokenManager
}

func NewAuth(s domain.UserService, tm *auth_token.TokenManager) *AuthHandler {
	return &AuthHandler{
		svc: s,
		tm:  tm,
	}
}

func (h *AuthHandler) AuthRoutes(g *gin.Engine) {
	g.POST("/sign-up", h.SignUp)
	g.POST("/sign-in", h.SignIn)
	g.GET("/sign-out", middleware.AuthMiddleware(h.tm), h.SignOut)
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

func (h *AuthHandler) SignOut(c *gin.Context) {
	ref, _ := c.Cookie("refresh_token")

	h.svc.SignOut(c.Request.Context(), ref)

	setCookie(c, "access_token", "", -1)
	setCookie(c, "refresh_token", "", -1)

	c.Status(http.StatusOK)
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
