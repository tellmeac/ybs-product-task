package storage

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"yandex-team.ru/bstask/internal/core/entities"
	"yandex-team.ru/bstask/internal/pkg/types"
)

type CourierMapper struct {
	Storage *Storage
}

type CourierCreateParams struct {
	Type         entities.CourierType
	Regions      []int32
	WorkingHours []types.Interval
}

func (m *CourierMapper) executeQuery(ctx context.Context, query sq.Sqlizer) ([]entities.Courier, error) {
	rows, err := m.Storage.Database.QuerySq(ctx, query)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]entities.Courier, 0)
	for rows.Next() {
		order, err := toCourier(rows)
		if err != nil {
			return nil, err
		}

		result = append(result, order)
	}

	return result, nil
}

func (m *CourierMapper) All(ctx context.Context, limit uint64, offset uint64) ([]entities.Courier, error) {
	return m.executeQuery(ctx, sq.Select("*").From("couriers").
		PlaceholderFormat(sq.Dollar).
		Limit(limit).Offset(offset))
}

func (m *CourierMapper) Get(ctx context.Context, id int64) (*entities.Courier, error) {
	result, err := m.executeQuery(ctx, sq.Select("*").From("couriers").
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

func (m *CourierMapper) Insert(ctx context.Context, params CourierCreateParams) (*entities.Courier, error) {
	result, err := m.executeQuery(ctx, sq.Insert("couriers").
		PlaceholderFormat(sq.Dollar).
		Columns("courier_type", "regions", "working_hours").
		Values(params.Type, params.Regions, params.WorkingHours).
		Suffix("RETURNING *"))
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}

func toCourier(rows pgx.Rows) (entities.Courier, error) {
	var courier entities.Courier
	err := rows.Scan(&courier.ID, &courier.Type, &courier.Regions, &courier.WorkingHours)
	if err != nil {
		return entities.Courier{}, err
	}

	return courier, nil
}
