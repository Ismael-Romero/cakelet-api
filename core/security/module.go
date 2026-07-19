package security

import (
	"github.com/Ismael-Romero/cakelet-suite/core/security/factory"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		factory.NewPasswordFactory,
		factory.NewTokenFactory,
		factory.NewMFAFactory,
	),
)
