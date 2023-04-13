package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"yandex-team.ru/bstask/internal/core"
	"yandex-team.ru/bstask/internal/core/entities"
	"yandex-team.ru/bstask/internal/pkg/types"
	"yandex-team.ru/bstask/internal/pkg/web/json"
)

func GetCouriers(ctx *gin.Context, r *core.Repository) error {
	var queryParams = struct {
		Limit  uint64 `form:"limit"`
		Offset uint64 `form:"offset"`
	}{
		Limit:  1,
		Offset: 0,
	}

	if err := ctx.BindQuery(&queryParams); err != nil {
		ctx.JSON(http.StatusBadRequest, json.BadRequestResponse{Message: err.Error()})
		return nil
	}

	couriers, err := r.Actions.GetCouriers(ctx, queryParams.Limit, queryParams.Offset)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, gin.H{
		"couriers": couriers,
		"limit":    queryParams.Limit,
		"offset":   queryParams.Offset,
	})
	return nil
}

func GetCourier(ctx *gin.Context, r *core.Repository) error {
	courierId, err := strconv.ParseInt(ctx.Param("courier_id"), 10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, json.BadRequestResponse{Message: err.Error()})
		return nil
	}

	courier, err := r.Actions.GetCourier(ctx, courierId)
	if err != nil {
		return err
	}

	if courier == nil {
		ctx.JSON(http.StatusNotFound, json.NotFoundResponse{})
		return nil
	}

	ctx.JSON(http.StatusOK, courier)
	return nil
}

func CreateCourier(ctx *gin.Context, r *core.Repository) error {
	var request struct {
		Couriers []struct {
			Type         entities.CourierType `json:"courier_type"`
			Regions      []int32              `json:"regions"`
			WorkingHours []types.Interval     `json:"working_hours"`
		} `json:"couriers"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, json.BadRequestResponse{Message: err.Error()})
		return nil
	}

	createCouriers := make([]entities.Courier, 0, len(request.Couriers))
	for i := range request.Couriers {
		createCouriers = append(createCouriers, entities.Courier{
			Type:         request.Couriers[i].Type,
			Regions:      request.Couriers[i].Regions,
			WorkingHours: request.Couriers[i].WorkingHours,
		})
	}

	if err := r.Actions.ValidateCreateCouriers(createCouriers); err != nil {
		ctx.JSON(http.StatusBadRequest, json.BadRequestResponse{Message: err.Error()})
		return nil
	}

	couriers, err := r.Actions.CreateCouriers(ctx, createCouriers)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, gin.H{
		"couriers": couriers,
	})
	return nil
}
