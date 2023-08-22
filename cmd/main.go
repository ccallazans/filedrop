package main

import (
	"log"
	"os"

	"github.com/ccallazans/filedrop/internal/api"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pgdb, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	router := api.NewRouter(pgdb)
	router.Start(":8080")
}
