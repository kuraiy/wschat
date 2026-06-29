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
	g.POST("/login", h.Login)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var json AuthDTO

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := h.svc.CreateUser(c.Request.Context(), json.Username, json.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Login(c *gin.Context) {

}
