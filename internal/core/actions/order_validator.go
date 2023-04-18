package actions

import (
	"errors"
	"yandex-team.ru/bstask/internal/core/entities"
)

func ValidateOrder(order *entities.Order) error {
	if err := validateHours(order); err != nil {
		return err
	}

	return nil
}

func validateHours(order *entities.Order) error {
	if len(order.DeliveryHours) == 0 {
		return errors.New("must provide at least one delivery_hours value")
	}

	return nil
}
