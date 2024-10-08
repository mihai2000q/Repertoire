package router

import (
	"repertoire/api/handler"
	"repertoire/api/server"
)

type AuthRouter struct {
	// logger      server.Logger
	requestHandler *server.RequestHandler
	handler        *handler.AuthHandler
}

func (u AuthRouter) SetupRoutes() {
	// u.logger.Debug("Setting up auth routes")
	api := u.requestHandler.BaseRouter.Group("/auth")
	{
		api.POST("/sign-up", u.handler.SignUp)
	}
}

func NewAuthRouter(
	//logger server.Logger,
	requestHandler *server.RequestHandler,
	handler *handler.AuthHandler,
) AuthRouter {
	return AuthRouter{
		//logger: 		logger,
		handler:        handler,
		requestHandler: requestHandler,
	}
}
