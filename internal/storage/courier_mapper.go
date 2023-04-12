package storage

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"yandex-team.ru/bstask/internal/core/entities"
)

type CourierMapper struct {
	Storage *Storage
}

func (m CourierMapper) All(ctx context.Context, limit uint64, offset uint64) ([]entities.Courier, error) {
	query, _, err := squirrel.Select("*").From("couriers").
		PlaceholderFormat(squirrel.Dollar).Limit(limit).Offset(offset).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := m.Storage.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	result := make([]entities.Courier, 0)
	for rows.Next() {
		courier, err := toCourier(rows)
		if err != nil {
			return nil, err
		}

		result = append(result, courier)
	}
	return result, nil
}

func (m CourierMapper) Get(ctx context.Context, id int64) (*entities.Courier, error) {
	query, args, err := squirrel.Select("*").From("couriers").
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

	result, err := toCourier(rows)
	return &result, err
}

func (m CourierMapper) Create(ctx context.Context, couriers []entities.Courier) ([]entities.Courier, error) {
	builder := squirrel.Insert("couriers").
		Columns("courier_type", "regions", "working_hours").
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id")
	for i := range couriers {
		builder = builder.Values(couriers[i].Type, couriers[i].Regions, couriers[i].WorkingHours)
	}
	query, args, err := builder.ToSql()
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

		couriers[ind].ID = id
	}

	return couriers, nil
}

func toCourier(rows pgx.Rows) (entities.Courier, error) {
	var courier entities.Courier
	err := rows.Scan(&courier.ID, &courier.Type, &courier.Regions, &courier.WorkingHours)
	if err != nil {
		return entities.Courier{}, err
	}

	return courier, nil
}
