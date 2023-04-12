package actions

import (
	"context"
	"yandex-team.ru/bstask/internal/core/entities"
	"yandex-team.ru/bstask/internal/storage"
)

type Actions struct {
	storage *storage.Storage
}

func NewActions(storage *storage.Storage) *Actions {
	return &Actions{
		storage: storage,
	}
}

func (a *Actions) GetOrders(ctx context.Context, limit, offset uint64) ([]entities.Order, error) {
	return a.storage.Orders.All(ctx, limit, offset)
}

func (a *Actions) GetOrder(ctx context.Context, id int64) (*entities.Order, error) {
	return a.storage.Orders.Get(ctx, id)
}

func (a *Actions) CreateOrders(ctx context.Context, requests []entities.Order) ([]entities.Order, error) {
	return a.CreateOrders(ctx, requests)
}

func (a *Actions) CreateCouriers(ctx context.Context, requests []entities.Courier) ([]entities.Courier, error) {
	return a.storage.Couriers.Create(ctx, requests)
}

func (a *Actions) GetCouriers(ctx context.Context, limit uint64, offset uint64) ([]entities.Courier, error) {
	return a.storage.Couriers.All(ctx, limit, offset)
}

func (a *Actions) GetCourier(ctx context.Context, id int64) (*entities.Courier, error) {
	return a.storage.Couriers.Get(ctx, id)
}
