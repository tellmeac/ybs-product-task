package storage

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"time"
	"yandex-team.ru/bstask/internal/core/entities"
	"yandex-team.ru/bstask/internal/pkg/types"
)

type OrderCreateParams struct {
	Weight        float64
	Region        int32
	DeliveryHours []types.Interval
	Cost          int32
}

type OrderSaveParams struct {
	OrderID      int64
	CourierID    int64
	CompleteTime time.Time
}

type OrderMapper struct {
	Storage *Storage
}

func (m *OrderMapper) executeQuery(ctx context.Context, query sq.Sqlizer) ([]entities.Order, error) {
	rows, err := m.Storage.Database.QuerySq(ctx, query)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close()

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

type OrderFindParams struct {
	From, To  time.Time
	CourierID int64
}

func (m *OrderMapper) Find(ctx context.Context, params OrderFindParams) ([]entities.Order, error) {
	return m.executeQuery(ctx, sq.Select("*").From("orders").
		PlaceholderFormat(sq.Dollar).
		Where(sq.And{
			sq.Eq{"courier_id": params.CourierID},
			sq.LtOrEq{"completed_time": params.To},
			sq.GtOrEq{"completed_time": params.From},
		}))
}

func (m *OrderMapper) All(ctx context.Context, limit, offset uint64) ([]entities.Order, error) {
	return m.executeQuery(ctx, sq.Select("*").From("orders").
		PlaceholderFormat(sq.Dollar).
		Limit(limit).Offset(offset))
}

func (m *OrderMapper) Get(ctx context.Context, id int64) (*entities.Order, error) {
	result, err := m.executeQuery(ctx, sq.Select("*").From("orders").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}))
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}

func (m *OrderMapper) Insert(ctx context.Context, params OrderCreateParams) (*entities.Order, error) {
	result, err := m.executeQuery(ctx, sq.Insert("orders").
		Columns("weight", "region", "delivery_hours", "cost").
		Values(params.Weight, params.Region, params.DeliveryHours, params.Cost).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *"))
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}

func (m *OrderMapper) Save(ctx context.Context, params OrderSaveParams) (*entities.Order, error) {
	result, err := m.executeQuery(ctx, sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{
			"id":         params.OrderID,
			"courier_id": params.CourierID, // TODO: what about complete time assertion ? (idempotency)
		}).
		Set("completed_time", params.CompleteTime).
		Suffix("RETURNING *"))
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
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
