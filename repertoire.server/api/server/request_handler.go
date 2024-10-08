package server

import (
	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	Gin        *gin.Engine
	BaseRouter *gin.RouterGroup
}

func NewRequestHandler() *RequestHandler {
	engine := gin.Default()
	return &RequestHandler{
		Gin:        engine,
		BaseRouter: engine.Group("/api"),
	}
}
