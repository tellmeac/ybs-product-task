package main

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"yandex-team.ru/bstask/internal/core"
	"yandex-team.ru/bstask/internal/pkg/config"
)

func main() {
	loader := config.PrepareLoader(
		config.WithConfigPath("./config.yaml"),
	)

	cfg, err := core.ParseConfig(loader)
	if err != nil {
		log.Fatalf("Failed to parse config: %s", err)
	}

	m, err := migrate.New(
		"file://migrations",
		cfg.Storage.ConnectionStr,
	)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Failed to apply migrations: %s", err)
	}
}
