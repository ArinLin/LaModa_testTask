package handlers

import (
	"time"

	"lamoda/internal/service/hub"
	"lamoda/pkg/logger"

	"github.com/labstack/echo/v4"
)

type Resolver struct {
	server  *echo.Echo
	service hub.Service
}

const (
	servicePrefix = "/hub"
	APIPrefix     = servicePrefix + "/api"
	APIVersion    = "v1"
	pathPrefix    = APIPrefix + "/" + APIVersion
)

func NewResolver(e *echo.Echo, service hub.Service) *Resolver {
	resolver := &Resolver{
		server:  e,
		service: service,
	}

	resolver.server.Use(LoggerMiddleware())
	resolver.initRoutes()

	return resolver
}

func (r *Resolver) initRoutes() {
	// goods
	r.server.GET(pathPrefix+"/good/:id", r.getGoodByID)
	r.server.GET(pathPrefix+"/goods", r.getGoods)
	r.server.POST(pathPrefix+"/good", r.createGood)
	r.server.PUT(pathPrefix+"/good/:id", r.updateGood)
	r.server.DELETE(pathPrefix+"/good/:id", r.deleteGood)

	// stocks
	r.server.GET(pathPrefix+"/stocks/warehouse/:id", r.getStockByWarehouseID)
	r.server.PUT(pathPrefix+"/stocks/reserve", r.reserveStock)
	r.server.PUT(pathPrefix+"/stocks/release", r.releaseStock)

	// warehouses
	r.server.GET(pathPrefix+"/warehouse/:id", r.getWarehouseByID)
	r.server.GET(pathPrefix+"/warehouses", r.getWarehouses)
	r.server.POST(pathPrefix+"/warehouse", r.createWarehouse)
	r.server.PUT(pathPrefix+"/warehouse/:id", r.updateWarehouse)
	r.server.DELETE(pathPrefix+"/warehouse/:id", r.deleteWarehouse)
}

func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			duration := time.Since(start)

			logger.Log.Info("Request details",
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", c.Response().Status,
				"duration", duration)

			return err
		}
	}
}
