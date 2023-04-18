package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"yandex-team.ru/bstask/internal/core"
	"yandex-team.ru/bstask/internal/pkg/web"
	"yandex-team.ru/bstask/internal/server/handlers"
)

type App struct {
	Server     web.Server
	Router     *gin.Engine
	Repository *core.Repository
}

func New(repository *core.Repository) *App {
	app := &App{
		Repository: repository,
	}

	app.initRoutes()
	app.Server = web.NewServer(repository.Config.Server, app.Router)

	return app
}

func (app *App) Start(ctx context.Context) error {
	return app.Server.Start(ctx)
}

func (app *App) initRoutes() {
	app.Router = gin.Default()

	app.Router.GET("/orders", app.mappedHandler(handlers.GetOrders))
	app.Router.POST("/orders", app.mappedHandler(handlers.CreateOrder))
	app.Router.GET("/orders/:order_id", app.mappedHandler(handlers.GetOrder))
	app.Router.POST("/orders/complete", app.mappedHandler(handlers.CompleteOrder))

	app.Router.GET("/couriers", app.mappedHandler(handlers.GetCouriers))
	app.Router.POST("/couriers", app.mappedHandler(handlers.CreateCourier))
	app.Router.GET("/couriers/:courier_id", app.mappedHandler(handlers.GetCourier))
	app.Router.GET("/couriers/meta-info/:courier_id", app.mappedHandler(handlers.GetCourierMetaInfo))
}

func (app *App) mappedHandler(handler func(*gin.Context, *core.Repository) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := handler(ctx, app.Repository); err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}
}
