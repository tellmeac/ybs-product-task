package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	ConnectionStr string `yaml:"connectionStr"`
}

type Storage struct {
	Config Config
	Pool   *pgxpool.Pool

	Couriers CourierMapper
	Orders   OrderMapper
}

func NewStorage(ctx context.Context, cfg Config) (*Storage, error) {
	var err error

	storage := &Storage{
		Config: cfg,
	}

	storage.Pool, err = pgxpool.New(ctx, cfg.ConnectionStr)
	if err != nil {
		return nil, err
	}

	storage.Couriers = CourierMapper{Storage: storage}
	storage.Orders = OrderMapper{Storage: storage}

	return storage, nil
}
