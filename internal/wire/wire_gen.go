// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"bookstack/config"
	"bookstack/internal/controller"
	"bookstack/internal/middleware"
	"bookstack/internal/repository"
	"bookstack/internal/service"
	"github.com/google/wire"
)

// Injectors from injector.go:

// InitializeUserService khởi tạo UserService tự động
func InitializeApp() (*App, error) {
	configConfig, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	db := config.ConnectDB(configConfig)
	userRepository := repository.NewUserRepositoryImpl(db, configConfig)
	authService := service.NewAuthServiceImpl(userRepository, configConfig)
	authenticationController := controller.NewAuthenticationController(authService)
	userService := service.NewUserServiceImpl(userRepository)
	userController := controller.NewUserController(userService)
	bookRepository := repository.NewBookRepositoryImpl(db)
	bookService := service.NewBookServiceImpl(bookRepository)
	bookController := controller.NewBookController(bookService, userService)
	orderRepository := repository.NewOrderRepositoryImpl(db)
	orderService := service.NewOrderServiceImpl(orderRepository)
	orderController := controller.NewOrderController(orderService, userService)
	permissionRepository := repository.NewPermissionRepositoryImpl(db)
	middlewareMiddleware := middleware.NewAuthorizeMiddleware(userRepository, permissionRepository, configConfig)
	shipperRepository := repository.NewShipperRepository(db)
	shipperOrderManageService := service.NewOrderManageService(shipperRepository)
	shipperController := controller.NewShipperController(shipperOrderManageService, userService)
	app := &App{
		AuthenticationController: authenticationController,
		UserController:           userController,
		BookController:           bookController,
		OrderController:          orderController,
		Middleware:               middlewareMiddleware,
		ShipperController:        shipperController,
	}
	return app, nil
}

// injector.go:

var AppSet = wire.NewSet(config.LoadConfig, config.ConnectDB, config.ConnectRedis, RepositorySet,
	MiddlerwareSet,
	ServiceSet,
	ControllerSet, wire.Struct(new(App), "*"),
)

type App struct {
	AuthenticationController *controller.AuthenticationController
	UserController           *controller.UserController
	BookController           *controller.BookController
	OrderController          *controller.OrderController
	Middleware               *middleware.Middleware
	ShipperController        *controller.ShipperController
}
