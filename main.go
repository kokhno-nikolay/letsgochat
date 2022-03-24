package main

import (
	"log"
	"os"

	"github.com/kokhno-nikolay/letsgochat/api"
	"github.com/kokhno-nikolay/letsgochat/repository"
	"github.com/kokhno-nikolay/letsgochat/repository/postgres"
)

func main() {
	db, err := postgres.NewClient(os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Println(err.Error())
	}

	repos := repository.NewRepositories(db)
	handler := api.NewHandler(api.Deps{Repos: repos})
	router := handler.Init()

	router.Run(":" + os.Getenv("PORT"))
}
