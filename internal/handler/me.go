package handler

import (
	"net/http"
	"wschat/internal/domain"
	"wschat/internal/middleware"
	"wschat/internal/token"

	"github.com/gin-gonic/gin"
)

type MeHandler struct {
	svc domain.UserService
	tm  *token.TokenManager
}

func NewMe(s domain.UserService, manager *token.TokenManager) *MeHandler {
	return &MeHandler{
		svc: s,
		tm:  manager,
	}
}

func (m *MeHandler) MeRoutes(g *gin.Engine) {
	g.GET("/me", middleware.AuthMiddleware(m.tm), m.GetMe)
}

func (m *MeHandler) GetMe(c *gin.Context) {
	id := c.MustGet("userID")

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
