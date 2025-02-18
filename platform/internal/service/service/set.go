package service

import (
	"github.com/google/wire"
	"github.com/novychok/flagroll/platform/internal/service/apikeys"
	"github.com/novychok/flagroll/platform/internal/service/authorization"
	featureflag "github.com/novychok/flagroll/platform/internal/service/featureFlag"
	"github.com/novychok/flagroll/platform/internal/service/realtime"
	"github.com/novychok/flagroll/platform/internal/service/users"
)

var Set = wire.NewSet(
	users.New,
	authorization.New,
	featureflag.New,
	apikeys.New,
	realtime.New,
)
