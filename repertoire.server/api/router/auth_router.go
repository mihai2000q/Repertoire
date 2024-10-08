package router

import (
	"repertoire/api/handler"
	"repertoire/api/server"
)

type AuthRouter struct {
	// logger      logger.Logger
	requestHandler *server.RequestHandler
	handler        *handler.AuthHandler
}

func (u AuthRouter) RegisterRoutes() {
	// u.logger.Info("Setting up auth routes")
	api := u.requestHandler.PublicRouter.Group("/auth")
	{
		api.PUT("/sign-in", u.handler.SignIn)
		api.POST("/sign-up", u.handler.SignUp)
	}
}

func NewAuthRouter(
	//logger logger.Logger,
	requestHandler *server.RequestHandler,
	handler *handler.AuthHandler,
) AuthRouter {
	return AuthRouter{
		//logger: 		logger,
		handler:        handler,
		requestHandler: requestHandler,
	}
}
