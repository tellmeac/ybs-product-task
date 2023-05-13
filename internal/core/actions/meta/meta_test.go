package meta

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
	"yandex-team.ru/bstask/internal/core/entities"
)

func makeCourier(t entities.CourierType) entities.Courier {
	return entities.Courier{Type: t}
}

func manyOrdersWith(cost int32, count int) []entities.Order {
	result := make([]entities.Order, 0, count)
	for ; count > 0; count-- {
		result = append(result, entities.Order{Cost: cost})
	}
	return result
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func TestGetCourierMeta(t *testing.T) {
	for _, tc := range []struct {
		Name            string
		Courier         entities.Courier
		CompletedOrders []entities.Order
		From            time.Time
		To              time.Time
		Expected        entities.CourierMeta
	}{
		{
			Name:            "Empty completed orders",
			Courier:         makeCourier(entities.AutoCourier),
			CompletedOrders: nil,
			Expected: entities.CourierMeta{
				Type: entities.AutoCourier,
			},
		},
		{
			Name:            "Default case",
			Courier:         makeCourier(entities.AutoCourier),
			CompletedOrders: manyOrdersWith(100, 24),
			Expected: entities.CourierMeta{
				Type:     entities.AutoCourier,
				Earnings: ref(int32(9600)), // (100 * 24) * 4;
				Rating:   ref(int32(1)),    // total two completed orders in one requested day;
			},
			From: date(2023, 4, 17),
			To:   date(2023, 4, 18),
		},
		{
			Name:            "Default case with FOOT courier",
			Courier:         makeCourier(entities.FootCourier),
			CompletedOrders: manyOrdersWith(100, 24),
			Expected: entities.CourierMeta{
				Type:     entities.FootCourier,
				Earnings: ref(int32(4800)), // (100 * 24) * 2;
				Rating:   ref(int32(3)),    // total two completed orders in one requested day;
			},
			From: date(2023, 4, 17),
			To:   date(2023, 4, 18),
		},
		{
			Name:            "Round date if range is too small",
			Courier:         makeCourier(entities.AutoCourier),
			CompletedOrders: manyOrdersWith(100, 24),
			Expected: entities.CourierMeta{
				Type:     entities.AutoCourier,
				Earnings: ref(int32(9600)), // (200 + 300) * 4;
				Rating:   ref(int32(0)),    // total two completed orders in two days;
			},
			From: date(2023, 4, 18),
			To:   date(2023, 4, 20),
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			got := GetCourierMeta(&tc.Courier, tc.CompletedOrders, tc.From, tc.To)

			require.True(t, reflect.DeepEqual(tc.Expected, *got))
		})
	}
}

func ref[T any](v T) *T {
	return &v
}
