package wire

import (
	"bookstack/internal/service"

	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(
	service.NewUserServiceImpl,
	service.NewAuthServiceImpl,
	service.NewPermissionRepositoryImpl,
	service.NewBookServiceImpl,
	service.NewOrderServiceImpl,
	service.NewOrderManageService,
)
