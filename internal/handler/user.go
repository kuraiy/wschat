package handler

import (
	"net/http"
	"wschat/internal/domain"
	"wschat/internal/dto"
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
}

func (m *UserHandler) GetMe(c *gin.Context) {
	id := c.MustGet("userID").(int64)

	user, err := m.svc.GetUser(c.Request.Context(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}
