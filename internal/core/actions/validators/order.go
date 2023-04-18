package validators

import (
	"errors"
	"yandex-team.ru/bstask/internal/core/entities"
	"yandex-team.ru/bstask/internal/pkg/types"
)

func ValidateOrder(order *entities.Order) error {
	if err := ValidateOrderDeliveryHours(order); err != nil {
		return err
	}

	return nil
}

func ValidateOrderDeliveryHours(order *entities.Order) error {
	if len(order.DeliveryHours) == 0 {
		return errors.New("must provide at least one delivery_hours value")
	}

	if collide := types.IntervalsCollide(order.DeliveryHours...); collide {
		return errors.New("invalid delivery hours with collisions")
	}

	return nil
}
