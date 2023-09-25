package config

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormPg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestContainerStruct struct {
	PostgresContainer *postgres.PostgresContainer
	DB                *gorm.DB
}

var activeTestContainer *TestContainerStruct

func NewTestContainerStruct() *TestContainerStruct {
	if activeTestContainer == nil {

		container := &TestContainerStruct{}
		container.activatePostgres()
		container.activateDB()

		RunMigrations(container.DB)
	}

	return activeTestContainer
}

func (t *TestContainerStruct) activatePostgres() {
	ctx := context.Background()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	for !strings.HasSuffix(wd, "filedrop") {
		wd = filepath.Dir(wd)
	}

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:latest"),
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatal(err)
	}

	t.PostgresContainer = postgresContainer
}

func (t *TestContainerStruct) activateDB() {
	connStr, err := t.PostgresContainer.ConnectionString(context.TODO(), "sslmode=disable", "application_name=test")
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(gormPg.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	t.DB = db
}
