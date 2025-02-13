package platformapiv1

import "github.com/google/wire"

var Set = wire.NewSet(
	NewHandler,
	NewServer,
)
