package wire

import (
	"bookstack/internal/service"

	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(
	service.NewUserServiceImpl,
	service.NewAuthServiceImpl,
)
