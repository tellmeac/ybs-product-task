package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"log"
	"yandex-team.ru/bstask/internal/core"
	"yandex-team.ru/bstask/internal/pkg/config"
	"yandex-team.ru/bstask/internal/server"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	ctx := context.Background()

	loader := config.PrepareLoader(
		config.WithConfigPath("./config.yaml"),
	)

	cfg, err := core.ParseConfig(loader)
	if err != nil {
		log.Fatalf("Failed to parse config: %s", err)
	}

	UpMigrations(cfg)

	repository, err := core.NewRepository(ctx, cfg)
	if err != nil {
		log.Fatalf("Init repository: %s", err)
	}

	app := server.New(repository)

	if err := app.Start(ctx); err != nil {
		log.Fatalf(err.Error())
	}
}

func UpMigrations(cfg *core.Config) {
	db, err := sql.Open("pgx", cfg.Storage.URL)
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf(err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Failed to apply migrations: %s", err)
	}
}
