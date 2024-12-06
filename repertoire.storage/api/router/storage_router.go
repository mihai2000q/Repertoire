package router

import (
	"repertoire/storage/api/handler"
	"repertoire/storage/api/server"
)

type StorageRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.StorageHandler
}

func (s StorageRouter) RegisterRoutes() {
	api := s.requestHandler.PublicRouter.Group("/storage")
	{
		api.GET("/files/*filePath", s.handler.Get)
	}

	privateApi := s.requestHandler.PrivateRouter.Group("/storage")
	{
		privateApi.PUT("/upload", s.handler.Upload)
		privateApi.DELETE("/files/*filePath", s.handler.DeleteFile)
		privateApi.DELETE("/directories/*directoryPath", s.handler.DeleteDirectory)
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
