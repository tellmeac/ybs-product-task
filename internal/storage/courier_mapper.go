package storage

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"yandex-team.ru/bstask/internal/core/entities"
)

type CourierMapper struct {
	Storage *Storage
}

func (m CourierMapper) All(ctx context.Context, limit uint64, offset uint64) ([]entities.Courier, error) {
	rows, err := m.Storage.Database.Select(ctx,
		sq.Select("*").From("couriers").
			Limit(limit).Offset(offset))
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
	rows, err := m.Storage.Database.Select(ctx, sq.Select("*").From("couriers").
		Where(sq.Eq{
			"id": id,
		}))
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, nil
	}

	result, err := toCourier(rows)
	return &result, err
}

func (m CourierMapper) Insert(ctx context.Context, couriers []entities.Courier) ([]entities.Courier, error) {
	builder := sq.Insert("couriers").
		Columns("courier_type", "regions", "working_hours").
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")
	for i := range couriers {
		builder = builder.Values(couriers[i].Type, couriers[i].Regions, couriers[i].WorkingHours)
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
