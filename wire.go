//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"
	"github.com/google/wire"
	"github.com/kokhno-nikolay/letsgochat/api"
	"github.com/kokhno-nikolay/letsgochat/repository"
	"github.com/kokhno-nikolay/letsgochat/services"
)

var (
	ProviderSet wire.ProviderSet = wire.NewSet(
		api.NewHandler,
		services.NewServices,
		repository.NewRepositories,
	)
)

func Wire(db *sql.DB) *repository.Repositories {
	panic(wire.Build(ProviderSet))
}
