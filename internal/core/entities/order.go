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
	CompletedTime *time.Time       `json:"completed_time,omitempty"`
	CourierID     *int64
}

type CompleteInfo struct {
	OrderID      int64     `json:"order_id"`
	CourierID    int64     `json:"courier_id"`
	CompleteTime time.Time `json:"complete_time"`
}
