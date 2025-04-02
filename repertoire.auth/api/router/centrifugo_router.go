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
	publicApi := c.requestHandler.PublicRouter.Group("/centrifugo")
	{
		publicApi.POST("/public-token", c.handler.PublicToken)
	}

	privateApi := c.requestHandler.PrivateRouter.Group("/centrifugo")
	{
		privateApi.GET("/token", c.handler.Token)
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
