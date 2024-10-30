package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"repertoire/storage/utils"
)

type CorsMiddleware struct {
	env utils.Env
}

func NewCorsMiddleware(env utils.Env) CorsMiddleware {
	return CorsMiddleware{
		env: env,
	}
}

func (m CorsMiddleware) Handler() gin.HandlerFunc {
	allowOrigins := []string{"https://yourdomain.com"}
	if m.env.Environment == utils.DevelopmentEnvironment {
		allowOrigins = []string{"*"}
	}
	config := cors.Config{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{"POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"*"}, // []string{"Origin", "Content-Length", "Content-Type"},
	}

	return cors.New(config)
}
