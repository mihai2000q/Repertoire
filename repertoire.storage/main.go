package main

import (
	"go.uber.org/fx"
	"repertoire/storage/api"
	"repertoire/storage/utils"
)

func main() {
	fx.New(
		fx.Provide(utils.NewEnv),
		api.Module,
	).Run()
}
