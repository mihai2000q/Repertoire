package api

import (
	"net/http"
	"repertoire/server/api/handler"
	"repertoire/server/api/middleware"
	"repertoire/server/api/router"
	"repertoire/server/api/routes"
	"repertoire/server/api/server"
	"repertoire/server/api/validation"

	"go.uber.org/fx"
)

var middlewares = fx.Options(
	fx.Provide(middleware.NewCorsMiddleware),
	fx.Provide(middleware.NewErrorHandlerMiddleware),
	fx.Provide(middleware.NewJWTAuthMiddleware),
	fx.Provide(middleware.NewMeiliAuthMiddleware),
)

var handlers = fx.Options(
	fx.Provide(handler.NewAlbumHandler),
	fx.Provide(handler.NewArtistHandler),
	fx.Provide(handler.NewPlaylistHandler),
	fx.Provide(handler.NewSearchHandler),
	fx.Provide(handler.NewSongHandler),
	fx.Provide(handler.NewUserDataHandler),
	fx.Provide(handler.NewUserHandler),
)

var routers = fx.Options(
	fx.Provide(router.NewAlbumRouter),
	fx.Provide(router.NewArtistRouter),
	fx.Provide(router.NewPlaylistRouter),
	fx.Provide(router.NewSearchRouter),
	fx.Provide(router.NewSongRouter),
	fx.Provide(router.NewSongSectionRouter),
	fx.Provide(router.NewUserDataRouter),
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
