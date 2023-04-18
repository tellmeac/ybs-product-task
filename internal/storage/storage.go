package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"yandex-team.ru/bstask/internal/pkg/pgdb"
)

type Config struct {
	ConnectionStr string `yaml:"connectionStr"`
}

type Storage struct {
	Config   Config
	Database *pgdb.Database

	Couriers CourierMapper
	Orders   OrderMapper
}

func NewStorage(ctx context.Context, cfg Config) (*Storage, error) {
	var err error

	storage := &Storage{
		Config: cfg,
	}

	pool, err := pgxpool.New(ctx, cfg.ConnectionStr)
	if err != nil {
		return nil, err
	}
	storage.Database = pgdb.NewDatabase(pool)

	storage.Couriers = CourierMapper{Storage: storage}
	storage.Orders = OrderMapper{Storage: storage}

	return storage, nil
}
