package repository

import (
	"github.com/google/wire"
	"github.com/novychok/flagroll/platform/internal/repository/apikey"
	featureflag "github.com/novychok/flagroll/platform/internal/repository/featureFlag"
	"github.com/novychok/flagroll/platform/internal/repository/user"
)

var Set = wire.NewSet(
	user.NewPostgres,
	featureflag.NewPostgres,
	apikey.NewPostgres,
)
