package server

import (
	"github.com/gin-gonic/gin"
	"repertoire/api/middleware"
)

type RequestHandler struct {
	Gin           *gin.Engine
	PublicRouter  *gin.RouterGroup
	PrivateRouter *gin.RouterGroup
}

func NewRequestHandler(jwtAuthMiddleware middleware.JWTAuthMiddleware) *RequestHandler {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	publicRouter := engine.Group("/api")

	privateRouter := engine.Group("/api")
	privateRouter.Use(jwtAuthMiddleware.Handler())

	return &RequestHandler{
		Gin:           engine,
		PublicRouter:  publicRouter,
		PrivateRouter: privateRouter,
	}
}
