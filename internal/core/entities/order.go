package entities

import (
	"time"
	"yandex-team.ru/bstask/internal/pkg/types"
)

type Order struct {
	ID            int          `json:"order_id"`
	Weight        int          `json:"weight"`
	Region        int          `json:"regions"` // NOTE: One region for order.
	DeliveryHours []types.Hour `json:"delivery_hours"`
	Cost          int          `json:"cost"`
	CompletedTime *time.Time   `json:"completed_time"`
}
