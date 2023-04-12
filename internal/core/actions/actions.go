package actions

import (
	"context"
	"errors"
	"fmt"
	"yandex-team.ru/bstask/internal/core/actions/courier"
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

func (a *Actions) ValidateCreateOrders(requests []entities.Order) error {
	// TODO: validation for request: Weight, DeliveryHours, Cost, Region
	if len(requests) == 0 {
		return errors.New("must provide at least one request item")
	}

	return nil
}

func (a *Actions) CreateOrders(ctx context.Context, requests []entities.Order) ([]entities.Order, error) {
	return a.CreateOrders(ctx, requests)
}

func (a *Actions) ValidateCreateCouriers(requests []entities.Courier) error {
	if len(requests) == 0 {
		return errors.New("must provide at least one request item")
	}

	for i := range requests {
		if err := courier.Validate(&requests[i]); err != nil {
			return fmt.Errorf("validate %d's request: %w", i, err)
		}
	}

	return nil
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
