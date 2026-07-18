package main

import "go.uber.org/fx"

func main() {
	app := fx.New(
		fx.Provide(NewConfig),
		fx.Invoke(TestConfig),
		fx.Invoke(Setup),
	)
	app.Run()
}
