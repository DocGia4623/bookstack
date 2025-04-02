package wire

import (
	"bookstack/internal/controller"

	"github.com/google/wire"
)

var ControllerSet = wire.NewSet(
	controller.NewAuthenticationController,
	controller.NewUserController,
	controller.NewBookController,
	controller.NewOrderController,
	controller.NewShipperController,
)
