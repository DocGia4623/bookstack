//go:build wireinject
// +build wireinject

package wire

import (
	"bookstack/config"
	"bookstack/internal/controller"
	"bookstack/internal/middleware"

	"github.com/google/wire"
)

var AppSet = wire.NewSet(
	config.LoadConfig,
	config.ConnectDB,
	RepositorySet,
	MiddlerwareSet, // Inject Middleware
	ServiceSet,
	ControllerSet,
	wire.Struct(new(App), "*"),
)

type App struct {
	AuthenticationController *controller.AuthenticationController
	UserController           *controller.UserController
	BookController           *controller.BookController
	OrderController          *controller.OrderController
	Middleware               *middleware.Middleware
}

// InitializeUserService khởi tạo UserService tự động
func InitializeApp() (*App, error) {
	wire.Build(AppSet)
	return nil, nil
}
