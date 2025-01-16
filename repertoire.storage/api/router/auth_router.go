package router

import (
	"repertoire/storage/api/handler"
	"repertoire/storage/api/server"
)

type AuthRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.AuthHandler
}

func (a AuthRouter) RegisterRoutes() {
	api := a.requestHandler.PublicRouter.Group("/oauth")
	{
		api.POST("/token", a.handler.Token)
	}
}

func NewAuthRouter(
	requestHandler *server.RequestHandler,
	handler *handler.AuthHandler,
) AuthRouter {
	return AuthRouter{
		handler:        handler,
		requestHandler: requestHandler,
	}
}
