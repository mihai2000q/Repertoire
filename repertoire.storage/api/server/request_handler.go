package server

import (
	"github.com/gin-gonic/gin"
	"repertoire/storage/api/middleware"
)

type RequestHandler struct {
	Gin    *gin.Engine
	Router *gin.RouterGroup
}

func NewRequestHandler(
	corsMiddleware middleware.CorsMiddleware,
	errorHandlerMiddleware middleware.ErrorHandlerMiddleware,
) *RequestHandler {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(corsMiddleware.Handler())
	engine.Use(errorHandlerMiddleware.Handler())

	router := engine.Group("/api")

	return &RequestHandler{
		Gin:    engine,
		Router: router,
	}
}
