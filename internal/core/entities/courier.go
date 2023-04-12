package entities

import (
	"yandex-team.ru/bstask/internal/pkg/types"
)

type CourierType string

var (
	Foot CourierType = "FOOT"
	Bike CourierType = "BIKE"
	Auto CourierType = "AUTO"
)

type Courier struct {
	ID           int64        `json:"courier_id"`
	Type         CourierType  `json:"courier_type"`
	Regions      []int32      `json:"regions"`
	WorkingHours []types.Hour `json:"working_hours"`
}
