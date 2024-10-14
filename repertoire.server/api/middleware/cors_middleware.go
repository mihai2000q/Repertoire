package middleware

import (
	"context"
	"github.com/gin-contrib/cors"
	"go.uber.org/fx"
	"repertoire/api/server"
	"repertoire/utils"
	"strings"
)

type CorsMiddleware struct {
	handler server.RequestHandler
	env     utils.Env
}

func NewCorsMiddleware(lc fx.Lifecycle, handler server.RequestHandler, env utils.Env) *CorsMiddleware {
	middleware := &CorsMiddleware{
		handler: handler,
		env:     env,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			middleware.Setup()
			return nil
		},
	})

	return middleware
}

func (c *CorsMiddleware) Setup() {
	allowOriginsFunc := func(origin string) bool { return false }
	if c.env.Environment == utils.DevelopmentEnvironment {
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
	c.handler.Gin.Use(cors.New(config))
}
