package wire

import (
	"bookstack/internal/repository"

	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	repository.NewUserRepositoryImpl,
	repository.NewPermissionRepositoryImpl,
)
