package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"repertoire/utils"
	"strings"
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
	var allowOriginsFunc func(origin string) bool = nil
	if m.env.Environment == utils.DevelopmentEnvironment {
		allowOriginsFunc = func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost:") ||
				strings.HasPrefix(origin, "https://localhost:")
		}
	}
	config := cors.Config{
		AllowOriginFunc: allowOriginsFunc,
		AllowOrigins:    []string{"https://yourdomain.com"},
		AllowMethods:    []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"*"}, // []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
	}

	return cors.New(config)
}
