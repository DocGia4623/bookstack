//go:build wireinject
// +build wireinject

package wire

import (
	"bookstack/cmd/service2/internal/controller.go"
	"bookstack/config"

	"github.com/google/wire"
)

var AppSet = wire.NewSet(
	config.LoadConfig,
	config.ConnectDB,
	RepositorySet,
	ServiceSet,
	ControllerSet,
	wire.Struct(new(App), "*"),
)

type App struct {
	ShipperController *controller.ShipperController
}

func InitializeApp() (*App, error) {
	wire.Build(AppSet)
	return nil, nil
}
