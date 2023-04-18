package entities

import (
	"yandex-team.ru/bstask/internal/pkg/types"
)

type CourierType string

var CourierTypes = []CourierType{FootCourier, BikeCourier, AutoCourier}

var (
	FootCourier CourierType = "FOOT"
	BikeCourier CourierType = "BIKE"
	AutoCourier CourierType = "AUTO"
)

type Courier struct {
	ID           int64            `json:"courier_id"`
	Type         CourierType      `json:"courier_type"`
	Regions      []int32          `json:"regions"`
	WorkingHours []types.Interval `json:"working_hours"`
}

type CourierMeta struct {
	ID           int64            `json:"courier_id"`
	Type         CourierType      `json:"courier_type"`
	Regions      []int32          `json:"regions"`
	WorkingHours []types.Interval `json:"working_hours"`
	Earnings     int32            `json:"earnings,omitempty"`
	Rating       int32            `json:"rating,omitempty"`
}
