package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yandex-team.ru/bstask/internal/core"
)

// GetOrders ...
//
//	"/orders":
//	  get:
//	    tags:
//	      - order-controller
//	    operationId: getOrders
//	    parameters:
//	      - name: limit
//	        in: query
//	        description: Максимальное количество заказов в выдаче. Если параметр не передан,
//	          то значение по умолчанию равно 1.
//	        required: false
//	        schema:
//	          type: integer
//	          format: int32
//	        example: 10
//	      - name: offset
//	        in: query
//	        description: Количество заказов, которое нужно пропустить для отображения
//	          текущей страницы. Если параметр не передан, то значение по умолчанию равно
//	          0.
//	        required: false
//	        schema:
//	          type: integer
//	          format: int32
//	        example: 0
//	    responses:
//	      '200':
//	        description: ok
//	        content:
//	          application/json:
//	            schema:
//	              type: array
//	              items:
//	                "$ref": "#/components/schemas/OrderDto"
//	      '400':
//	        description: bad request
//	        content:
//	          application/json:
//	            schema:
//	              "$ref": "#/components/schemas/BadRequestResponse"
func GetOrders(ctx *gin.Context, r *core.Repository) error {
	q := struct {
		Limit  uint64 `form:"limit"`
		Offset uint64 `form:"offset"`
	}{
		Limit:  1,
		Offset: 0,
	}

	if err := ctx.BindQuery(&q); err != nil {
		// TODO: handle 400
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return err
	}

	orders, err := r.Storage.Orders.All(ctx, q.Limit, q.Offset)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, orders)

	return nil
}
