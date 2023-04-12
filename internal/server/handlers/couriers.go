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
		ctx.JSON(http.StatusBadRequest, json.BadRequestResponse{})
		return nil
	}

	couriers, err := r.Actions.GetCouriers(ctx, queryParams.Limit, queryParams.Offset)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, couriers)
	return nil
}

func GetCourier(ctx *gin.Context, r *core.Repository) error {
	CourierId, err := strconv.ParseInt(ctx.Param("courier_id"), 10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, json.BadRequestResponse{})
		return nil
	}

	Courier, err := r.Actions.GetCourier(ctx, CourierId)
	if err != nil {
		panic(err)
	}

	if Courier == nil {
		ctx.JSON(http.StatusNotFound, json.NotFoundResponse{})
		return nil
	}

	ctx.JSON(http.StatusOK, Courier)
	return nil
}

func CreateCourier(ctx *gin.Context, r *core.Repository) error {
	var request struct {
		Couriers []struct {
			Type         entities.CourierType `json:"courier_type"`
			Regions      []int32              `json:"regions"`
			WorkingHours []types.Hour         `json:"working_hours"`
		} `json:"Couriers"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, json.BadRequestResponse{})
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

	couriers, err := r.Actions.CreateCouriers(ctx, createCouriers)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, couriers)
	return nil
}
