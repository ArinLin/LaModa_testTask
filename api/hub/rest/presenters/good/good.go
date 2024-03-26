package good

import (
	"time"

	"lamoda/internal/core"
	"lamoda/internal/store/good"
	"lamoda/pkg/web"
)

type Presenter struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Size      string     `json:"size"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func PresentGood(entity good.Entity) Presenter {
	return Presenter{
		ID:        entity.ID,
		Name:      entity.Name,
		Size:      entity.Size,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}

func (p *Presenter) Response(msg string) web.Response {
	return web.OKResponse(msg, *p, nil)
}

type ListPresenter []Presenter

func PresentList(entities []good.Entity) ListPresenter {
	pres := ListPresenter{}

	for _, entity := range entities {
		goodPresenter := Presenter{
			ID:        entity.ID,
			Name:      entity.Name,
			Size:      entity.Size,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
			DeletedAt: entity.DeletedAt,
		}
		pres = append(pres, goodPresenter)
	}

	return pres
}

func (p *ListPresenter) Response() web.Response {
	return web.OKResponse(core.GoodsReceivedCode, *p, nil)
}
