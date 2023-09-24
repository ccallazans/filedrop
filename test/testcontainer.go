package test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var ActivePostgresTestContainer *postgres.PostgresContainer

func NewPostgresTestContainer() *postgres.PostgresContainer {
	ctx := context.Background()
	wd, _ := os.Getwd()
	for !strings.HasSuffix(wd, "filedrop") {
		wd = filepath.Dir(wd)
	}

	if ActivePostgresTestContainer == nil {
		postgresContainer, err := postgres.RunContainer(ctx,
			testcontainers.WithImage("postgres:latest"),
			postgres.WithInitScripts(filepath.Join(wd, "migrations/000001_initial.up.sql")),
			postgres.WithDatabase("test"),
			postgres.WithUsername("test"),
			postgres.WithPassword("test"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(5*time.Second)),
		)
		if err != nil {
			panic(err)
		}

		ActivePostgresTestContainer = postgresContainer
	}

	return ActivePostgresTestContainer
}