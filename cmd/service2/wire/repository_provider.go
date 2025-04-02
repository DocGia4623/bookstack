package wire

import (
	"bookstack/cmd/service2/internal/repository"

	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	repository.NewShipperRepository,
)
