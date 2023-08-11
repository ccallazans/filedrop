package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ccallazans/filedrop/internal/domain"
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
		log.Fatal("asdasdasd")
	}

	u := domain.User{}

	pgdb.Preload("Role").Find(&u, 1)

	fmt.Println(u.Role)
}
