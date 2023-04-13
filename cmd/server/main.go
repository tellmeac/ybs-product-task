package main

import (
	"context"
	"log"
	"yandex-team.ru/bstask/internal/core"
	"yandex-team.ru/bstask/internal/pkg/config"
	"yandex-team.ru/bstask/internal/server"
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

	repository, err := core.NewRepository(ctx, cfg)
	if err != nil {
		log.Fatalf("Init repository: %s", err)
	}

	app := server.New(repository)

	if err := app.Start(ctx); err != nil {
		log.Fatalf(err.Error())
	}
}
