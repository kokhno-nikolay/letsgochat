package main

import (
	"github.com/kokhno-nikolay/letsgochat/api"
	"os"
)

func main() {
	handler := api.NewHandler()
	router := handler.Init()

	router.Run(os.Getenv("PORT"))
}
