package main

import (
	"log"

	"github.com/ccallazans/filedrop/internal/api/v1"
	"github.com/ccallazans/filedrop/internal/config"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := config.NewPostgresConn()
	if err != nil {
		log.Fatal(err)
	}

	err = config.RunMigrations(db)
	if err != nil {
		log.Fatal(err)
	}

	api, err := api.NewApi(db)
	if err != nil {
		log.Fatal(err)
	}

	router := api.Routes()
	err = router.Start(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
