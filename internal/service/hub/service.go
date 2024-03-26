package hub

import (
	"lamoda/internal/store/good"
	"lamoda/internal/store/stock"
	"lamoda/internal/store/warehouse"
)

type (
	Service interface {
		goodsService
		stocksService
		warehousesService
	}

	serviceImpl struct {
		goodsStore      good.Store
		stocksStore     stock.Store
		warehousesStore warehouse.Store
	}
)

func New(
	goods good.Store,
	stocks stock.Store,
	warehouses warehouse.Store,
) Service {
	return &serviceImpl{
		goodsStore:      goods,
		stocksStore:     stocks,
		warehousesStore: warehouses,
	}
}
