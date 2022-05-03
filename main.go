package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/kokhno-nikolay/letsgochat/api"
	"github.com/kokhno-nikolay/letsgochat/repository"
	"github.com/kokhno-nikolay/letsgochat/repository/postgres"
)

// @title Lets go chat
// @version 1.0
// @description Online chat in golang

// @host https://letsgochat.herokuapp.com
// @BasePath /
func main() {
	db, err := postgres.NewClient(os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err.Error())
	}

	repos := repository.NewRepositories(db)
	handler := api.NewHandler(api.Deps{Repos: repos})
	router := handler.Init()

	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.Run(":" + os.Getenv("PORT"))
}
