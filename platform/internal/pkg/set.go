package pkg

import (
	"github.com/google/wire"
	"github.com/novychok/flagroll/platform/internal/pkg/slog"
)

var Set = wire.NewSet(
	// httpclient.New,
	// postgres.New,
	slog.New,
	// validator.New,
	// jwts.New,
)
