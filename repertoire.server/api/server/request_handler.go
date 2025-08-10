package server

import (
	"repertoire/server/api/middleware"
	"repertoire/server/data/logger"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	Gin           *gin.Engine
	PublicRouter  *gin.RouterGroup
	PrivateRouter *gin.RouterGroup
}

func NewRequestHandler(
	jwtAuthMiddleware middleware.JWTAuthMiddleware,
	corsMiddleware middleware.CorsMiddleware,
	errorHandlerMiddleware middleware.ErrorHandlerMiddleware,
	logger *logger.GinLogger,
) *RequestHandler {
	gin.DefaultWriter = logger
	gin.DefaultErrorWriter = logger
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(corsMiddleware.Handler())
	engine.Use(errorHandlerMiddleware.Handler())

	publicRouter := engine.Group("/api")

	var privateRouter = &gin.RouterGroup{}
	*privateRouter = *publicRouter
	privateRouter.Use(jwtAuthMiddleware.Handler())

	return &RequestHandler{
		Gin:           engine,
		PublicRouter:  publicRouter,
		PrivateRouter: privateRouter,
	}
}
