package main

import (
	"go.uber.org/fx"
	"repertoire/storage/api"
	"repertoire/storage/internal"
)

func main() {
	fx.New(
		fx.Provide(internal.NewEnv),
		api.Module,
	).Run()
}
