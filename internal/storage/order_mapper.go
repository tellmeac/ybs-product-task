package storage

import (
	"context"
	"fmt"
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

	orders := make([]entities.Order, 0)
	for rows.Next() {
		orders = append(orders, toOrder(rows))
	}
	return orders, nil
}

func toOrder(rows pgx.Rows) entities.Order {
	var order entities.Order
	err := rows.Scan(&order.ID, &order.Weight, &order.Region, &order.DeliveryHours, &order.Cost, &order.CompletedTime)
	if err != nil {
		panic(fmt.Errorf("scan order: %w", err))
	}

	return order
}
