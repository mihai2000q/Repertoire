package router

import (
	"repertoire/server/api/handler"
	"repertoire/server/api/server"
)

type SearchRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.SearchHandler
}

func (s SearchRouter) RegisterRoutes() {
	api := s.requestHandler.PrivateRouter.Group("/search")
	{
		api.GET("", s.handler.Get)
	}
}

func NewSearchRouter(
	requestHandler *server.RequestHandler,
	handler *handler.SearchHandler,
) SearchRouter {
	return SearchRouter{
		handler:        handler,
		requestHandler: requestHandler,
	}
}
