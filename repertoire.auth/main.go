package main

import (
	"go.uber.org/fx"
	"repertoire/auth/api"
	"repertoire/auth/data"
	"repertoire/auth/domain"
	"repertoire/auth/internal"
)

func main() {
	fx.New(
		internal.Module,
		data.Module,
		domain.Module,
		api.Module,
	).Run()
}
