package pkg

import (
	"github.com/google/wire"
	"github.com/novychok/flagroll/platform/internal/pkg/httpclient"
	"github.com/novychok/flagroll/platform/internal/pkg/jwts"
	"github.com/novychok/flagroll/platform/internal/pkg/postgres"
	"github.com/novychok/flagroll/platform/internal/pkg/slog"
	"github.com/novychok/flagroll/platform/internal/pkg/validator"
)

var Set = wire.NewSet(
	httpclient.New,
	postgres.New,
	slog.New,
	validator.New,
	jwts.New,
)
