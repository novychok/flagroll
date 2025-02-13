//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/novychok/flagroll/platform/internal/config"
	"github.com/novychok/flagroll/platform/internal/handler/platformapiv1"
	"github.com/novychok/flagroll/platform/internal/pkg"
	"github.com/novychok/flagroll/platform/internal/repository/repository"
	"github.com/novychok/flagroll/platform/internal/service/service"
)

func Init() (*App, func(), error) {
	wire.Build(
		config.Set,
		pkg.Set,
		repository.Set,
		service.Set,
		platformapiv1.Set,
		New,
	)

	return &App{}, func() {}, nil
}
