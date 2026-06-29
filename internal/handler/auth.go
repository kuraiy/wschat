package handler

import (
	"net/http"
	"wschat/internal/domain"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc domain.UserService
}

func New(s domain.UserService) *AuthHandler {
	return &AuthHandler{
		svc: s,
	}
}

func (h *AuthHandler) AuthRoutes(g *gin.Engine) {
	g.POST("/register", h.Register)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var json AuthDTO

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.svc.CreateUser(c.Request.Context(), json.Username, json.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.Status(http.StatusOK)
}
