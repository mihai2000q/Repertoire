package api

import (
	"net/http"
	"repertoire/api/handler"
	"repertoire/api/middleware"
	"repertoire/api/router"
	"repertoire/api/routes"
	"repertoire/api/server"
	"repertoire/api/validation"

	"go.uber.org/fx"
)

var middlewares = fx.Options(
	fx.Provide(middleware.NewErrorHandlerMiddleware),
	fx.Provide(middleware.NewJWTAuthMiddleware),
)

var handlers = fx.Options(
	fx.Provide(handler.NewAuthHandler),
	fx.Provide(handler.NewSongHandler),
	fx.Provide(handler.NewUserHandler),
)

var routers = fx.Options(
	fx.Provide(router.NewAuthRouter),
	fx.Provide(router.NewSongRouter),
	fx.Provide(router.NewUserRouter),
)

var Module = fx.Options(
	fx.Provide(validation.NewValidator),
	middlewares,
	fx.Provide(server.NewRequestHandler),
	handlers,
	routers,
	fx.Provide(routes.NewRoutes),
	fx.Provide(server.NewServer),
	fx.Invoke(func(*validation.Validator) {}),
	fx.Invoke(func(*routes.Routes) {}),
	fx.Invoke(func(*http.Server) {}),
)
