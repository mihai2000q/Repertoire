package main

import (
	"go.uber.org/fx"
	"repertoire/storage/api"
	"repertoire/storage/internal"
)

func main() {
	fx.New(
		internal.Module,
		api.Module,
	).Run()
}
