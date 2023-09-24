package main

import (
	"log"

	"github.com/ccallazans/filedrop/internal/api/v1"
	"github.com/ccallazans/filedrop/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConn, err := config.NewPostgresConn("")
	if err != nil {
		log.Fatal(err)
	}

	api, err := api.NewApi(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	router := api.Routes()
	err = router.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
