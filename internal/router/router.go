package router

import (
	"wschat/internal/handler"

	"github.com/gin-gonic/gin"
)

func New(authHandler *handler.AuthHandler, meHandler *handler.MeHandler) *gin.Engine {
	r := gin.Default()

	handler.RegisterValidators()

	authHandler.AuthRoutes(r)
	meHandler.MeRoutes(r)

	return r
}
