package repository

import (
	"github.com/google/wire"
	"github.com/novychok/flagroll/platform/internal/repository/user"
)

var Set = wire.NewSet(
	user.NewPostgres,
)
