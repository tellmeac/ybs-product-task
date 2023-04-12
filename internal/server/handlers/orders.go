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

func GetOrders(ctx *gin.Context, r *core.Repository) error {
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

	orders, err := r.Actions.GetOrders(ctx, queryParams.Limit, queryParams.Offset)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, orders)
	return nil
}

func GetOrder(ctx *gin.Context, r *core.Repository) error {
	orderId, err := strconv.ParseInt(ctx.Param("order_id"), 10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, json.BadRequestResponse{})
		return nil
	}

	order, err := r.Actions.GetOrder(ctx, orderId)
	if err != nil {
		panic(err)
	}

	if order == nil {
		ctx.JSON(http.StatusNotFound, json.NotFoundResponse{})
		return nil
	}

	ctx.JSON(http.StatusOK, order)
	return nil
}

func CreateOrder(ctx *gin.Context, r *core.Repository) error {
	var request struct {
		Orders []struct {
			Weight        float64      `json:"weight"`
			Region        int32        `json:"regions"` // NOTE: One region for order.
			DeliveryHours []types.Hour `json:"delivery_hours"`
			Cost          int32        `json:"cost"`
		} `json:"orders"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, json.BadRequestResponse{})
		return nil
	}

	createOrders := make([]entities.Order, 0, len(request.Orders))
	for i := range request.Orders {
		createOrders = append(createOrders, entities.Order{
			Weight:        request.Orders[i].Weight,
			Region:        request.Orders[i].Region,
			DeliveryHours: request.Orders[i].DeliveryHours,
			Cost:          request.Orders[i].Cost,
		})
	}

	orders, err := r.Actions.CreateOrders(ctx, createOrders)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, orders)
	return nil
}
