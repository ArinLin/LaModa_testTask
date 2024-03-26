package stock

import (
	"lamoda/internal/core"
	"lamoda/internal/store/stock"
	"lamoda/pkg/web"
)

type Presenter struct {
	GoodID   int    `json:"good_id"`
	GoodName string `json:"good_name"`
	Amount   int    `json:"amount"`
	Reserved int    `json:"reserved"`
}

type ListPresenter []Presenter

func PresentList(entities []stock.AmountEntity) ListPresenter {
	pres := ListPresenter{}

	for _, entity := range entities {
		stockPresenter := Presenter{
			GoodID:   entity.GoodID,
			GoodName: entity.GoodName,
			Amount:   entity.Amount,
			Reserved: entity.Reserved,
		}
		pres = append(pres, stockPresenter)
	}

	return pres
}

func (p *ListPresenter) Response() web.Response {
	return web.OKResponse(core.StocksReceivedCode, *p, nil)
}
