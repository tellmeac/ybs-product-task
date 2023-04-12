package core

import (
	"context"
	"yandex-team.ru/bstask/internal/storage"
)

type Repository struct {
	Config  *Config
	Storage *storage.Storage
}

func NewRepository(cfg *Config) (*Repository, error) {
	var err error

	r := &Repository{
		Config: cfg,
	}

	r.Storage, err = storage.NewStorage(context.Background(), cfg.Storage)
	if err != nil {
		return nil, err
	}

	return r, nil
}
