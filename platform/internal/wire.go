//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/novychok/flagroll/platform/internal/handler/platformapiv1"
	"github.com/novychok/flagroll/platform/internal/pkg"
)

func Init() (*App, func(), error) {
	wire.Build(
		pkg.Set,
		platformapiv1.Set,
		New,
	)

	return &App{}, func() {}, nil
}
