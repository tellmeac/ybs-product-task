package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/avast/retry-go/v4"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"log"
	"time"
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

	err = retry.Do(func() error {
		return UpMigrations(cfg)
	}, retry.Attempts(4), retry.Delay(2*time.Second))
	if err != nil {
		log.Fatalf(err.Error())
	}

	repository, err := core.NewRepository(ctx, cfg)
	if err != nil {
		log.Fatalf("Init repository: %s", err)
	}

	app := server.New(repository)

	if err := app.Start(ctx); err != nil {
		log.Fatalf(err.Error())
	}
}

func UpMigrations(cfg *core.Config) error {
	db, err := sql.Open("pgx", cfg.Storage.URL)
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil

}
