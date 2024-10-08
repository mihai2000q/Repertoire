package router

import (
	"repertoire/api/handler"
	"repertoire/api/server"
)

type UserRouter struct {
	// logger      server.Logger
	requestHandler *server.RequestHandler
	handler        *handler.UserHandler
}

func (u UserRouter) SetupRoutes() {
	// u.logger.Debug("Setting up user routes")
	api := u.requestHandler.BaseRouter.Group("/users")
	{
		api.GET("/", u.handler.GetUserByEmail)
		api.GET("/test", u.handler.Test)
	}
}

func NewUserRouter(
	//logger server.Logger,
	requestHandler *server.RequestHandler,
	handler *handler.UserHandler,
) UserRouter {
	return UserRouter{
		//logger: 		logger,
		handler:        handler,
		requestHandler: requestHandler,
	}
}
