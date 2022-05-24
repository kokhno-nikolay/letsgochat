//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/kokhno-nikolay/letsgochat/api"
	"github.com/kokhno-nikolay/letsgochat/repository"
	"github.com/kokhno-nikolay/letsgochat/services"
	"gorm.io/gorm"
)

var (
	ProviderSet wire.ProviderSet = wire.NewSet(
		api.NewHandler,
		services.NewServices,
		repository.NewRepositories,
	)
)

func Wire(db *gorm.DB) *repository.Repositories {
	panic(wire.Build(ProviderSet))
}
