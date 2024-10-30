package main

import (
	"go.uber.org/fx"
	"repertoire/server/api"
	"repertoire/server/data"
	"repertoire/server/domain"
	"repertoire/server/utils"
)

func main() {
	fx.New(
		fx.Provide(utils.NewEnv),
		data.Module,
		domain.Module,
		api.Module,
	).Run()
}
