package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/kokhno-nikolay/letsgochat/repository/postgres"
)

// @title Lets go chat
// @version 1.0
// @description Online chat in golang

// @host letsgochat.herokuapp.com
// @BasePath /
func main() {
	runtime.GOMAXPROCS(1)

	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DBNAME"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_SSLMODE"),
	)
	db, err := postgres.NewClient(dns)
	if err != nil {
		panic(err.Error())
	}

	handlers := Wire(db)
	router := handlers.Init()

	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.Run(":" + os.Getenv("PORT"))
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading local .env file. Now used production(heroku) env file.")
	}
}
