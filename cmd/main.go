package main

import (
	"log"
	"os"

	"github.com/ccallazans/filedrop/internal/api"
	"github.com/ccallazans/filedrop/internal/config"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Lugu *zap.Logger

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pgdb, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	api := api.NewApi(config.NewLogger("api"), pgdb)

	router := api.Routes()
	err = router.Start(":8081")
	if err != nil {
		log.Fatal(err)
	}
}
