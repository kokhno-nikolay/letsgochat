// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"database/sql"
	"github.com/google/wire"
	"github.com/kokhno-nikolay/letsgochat/api"
	"github.com/kokhno-nikolay/letsgochat/repository"
	"github.com/kokhno-nikolay/letsgochat/services"
)

// Injectors from wire.go:

func Wire(db *sql.DB) *api.Handler {
	repositories := repository.NewRepositories(db)
	services := services.NewServices(services.Deps{Repos: repositories})
	handlers := api.NewHandler(services)
	return handlers
}

// wire.go:

var (
	ProviderSet wire.ProviderSet = wire.NewSet(api.NewHandler, services.NewServices, repository.NewRepositories)
)