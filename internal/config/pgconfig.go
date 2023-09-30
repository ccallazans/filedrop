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
	"github.com/spf13/viper"
)

type DBConfig struct {
	host     string
	port     int
	db       string
	user     string
	password string
}

func NewPostgresConn() (*sql.DB, error) {
	config := DBConfig{
		host:     viper.GetString("db.host"),
		port:     viper.GetInt("db.port"),
		db:       viper.GetString("db.db"),
		user:     viper.GetString("db.user"),
		password: viper.GetString("db.password"),
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.user, config.password, config.host, config.port, config.db,
	)

	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(db *sql.DB) error {
	config := DBConfig{
		host:     viper.GetString("db.host"),
		port:     viper.GetInt("db.port"),
		db:       viper.GetString("db.db"),
		user:     viper.GetString("db.user"),
		password: viper.GetString("db.password"),
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.user, config.password, config.host, config.port, config.db,
	)

	workDir, _ := os.Getwd()
	for !strings.HasSuffix(workDir, "filedrop") {
		workDir = filepath.Dir(workDir)
	}

	migrationsPath := fmt.Sprintf("file:///%s/migrations", workDir)

	m, err := migrate.New(migrationsPath, connStr)
	if err != nil {
		return err
	}
	m.Up()

	return nil
}
