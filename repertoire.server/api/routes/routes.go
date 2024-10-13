package routes

import (
	"context"
	"repertoire/api/router"

	"go.uber.org/fx"
)

type Routes []Route

type Route interface {
	RegisterRoutes()
}

func NewRoutes(
	lc fx.Lifecycle,
	authRouter router.AuthRouter,
	songRouter router.SongRouter,
	userRouter router.UserRouter,
) *Routes {
	routes := &Routes{
		authRouter,
		songRouter,
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
		route.RegisterRoutes()
	}
}
