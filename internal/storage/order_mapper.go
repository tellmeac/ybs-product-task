package storage

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"yandex-team.ru/bstask/internal/core/entities"
)

type OrderMapper struct {
	Storage *Storage
}

func (m *OrderMapper) All(ctx context.Context, limit, offset uint64) ([]entities.Order, error) {
	rows, err := m.Storage.Database.Select(ctx,
		sq.Select("*").From("orders").Limit(limit).Offset(offset))
	if err != nil {
		return nil, err
	}

	result := make([]entities.Order, 0)
	for rows.Next() {
		order, err := toOrder(rows)
		if err != nil {
			return nil, err
		}

		result = append(result, order)
	}
	return result, nil
}

func (m *OrderMapper) Get(ctx context.Context, id int64) (*entities.Order, error) {
	rows, err := m.Storage.Database.Select(ctx,
		sq.Select("*").From("orders").
			Where(sq.Eq{
				"id": id,
			}))
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, nil
	}

	result, err := toOrder(rows)
	return &result, err
}

func (m *OrderMapper) Insert(ctx context.Context, orders []entities.Order) ([]entities.Order, error) {
	builder := sq.Insert("orders").
		Columns("weight", "region", "delivery_hours", "cost").
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	for i := range orders {
		builder = builder.Values(orders[i].Weight, orders[i].Region, orders[i].DeliveryHours, orders[i].Cost)
	}

	rows, err := m.Storage.Database.Insert(ctx, builder)
	if err != nil {
		return nil, err
	}

	var ind, id int64
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		orders[ind].ID = id
	}

	return orders, nil
}

func (m *OrderMapper) Update(ctx context.Context, orders []entities.Order) error {
	err := m.Storage.Database.Tx(ctx, func(ctx context.Context) error {
		// update only updatable fields (complete_time, courier_id) with returning other fields
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func toOrder(rows pgx.Rows) (entities.Order, error) {
	var order entities.Order
	err := rows.Scan(&order.ID, &order.Weight, &order.Region,
		&order.DeliveryHours, &order.Cost, &order.CompletedTime, &order.CourierID)
	if err != nil {
		return entities.Order{}, err
	}

	return order, nil
}
