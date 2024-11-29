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
	api := s.requestHandler.Router.Group("/storage")
	{
		api.GET("/files/*filePath", s.handler.Get)
		api.PUT("/upload", s.handler.Upload)
		api.DELETE("/files/*filePath", s.handler.DeleteFile)
		api.DELETE("/directories/*directoryPath", s.handler.DeleteDirectory)
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
