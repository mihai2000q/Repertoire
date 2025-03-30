package router

import (
	"repertoire/auth/api/handler"
	"repertoire/auth/api/server"
)

type StorageRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.StorageHandler
}

func (c StorageRouter) RegisterRoutes() {
	api := c.requestHandler.PrivateRouter.Group("/storage")
	{
		api.POST("/token", c.handler.Token)
	}
}

func NewStorageRouter(
	requestHandler *server.RequestHandler,
	handler *handler.StorageHandler,
) StorageRouter {
	return StorageRouter{
		handler:        handler,
		requestHandler: requestHandler,
	}
}
