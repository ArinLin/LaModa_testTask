package warehouse

import (
	"time"

	"lamoda/internal/core"
	"lamoda/internal/store/warehouse"
	"lamoda/pkg/web"
)

type Presenter struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	IsAvailable bool       `json:"is_available"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

func PresentWarehouse(entity warehouse.Entity) Presenter {
	return Presenter{
		ID:          entity.ID,
		Name:        entity.Name,
		IsAvailable: entity.IsAvailable,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		DeletedAt:   entity.DeletedAt,
	}
}

func (p *Presenter) Response(msg string) web.Response {
	return web.OKResponse(msg, *p, nil)
}

type ListPresenter []Presenter

func PresentList(entities []warehouse.Entity) ListPresenter {
	pres := ListPresenter{}

	for _, entity := range entities {
		warehousePresenter := Presenter{
			ID:          entity.ID,
			Name:        entity.Name,
			IsAvailable: entity.IsAvailable,
			CreatedAt:   entity.CreatedAt,
			UpdatedAt:   entity.UpdatedAt,
			DeletedAt:   entity.DeletedAt,
		}
		pres = append(pres, warehousePresenter)
	}

	return pres
}

func (p *ListPresenter) Response() web.Response {
	return web.OKResponse(core.WarehousesReceivedCode, *p, nil)
}
