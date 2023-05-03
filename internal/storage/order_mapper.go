package storage

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
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
	CourierID    *int64
	CompleteTime *time.Time
}

type OrderFilterParams struct {
	OrderID   int64
	CourierID *int64
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
			sq.GtOrEq{"completed_time": params.From},
			sq.Lt{"completed_time": params.To},
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

func (m *OrderMapper) Save(
	ctx context.Context, filter OrderFilterParams, params OrderSaveParams,
) (*entities.Order, error) {
	query := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": filter.OrderID}).
		Suffix("RETURNING *")

	if filter.CourierID != nil {
		query = query.Where(sq.Eq{"courier_id": filter.CourierID})
	}

	if params.CourierID != nil {
		query = query.Set("courier_id", *params.CourierID)
	}

	if params.CompleteTime != nil {
		query = query.Set("completed_time", *params.CompleteTime)
	}

	result, err := m.executeQuery(ctx, query)
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
