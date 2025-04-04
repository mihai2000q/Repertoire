package router

import (
	"github.com/gin-gonic/gin"
	"repertoire/server/api/handler"
	"repertoire/server/api/middleware"
	"repertoire/server/api/server"
)

type SearchRouter struct {
	requestHandler      *server.RequestHandler
	handler             *handler.SearchHandler
	meiliAuthMiddleware middleware.MeiliAuthMiddleware
}

func (s SearchRouter) RegisterRoutes() {
	privateApi := s.requestHandler.PrivateRouter.Group("/search")
	{
		privateApi.GET("", s.handler.Get)
	}

	var webHookGroup = &gin.RouterGroup{}
	*webHookGroup = *s.requestHandler.PublicRouter
	webHookGroup.Use(s.meiliAuthMiddleware.Handler())

	publicApi := webHookGroup.Group("/search/meili-webhook")
	{
		publicApi.POST("", s.handler.MeiliWebhook)
	}
}

func NewSearchRouter(
	requestHandler *server.RequestHandler,
	handler *handler.SearchHandler,
	meiliAuthMiddleware middleware.MeiliAuthMiddleware,
) SearchRouter {
	return SearchRouter{
		handler:             handler,
		requestHandler:      requestHandler,
		meiliAuthMiddleware: meiliAuthMiddleware,
	}
}
