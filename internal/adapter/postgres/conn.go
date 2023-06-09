package postgresql

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConn() (*gorm.DB, error) {

	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_host := os.Getenv("DB_HOST")
	db_database := os.Getenv("DB_DATABASE")

	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
		db_user, db_password, db_host, db_database,
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
