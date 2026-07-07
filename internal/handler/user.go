package handler

import (
	"net/http"
	"time"
	"wschat/internal/domain"
	"wschat/internal/dto"
	"wschat/internal/helpers"
	auth_token "wschat/internal/service/auth_token"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc domain.UserService
	tm  *auth_token.TokenManager
}

func NewUser(s domain.UserService, manager *auth_token.TokenManager) *UserHandler {
	return &UserHandler{
		svc: s,
		tm:  manager,
	}
}

func (m *UserHandler) UserRoutes(g *gin.RouterGroup) {
	g.GET("/me", m.GetMe)
	g.PATCH("/username", m.ChangeUsername)
	g.PATCH("/password", m.ChangePassword)
	g.DELETE("/delete", m.Delete)
}

func (m *UserHandler) GetMe(c *gin.Context) {
	id := c.MustGet("userID").(int64)

	user, err := m.svc.GetUser(c.Request.Context(), id)

	if err != nil {
		helpers.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (m *UserHandler) ChangeUsername(c *gin.Context) {
	var json dto.ChangeUsernameDTO

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id := c.MustGet("userID").(int64)
	err := m.svc.ChangeUsername(c.Request.Context(), id, json.Username)

	if err != nil {
		helpers.WriteError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (m *UserHandler) ChangePassword(c *gin.Context) {
	var json dto.ChangePasswordDTO

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id := c.MustGet("userID").(int64)

	err := m.svc.ChangePassword(c.Request.Context(), id, json)

	if err != nil {
		helpers.WriteError(c, err)
		return
	}

	refresh, _ := c.Cookie("refresh_token")
	m.svc.SignOut(c.Request.Context(), refresh)

	invalidateTokens(c)
	c.Status(http.StatusOK)
}

func (m *UserHandler) Delete(c *gin.Context) {
	id := c.MustGet("userID").(int64)

	err := m.svc.DeleteUser(c.Request.Context(), id)

	if err != nil {
		helpers.WriteError(c, err)
		return
	}

	refresh, _ := c.Cookie("refresh_token")
	m.svc.SignOut(c.Request.Context(), refresh)

	invalidateTokens(c)
	c.Status(http.StatusOK)
}

func invalidateTokens(c *gin.Context) {
	helpers.SetCookie(c, "refresh_token", "", -1, time.Hour)
	helpers.SetCookie(c, "access_token", "", -1, time.Minute)
}
