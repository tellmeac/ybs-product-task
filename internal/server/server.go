package server

import (
	"context"
	"github.com/gin-gonic/gin"
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

	app.Router.GET("/orders", app.handlerMapper(handlers.GetOrders))
	app.Router.POST("/orders", app.handlerMapper(handlers.CreateOrder))
	app.Router.GET("/orders/:order_id", app.handlerMapper(handlers.GetOrder))

	app.Router.GET("/couriers", app.handlerMapper(handlers.GetCouriers))
	app.Router.POST("/couriers", app.handlerMapper(handlers.CreateCourier))
	app.Router.GET("/couriers/:courier_id", app.handlerMapper(handlers.GetCourier))
}

func (app *App) handlerMapper(handler func(*gin.Context, *core.Repository) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_ = handler(ctx, app.Repository)
	}
}
