package handler

import (
	"errors"
	"net/http"
	"time"
	"wschat/internal/domain"
	"wschat/internal/dto"
	"wschat/internal/helpers"
	auth_token "wschat/internal/service/auth_token"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
}

var ErrNotFound = errors.New("user not found")
var ErrTaken = errors.New("username is already taken")

func (m *UserHandler) GetMe(c *gin.Context) {
	id := c.MustGet("userID").(int64)

	user, err := m.svc.GetUser(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": ErrNotFound.Error(),
			})
			return
		}
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	refresh, _ := c.Cookie("refresh_token")
	m.svc.SignOut(c.Request.Context(), refresh)

	helpers.SetCookie(c, "refresh_token", "", -1, time.Hour)
	helpers.SetCookie(c, "access_token", "", -1, time.Minute)
	c.Status(http.StatusOK)
}
