package router

import (
	"repertoire/api/handler"
	"repertoire/api/server"
)

type UserRouter struct {
	// logger      logger.Logger
	requestHandler *server.RequestHandler
	handler        *handler.UserHandler
}

func (u UserRouter) RegisterRoutes() {
	// u.logger.Info("Setting up user routes")
	api := u.requestHandler.PrivateRouter.Group("/users")
	{
		api.GET("/", u.handler.GetUserByEmail)
	}
}

func NewUserRouter(
	//logger logger.Logger,
	requestHandler *server.RequestHandler,
	handler *handler.UserHandler,
) UserRouter {
	return UserRouter{
		//logger: 		logger,
		handler:        handler,
		requestHandler: requestHandler,
	}
}
