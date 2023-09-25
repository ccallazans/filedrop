package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConn() (*gorm.DB, error) {
	connString := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	connString := os.Getenv("DATABASE_URL")

	workDir, _ := os.Getwd()
	for !strings.HasSuffix(workDir, "filedrop") {
		workDir = filepath.Dir(workDir)
	}

	migrationsPath := fmt.Sprintf("file:///%s/migrations", workDir)

	m, err := migrate.New(migrationsPath, connString)
	if err != nil {
		return err
	}
	m.Up()

	return nil
}
