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
	allowOriginsFunc := func(origin string) bool { return false }
	if m.env.Environment == utils.DevelopmentEnvironment {
		allowOriginsFunc = func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost:") ||
				strings.HasPrefix(origin, "https://localhost:")
		}
	}
	config := cors.Config{
		AllowMethods:    []string{"POST", "GET", "PUT", "DELETE"},
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowOriginFunc: allowOriginsFunc,
	}

	return cors.New(config)
}
