package validators

import (
	"errors"
	"fmt"
	"github.com/samber/lo"
	"yandex-team.ru/bstask/internal/core/entities"
	"yandex-team.ru/bstask/internal/pkg/types"
)

func ValidateCourier(courier *entities.Courier) error {
	if err := validateCourierType(courier); err != nil {
		return err
	}

	if err := validateCourierRegions(courier); err != nil {
		return err
	}

	if err := validateCourierWorkingHours(courier); err != nil {
		return err
	}

	return nil
}

func validateCourierType(courier *entities.Courier) error {
	_, valid := lo.Find(entities.CourierTypes, func(t entities.CourierType) bool {
		return courier.Type == t
	})

	if !valid {
		return fmt.Errorf("invalid courier type: %s", courier.Type)
	}
	return nil
}

func validateCourierRegions(courier *entities.Courier) error {
	if len(courier.Regions) == 0 {
		return errors.New("must provide at least one available region")
	}

	if len(lo.Uniq(courier.Regions)) != len(courier.Regions) {
		return errors.New("regions must be unique")
	}

	return nil
}

func validateCourierWorkingHours(courier *entities.Courier) error {
	if len(courier.WorkingHours) == 0 {
		return errors.New("must provide at least one delivery_hours value")
	}

	if collide := types.IntervalsCollide(courier.WorkingHours...); collide {
		return errors.New("invalid delivery hours with collisions")
	}

	return nil
}
