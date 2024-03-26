package hub

import (
	"context"
	"lamoda/internal/store/stock"
)

type (
	stocksService interface {
		GetStockByWarehouseID(context.Context, int) ([]stock.AmountEntity, error)
		ReserveStocks(context.Context, ChangeStocksAmountModel) error
		ReleaseStocks(context.Context, ChangeStocksAmountModel) error
	}

	ChangeStocksAmountModel struct {
		Goods []Item `json:"goods" validate:"required"`
	}

	Item struct {
		ID     int `json:"id" validate:"required"`
		Amount int `json:"amount" validate:"required"`
	}
)

func (s *serviceImpl) GetStockByWarehouseID(ctx context.Context, warehouseID int) ([]stock.AmountEntity, error) {
	data, err := s.stocksStore.GetGoodsAmountByWarehouseID(ctx, warehouseID)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *serviceImpl) ReserveStocks(ctx context.Context, model ChangeStocksAmountModel) error {
	stocks := model.toChangeAmountEntity()

	return s.stocksStore.Reserve(ctx, stocks)
}

func (s *serviceImpl) ReleaseStocks(ctx context.Context, model ChangeStocksAmountModel) error {
	stocks := model.toChangeAmountEntity()

	return s.stocksStore.Release(ctx, stocks)
}

func (m ChangeStocksAmountModel) toChangeAmountEntity() []stock.ChangeAmountEntity {
	stocks := make([]stock.ChangeAmountEntity, len(m.Goods))
	for i, good := range m.Goods {
		stocks[i] = stock.ChangeAmountEntity{GoodID: good.ID, Amount: good.Amount}
	}

	return stocks
}
