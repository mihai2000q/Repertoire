package router

import (
	"repertoire/auth/api/handler"
	"repertoire/auth/api/server"
)

type CentrifugoRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.CentrifugoHandler
}

func (c CentrifugoRouter) RegisterRoutes() {
	api := c.requestHandler.PrivateRouter.Group("/centrifugo")
	{
		api.GET("/token", c.handler.Token)
	}
}

func NewCentrifugoRouter(
	requestHandler *server.RequestHandler,
	handler *handler.CentrifugoHandler,
) CentrifugoRouter {
	return CentrifugoRouter{
		handler:        handler,
		requestHandler: requestHandler,
	}
}
