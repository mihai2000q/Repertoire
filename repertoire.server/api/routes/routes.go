package routes

import (
	"context"
	"go.uber.org/fx"
	"repertoire/api/router"
)

type Routes []Route

type Route interface {
	SetupRoutes()
}

func NewRoutes(
	lc fx.Lifecycle,
	authRouter router.AuthRouter,
	userRouter router.UserRouter,
) *Routes {
	routes := &Routes{
		authRouter,
		userRouter,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			routes.setup()
			return nil
		},
	})

	return routes
}

func (r Routes) setup() {
	for _, route := range r {
		route.SetupRoutes()
	}
}
