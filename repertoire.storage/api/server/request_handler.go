package server

import (
	"github.com/gin-gonic/gin"
	"repertoire/storage/api/middleware"
)

type RequestHandler struct {
	Gin           *gin.Engine
	PrivateRouter *gin.RouterGroup
	PublicRouter  *gin.RouterGroup
}

func NewRequestHandler(
	corsMiddleware middleware.CorsMiddleware,
	errorHandlerMiddleware middleware.ErrorHandlerMiddleware,
	authMiddleware middleware.AuthMiddleware,
) *RequestHandler {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(corsMiddleware.Handler())
	engine.Use(errorHandlerMiddleware.Handler())

	publicRouter := engine.Group("")

	var privateRouter = &gin.RouterGroup{}
	*privateRouter = *publicRouter
	privateRouter.Use(authMiddleware.Handler())

	return &RequestHandler{
		Gin:           engine,
		PublicRouter:  publicRouter,
		PrivateRouter: publicRouter,
	}
}
