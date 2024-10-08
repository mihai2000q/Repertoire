package router

import (
	"repertoire/api/handler"
	"repertoire/api/server"
)

type UserRouter struct {
	requestHandler *server.RequestHandler
	handler        *handler.UserHandler
}

func (u UserRouter) RegisterRoutes() {
	api := u.requestHandler.PrivateRouter.Group("/users")
	{
		api.GET("/current", u.handler.GetCurrentUser)
		api.GET("/", u.handler.GetUserByEmail)
	}
}

func NewUserRouter(
	requestHandler *server.RequestHandler,
	handler *handler.UserHandler,
) UserRouter {
	return UserRouter{
		handler:        handler,
		requestHandler: requestHandler,
	}
}
