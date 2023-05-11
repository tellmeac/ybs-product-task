package actions

import (
	"context"
	"errors"
	"fmt"
	"time"
	"yandex-team.ru/bstask/internal/core/actions/meta"
	"yandex-team.ru/bstask/internal/core/actions/validators"
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
	if len(requests) == 0 {
		return errors.New("must provide at least one request item")
	}

	for i := range requests {
		if err := validators.ValidateOrder(&requests[i]); err != nil {
			return fmt.Errorf("validate %d's request: %w", i, err)
		}
	}

	return nil
}

func (a *Actions) CreateOrders(ctx context.Context, requests []entities.Order) ([]entities.Order, error) {
	result := make([]entities.Order, 0)

	err := a.storage.Database.Tx(ctx, func(ctx context.Context) error {
		for _, r := range requests {
			order, err := a.storage.Orders.Insert(ctx, storage.OrderCreateParams{
				Weight:        r.Weight,
				Region:        r.Region,
				DeliveryHours: r.DeliveryHours,
				Cost:          r.Cost,
			})
			if err != nil {
				return err
			}

			result = append(result, *order)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *Actions) CompleteOrder(ctx context.Context, requests []entities.CompleteInfo) ([]entities.Order, error) {
	result := make([]entities.Order, 0, len(requests))

	err := a.storage.Database.Tx(ctx, func(ctx context.Context) error {
		for _, r := range requests {
			order, err := a.storage.Orders.Save(ctx,
				storage.OrderFilterParams{
					OrderID:   r.OrderID,
					CourierID: &r.CourierID,
				},
				storage.OrderSaveParams{
					CompleteTime: &r.CompleteTime,
				},
			)
			if err != nil {
				return err
			}
			if order == nil {
				return fmt.Errorf("failed to complete order %d: %w", r.OrderID, ErrCompleteOrder)
			}

			result = append(result, *order)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *Actions) ValidateCreateCouriers(requests []entities.Courier) error {
	if len(requests) == 0 {
		return errors.New("must provide at least one request item")
	}

	for i := range requests {
		if err := validators.ValidateCourier(&requests[i]); err != nil {
			return fmt.Errorf("validate %d's request: %w", i, err)
		}
	}

	return nil
}

func (a *Actions) CreateCouriers(ctx context.Context, requests []entities.Courier) ([]entities.Courier, error) {
	result := make([]entities.Courier, 0)

	err := a.storage.Database.Tx(ctx, func(ctx context.Context) error {
		for _, r := range requests {
			order, err := a.storage.Couriers.Insert(ctx, storage.CourierCreateParams{
				Type:         r.Type,
				Regions:      r.Regions,
				WorkingHours: r.WorkingHours,
			})
			if err != nil {
				return err
			}

			result = append(result, *order)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *Actions) GetCouriers(ctx context.Context, limit uint64, offset uint64) ([]entities.Courier, error) {
	return a.storage.Couriers.All(ctx, limit, offset)
}

func (a *Actions) GetCourier(ctx context.Context, id int64) (*entities.Courier, error) {
	return a.storage.Couriers.Get(ctx, id)
}

func (a *Actions) GetCourierMetaInfo(
	ctx context.Context, courierId int64, startDate, endDate time.Time,
) (*entities.CourierMeta, error) {
	var result *entities.CourierMeta
	err := a.storage.Database.ReadonlyTx(ctx, func(ctx context.Context) error {
		courier, err := a.storage.Couriers.Get(ctx, courierId)
		if err != nil {
			return err
		}
		if courier == nil {
			return nil
		}

		orders, err := a.storage.Orders.Find(ctx, storage.OrderFindParams{
			From:      startDate,
			To:        endDate,
			CourierID: courier.ID,
		})
		if err != nil {
			return err
		}

		result = meta.GetCourierMeta(courier, orders, startDate, endDate)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
