package hub

import (
	"context"

	"lamoda/internal/store/good"
)

type (
	goodsService interface {
		GetGoods(context.Context) ([]good.Entity, error)
		GetGoodByID(context.Context, int) (good.Entity, error)
		CreateGood(context.Context, CreateGoodModel) (good.Entity, error)
		UpdateGood(context.Context, int, UpdateGoodModel) (good.Entity, error)
		DeleteGood(context.Context, int) error
	}

	CreateGoodModel struct {
		Name string `json:"name" validate:"required,min=4,max=127"`
		Size string `json:"size" validate:"required,oneof=s m l"`
	}

	UpdateGoodModel struct {
		Name *string `json:"name" validate:"required_without=Size,omitempty,min=4,max=127"`
		Size *string `json:"size" validate:"omitempty,oneof=s m l"`
	}
)

func (s *serviceImpl) GetGoods(ctx context.Context) ([]good.Entity, error) {
	data, err := s.goodsStore.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *serviceImpl) GetGoodByID(ctx context.Context, id int) (good.Entity, error) {
	data, err := s.goodsStore.GetByID(ctx, id)
	if err != nil {
		return good.Entity{}, err
	}

	return data, nil
}

func (s *serviceImpl) CreateGood(ctx context.Context, model CreateGoodModel) (good.Entity, error) {
	entity := model.toCreateGoodEntity()

	data, err := s.goodsStore.Create(ctx, entity)
	if err != nil {
		return good.Entity{}, err
	}

	return data, nil
}

func (s *serviceImpl) UpdateGood(ctx context.Context, id int, model UpdateGoodModel) (good.Entity, error) {
	entity := model.toUpdateGoodEntity()

	data, err := s.goodsStore.Update(ctx, id, entity)
	if err != nil {
		return good.Entity{}, err
	}

	return data, nil
}

func (s *serviceImpl) DeleteGood(ctx context.Context, id int) error {
	err := s.goodsStore.Delete(ctx, id)

	return err
}

func (m CreateGoodModel) toCreateGoodEntity() good.CreateEntity {
	return good.CreateEntity{
		Name: m.Name,
		Size: m.Size,
	}
}

func (m UpdateGoodModel) toUpdateGoodEntity() good.UpdateEntity {
	return good.UpdateEntity{
		Name: m.Name,
		Size: m.Size,
	}
}
