package wire

import (
	"bookstack/cmd/service2/internal/controller.go"

	"github.com/google/wire"
)

var ControllerSet = wire.NewSet(
	controller.NewShipperController,
)
