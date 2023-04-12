package storage

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"yandex-team.ru/bstask/internal/core/entities"
)

type OrderMapper struct {
	Storage *Storage
}

func (m *OrderMapper) All(ctx context.Context, limit, offset uint64) ([]entities.Order, error) {
	query, _, err := squirrel.Select("*").From("orders").
		Limit(limit).Offset(offset).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := m.Storage.Pool.Query(ctx, query)
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
	query, args, err := squirrel.Select("*").From("orders").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{
			"id": id,
		}).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := m.Storage.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, nil
	}

	result, err := toOrder(rows)
	return &result, err
}

func (m *OrderMapper) Create(ctx context.Context, orders []entities.Order) ([]entities.Order, error) {
	builder := squirrel.Insert("orders").Columns("weight", "region", "delivery_hours", "cost")
	for i := range orders {
		builder.Values(orders[i].Weight, orders[i].Region, orders[i].DeliveryHours, orders[i].Cost)
	}
	query, args, err := builder.Suffix("returning id").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := m.Storage.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var (
		ind, id int64
	)
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		orders[ind].ID = id
	}

	return orders, nil
}

func toOrder(rows pgx.Rows) (entities.Order, error) {
	var order entities.Order
	err := rows.Scan(&order.ID, &order.Weight, &order.Region, &order.DeliveryHours, &order.Cost, &order.CompletedTime)
	if err != nil {
		return entities.Order{}, err
	}

	return order, nil
}
