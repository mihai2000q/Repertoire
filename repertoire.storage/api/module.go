package api

import (
	"go.uber.org/fx"
	"net/http"
	"repertoire/storage/api/handler"
	"repertoire/storage/api/middleware"
	"repertoire/storage/api/router"
	"repertoire/storage/api/routes"
	"repertoire/storage/api/server"
)

var Module = fx.Options(
	fx.Provide(middleware.NewCorsMiddleware),
	fx.Provide(middleware.NewErrorHandlerMiddleware),
	fx.Provide(server.NewRequestHandler),
	fx.Provide(handler.NewStorageHandler),
	fx.Provide(router.NewStorageRouter),
	fx.Provide(routes.NewRoutes),
	fx.Provide(server.NewServer),
	fx.Invoke(func(*routes.Routes) {}),
	fx.Invoke(func(*http.Server) {}),
)
