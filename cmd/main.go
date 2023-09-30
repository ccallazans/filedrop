package main

import (
	"log"

	"github.com/ccallazans/filedrop/internal/api/v1"
	"github.com/ccallazans/filedrop/internal/config"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("Error loading config.json")
	}

	err = godotenv.Load()
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
	err = router.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
