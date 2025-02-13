// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package internal

import (
	"github.com/novychok/flagroll/platform/internal/config"
	"github.com/novychok/flagroll/platform/internal/handler/platformapiv1"
	"github.com/novychok/flagroll/platform/internal/pkg/jwts"
	"github.com/novychok/flagroll/platform/internal/pkg/postgres"
	"github.com/novychok/flagroll/platform/internal/pkg/slog"
	"github.com/novychok/flagroll/platform/internal/pkg/validator"
	"github.com/novychok/flagroll/platform/internal/repository/user"
	"github.com/novychok/flagroll/platform/internal/service/authorization"
)

// Injectors from wire.go:

func Init() (*App, func(), error) {
	logger := slog.New()
	configConfig, err := config.New()
	if err != nil {
		return nil, nil, err
	}
	platformapiv1Config := config.GetPlatfromAPIV1(configConfig)
	validate := validator.New()
	postgresConfig := config.GetPostgres(configConfig)
	connection, cleanup, err := postgres.New(postgresConfig)
	if err != nil {
		return nil, nil, err
	}
	repositoryUser := user.NewPostgres(connection)
	secretManager, err := jwts.New()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	serviceAuthorization := authorization.New(logger, validate, repositoryUser, secretManager)
	serverInterface := platformapiv1.NewHandler(serviceAuthorization)
	server := platformapiv1.NewServer(logger, platformapiv1Config, serviceAuthorization, serverInterface)
	app := New(server)
	return app, func() {
		cleanup()
	}, nil
}
