package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
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
		zap.L().Debug("Get couriers", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, json.BadRequestResponse{})
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
		zap.L().Debug("Get courier", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, json.BadRequestResponse{})
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
		zap.L().Debug("Create courier", zap.Error(err))
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

	if err := r.Actions.ValidateCreateCouriers(createCouriers); err != nil {
		zap.L().Debug("Validate create couriers request", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, json.BadRequestResponse{})
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

func GetCourierMetaInfo(ctx *gin.Context, r *core.Repository) error {
	courierId, err := strconv.ParseInt(ctx.Param("courier_id"), 10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, json.BadRequestResponse{})
		return nil
	}

	var queryParams struct {
		StartDate time.Time `form:"startDate" time_format:"2006-01-02"`
		EndDate   time.Time `form:"endDate" time_format:"2006-01-02" binding:"gtfield=StartDate"`
	}
	if err := ctx.BindQuery(&queryParams); err != nil {
		zap.L().Debug("Get courier meta info", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, json.BadRequestResponse{})
		return nil
	}

	result, err := r.Actions.GetCourierMetaInfo(ctx, courierId, queryParams.StartDate, queryParams.EndDate)
	if err != nil {
		return err
	}

	// NOTE: 404 ответ не предусмотрен спецификацией.
	if result == nil {
		zap.L().Error("Courier meta: requested courier was not found")
		return errors.New("requested courier was not found")
	}

	ctx.JSON(http.StatusOK, result)
	return nil
}
