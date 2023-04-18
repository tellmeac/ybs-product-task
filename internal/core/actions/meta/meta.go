package meta

import (
	"fmt"
	"time"
	"yandex-team.ru/bstask/internal/core/entities"
)

func GetCourierMeta(
	courier *entities.Courier, completedOrders []entities.Order, from, to time.Time,
) *entities.CourierMeta {
	meta := entities.CourierMeta{
		ID:           courier.ID,
		Type:         courier.Type,
		Regions:      courier.Regions,
		WorkingHours: courier.WorkingHours,
	}

	// Ignore, if courier hasn't completed any orders.
	if len(completedOrders) == 0 {
		return &meta
	}

	for i := range completedOrders {
		meta.Earnings += completedOrders[i].Cost
	}
	meta.Earnings *= CourierEarningFactor(courier)

	meta.Rating = int32(len(completedOrders)) / int32(to.Sub(from).Hours()) * CourierRatingFactor(courier)

	return &meta
}

func CourierEarningFactor(c *entities.Courier) int32 {
	switch c.Type {
	case entities.FootCourier:
		return 2
	case entities.BikeCourier:
		return 3
	case entities.AutoCourier:
		return 4
	default:
		panic(fmt.Sprintf("unknown courier type: %q", c.Type))
	}
}

func CourierRatingFactor(c *entities.Courier) int32 {
	switch c.Type {
	case entities.FootCourier:
		return 3
	case entities.BikeCourier:
		return 2
	case entities.AutoCourier:
		return 1
	default:
		panic(fmt.Sprintf("unknown courier type: %q", c.Type))
	}
}
