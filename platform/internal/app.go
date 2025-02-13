package internal

import (
	"context"

	"github.com/novychok/flagroll/platform/internal/handler/platformapiv1"
)

type App struct {
	platformAPIV1 *platformapiv1.Server
}

func (a *App) Start(ctx context.Context) error {
	return nil
}

func New(platformAPIV1 *platformapiv1.Server) *App {
	return &App{
		platformAPIV1: platformAPIV1,
	}
}
