package api

import (
	"go.uber.org/fx"
	"net/http"
	"repertoire/auth/api/handler"
	"repertoire/auth/api/middleware"
	"repertoire/auth/api/router"
	"repertoire/auth/api/routes"
	"repertoire/auth/api/server"
	"repertoire/auth/api/validation"
)

var middlewares = fx.Options(
	fx.Provide(middleware.NewAuthMiddleware),
	fx.Provide(middleware.NewCorsMiddleware),
	fx.Provide(middleware.NewErrorHandlerMiddleware),
)

var handlers = fx.Options(
	fx.Provide(handler.NewCentrifugoHandler),
	fx.Provide(handler.NewMainHandler),
	fx.Provide(handler.NewStorageHandler),
)

var routers = fx.Options(
	fx.Provide(router.NewCentrifugoRouter),
	fx.Provide(router.NewMainRouter),
	fx.Provide(router.NewStorageRouter),
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
