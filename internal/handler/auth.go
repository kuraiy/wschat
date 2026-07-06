package handler

import (
	"errors"
	"net/http"
	"time"
	"wschat/internal/domain"
	"wschat/internal/dto"
	"wschat/internal/helpers"
	"wschat/internal/middleware"
	"wschat/internal/service"
	"wschat/internal/service/auth_token"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": ErrTaken.Error(),
			})
			return
		}
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
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": service.ErrInvalidCredentials.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)

	helpers.SetCookie(c, "access_token", creds.AccessToken, creds.AccessExp, time.Minute)
	helpers.SetCookie(c, "refresh_token", creds.RefreshToken, creds.RefreshExp, time.Hour)

	c.JSON(http.StatusOK, creds.ID)
}

func (h *AuthHandler) SignOut(c *gin.Context) {
	ref, _ := c.Cookie("refresh_token")

	h.svc.SignOut(c.Request.Context(), ref)

	helpers.SetCookie(c, "access_token", "", -1, time.Minute)
	helpers.SetCookie(c, "refresh_token", "", -1, time.Hour)

	c.Status(http.StatusOK)
}
