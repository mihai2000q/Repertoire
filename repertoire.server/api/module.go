package api

import (
	"go.uber.org/fx"
	"net/http"
	"repertoire/api/handler"
	"repertoire/api/middleware"
	"repertoire/api/router"
	"repertoire/api/routes"
	"repertoire/api/server"
)

var middlewares = fx.Options(
	fx.Provide(middleware.NewErrorHandlerMiddleware),
	fx.Provide(middleware.NewJWTAuthMiddleware),
)

var handlers = fx.Options(
	fx.Provide(handler.NewAuthHandler),
	fx.Provide(handler.NewUserHandler),
)

var routers = fx.Options(
	fx.Provide(router.NewAuthRouter),
	fx.Provide(router.NewUserRouter),
)

var Module = fx.Options(
	middlewares,
	fx.Provide(server.NewRequestHandler),
	handlers,
	routers,
	fx.Provide(routes.NewRoutes),
	fx.Provide(server.NewServer),
	fx.Invoke(func(*routes.Routes) {}),
	fx.Invoke(func(*http.Server) {}),
)
