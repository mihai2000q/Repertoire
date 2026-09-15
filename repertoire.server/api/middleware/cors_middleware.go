package middleware

import (
	"repertoire/server/internal/env"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CorsMiddleware struct {
	env env.Env
}

func NewCorsMiddleware(env env.Env) CorsMiddleware {
	return CorsMiddleware{
		env: env,
	}
}

func (m CorsMiddleware) Handler() gin.HandlerFunc {
	allowOrigins := []string{"https://yourdomain.com"}
	if m.env.Environment == env.DevelopmentEnvironment {
		allowOrigins = []string{"*"}
	}
	config := cors.Config{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"*"}, // []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
	}

	return cors.New(config)
}
