package wire

import (
	"bookstack/cmd/service2/internal/service"

	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(
	service.NewOrderManageService,
)
