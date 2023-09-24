package config

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConn(conn string) (*gorm.DB, error) {
	if conn == "" {
		conn = os.Getenv("DATABASE_URL")
	}

	pgdb, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return pgdb, nil
}
