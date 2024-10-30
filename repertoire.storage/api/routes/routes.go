package routes

import (
	"context"
	"repertoire/storage/api/router"

	"go.uber.org/fx"
)

type Routes []Route

type Route interface {
	RegisterRoutes()
}

func NewRoutes(
	lc fx.Lifecycle,
	storageRouter router.StorageRouter,
) *Routes {
	routes := &Routes{
		storageRouter,
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
