package wire

import (
	"bookstack/internal/middleware"

	"github.com/google/wire"
)

var MiddlerwareSet = wire.NewSet(
	middleware.NewAuthorizeMiddleware,
)
