package main

import (
	"go.uber.org/fx"
	"repertoire/storage/api"
	"repertoire/storage/domain"
	"repertoire/storage/utils"
)

func main() {
	fx.New(
		fx.Provide(utils.NewEnv),
		domain.Module,
		api.Module,
	).Run()
}
