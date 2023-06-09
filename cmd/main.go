package main

import (
	"log"

	"github.com/ccallazans/filedrop/internal/adapter/handler"
	postgresql "github.com/ccallazans/filedrop/internal/adapter/postgres"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pgdb, err := postgresql.NewPostgresConn()
	if err != nil {
		log.Fatal(err)
	}

	router := handler.NewRouter(pgdb)
	router.Start(":8080")
}
