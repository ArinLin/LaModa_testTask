package hub

import (
	"context"
	"lamoda/internal/store/warehouse"
)

type (
	warehousesService interface {
		GetWarehouses(context.Context) ([]warehouse.Entity, error)
		GetWarehouseByID(context.Context, int) (warehouse.Entity, error)
		CreateWarehouse(context.Context, CreateWarehouseModel) (warehouse.Entity, error)
		UpdateWarehouse(context.Context, int, UpdateWarehouseModel) (warehouse.Entity, error)
		DeleteWarehouse(context.Context, int) error
	}

	CreateWarehouseModel struct {
		Name        string `json:"name" validate:"required,min=4,max=64"`
		IsAvailable bool   `json:"is_available" validate:"required"`
	}

	UpdateWarehouseModel struct {
		Name        *string `json:"name" validate:"required_without=IsAvailable,omitempty,min=4,max=64"`
		IsAvailable *bool   `json:"is_available,omitempty"`
	}
)

func (s *serviceImpl) GetWarehouses(ctx context.Context) ([]warehouse.Entity, error) {
	data, err := s.warehousesStore.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *serviceImpl) GetWarehouseByID(ctx context.Context, id int) (warehouse.Entity, error) {
	data, err := s.warehousesStore.GetByID(ctx, id)
	if err != nil {
		return warehouse.Entity{}, err
	}

	return data, nil
}

func (s *serviceImpl) CreateWarehouse(ctx context.Context, model CreateWarehouseModel) (warehouse.Entity, error) {
	entity := model.toCreateWarehouseEntity()

	data, err := s.warehousesStore.Create(ctx, entity)
	if err != nil {
		return warehouse.Entity{}, err
	}

	return data, nil
}

func (s *serviceImpl) UpdateWarehouse(ctx context.Context, id int, model UpdateWarehouseModel) (warehouse.Entity, error) {
	entity := model.toUpdateWarehouseEntity()

	data, err := s.warehousesStore.Update(ctx, id, entity)
	if err != nil {
		return warehouse.Entity{}, err
	}

	return data, nil
}

func (s *serviceImpl) DeleteWarehouse(ctx context.Context, id int) error {
	err := s.warehousesStore.Delete(ctx, id)

	return err
}

func (m CreateWarehouseModel) toCreateWarehouseEntity() warehouse.CreateEntity {
	return warehouse.CreateEntity{
		Name:        m.Name,
		IsAvailable: m.IsAvailable,
	}
}

func (m UpdateWarehouseModel) toUpdateWarehouseEntity() warehouse.UpdateEntity {
	return warehouse.UpdateEntity{
		Name:        m.Name,
		IsAvailable: m.IsAvailable,
	}
}
