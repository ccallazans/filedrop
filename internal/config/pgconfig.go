package config

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
)

func NewPostgresConn() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")

	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(db *sql.DB) error {
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
