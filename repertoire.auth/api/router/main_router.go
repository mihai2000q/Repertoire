package router

import (
	"repertoire/auth/api/handler"
	"repertoire/auth/api/server"
)

type MainRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.MainHandler
}

func (m MainRouter) RegisterRoutes() {
	api := m.requestHandler.PublicRouter.Group("")
	{
		api.PUT("/refresh", m.handler.Refresh)
		api.PUT("/sign-in", m.handler.SignIn)
	}
}

func NewMainRouter(
	requestHandler *server.RequestHandler,
	handler *handler.MainHandler,
) MainRouter {
	return MainRouter{
		handler:        handler,
		requestHandler: requestHandler,
	}
}
