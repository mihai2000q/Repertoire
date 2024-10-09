package main

import (
	"go.uber.org/fx"
	"repertoire/api"
	"repertoire/data"
	"repertoire/domain"
	"repertoire/utils"
)

func main() {
	fx.New(
		fx.Provide(utils.NewEnv),
		data.Module,
		domain.Module,
		api.Module,
	).Run()
}
