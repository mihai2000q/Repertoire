package main

import (
	"repertoire/auth/api"
	"repertoire/auth/data"
	"repertoire/auth/domain"
	"repertoire/auth/internal"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		internal.Module,
		data.Module,
		domain.Module,
		api.Module,
	).Run()
}
