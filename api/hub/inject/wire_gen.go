// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package inject

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/urfave/cli/v2"
	"lamoda/api/hub"
	"lamoda/api/hub/rest/handlers"
	"lamoda/internal/service/hub"
	"lamoda/internal/store/good"
	"lamoda/internal/store/stock"
	"lamoda/internal/store/warehouse"
	"lamoda/pkg/postgres"
)

// Injectors from wire.go:

func InitializeApplication(c *cli.Context, appCtx context.Context) (api.Container, error) {
	echoEcho := echo.New()
	config := providePostgresConfig(c)
	client, err := postgres.NewClient(config)
	if err != nil {
		return api.Container{}, err
	}
	store := good.New(client)
	stockStore := stock.New(client)
	warehouseStore := warehouse.New(client)
	service := hub.New(store, stockStore, warehouseStore)
	resolver := handlers.NewResolver(echoEcho, service)
	container := api.NewContainer(echoEcho, resolver)
	return container, nil
}
