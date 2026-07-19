package main

import (
	"github.com/Ismael-Romero/cakelet-suite/config"
	"github.com/Ismael-Romero/cakelet-suite/core/security"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(config.NewConfig),
		fx.Invoke(config.TestConfig),
		fx.Invoke(Setup),
		security.Module,
	)
	app.Run()
}
