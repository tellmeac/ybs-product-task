package entities

import (
	"time"
	"yandex-team.ru/bstask/internal/pkg/types"
)

type Order struct {
	ID            int64            `json:"order_id"`
	Weight        float64          `json:"weight"`
	Region        int32            `json:"regions"` // NOTE: One region for order.
	DeliveryHours []types.Interval `json:"delivery_hours"`
	Cost          int32            `json:"cost"`
	CompletedTime *time.Time       `json:"completed_time"`
}
