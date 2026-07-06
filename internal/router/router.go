package router

import (
	"wschat/internal/handler"
	"wschat/internal/middleware"
	"wschat/internal/service/auth_token"

	"github.com/gin-gonic/gin"
)

func New(authHandler *handler.AuthHandler, userHandler *handler.UserHandler, tm *auth_token.TokenManager) *gin.Engine {
	r := gin.Default()

	handler.RegisterValidators()

	authHandler.AuthRoutes(r)

	uRoutes := r.Group("/users", middleware.AuthMiddleware(tm))
	userHandler.UserRoutes(uRoutes)

	return r
}
