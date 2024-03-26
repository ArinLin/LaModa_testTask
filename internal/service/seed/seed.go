package seed

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"

	"lamoda/internal/store/good"
	"lamoda/internal/store/stock"
	"lamoda/internal/store/warehouse"
	"lamoda/pkg/logger"
	"lamoda/pkg/postgres"

	"github.com/brianvoe/gofakeit/v6"
)

type Counter struct {
	Goods      int
	Stocks     int
	Warehouses int
}

func Seed(ctx context.Context, postgresClient *postgres.Client, counter Counter) error {
	goodsID, err := seedGoods(ctx, postgresClient, counter.Goods)
	if err != nil {
		return err
	}

	warehousesID, err := seedWarehouses(ctx, postgresClient, counter.Warehouses)
	if err != nil {
		return err
	}

	return seedStocks(ctx, postgresClient, counter.Stocks, goodsID, warehousesID)
}

func seedGoods(ctx context.Context, db *postgres.Client, count int) ([]int, error) {
	store := good.New(db)

	goods := make([]int, 0, count)
	for i := 0; i < count; i++ {
		goodName := fmt.Sprintf("%s %s %s", gofakeit.Color(), gofakeit.MinecraftArmorPart(), gofakeit.BeerName())
		goodSize := []string{"s", "m", "l"}
		createdGood, err := store.Create(ctx, good.CreateEntity{Name: goodName, Size: goodSize[i%len(goodSize)]})
		if err != nil {
			logger.Log.Error("generating good",
				"error", err.Error())
			return nil, err
		}
		goods = append(goods, createdGood.ID)

	}
	logger.Log.Debug("goods created successfully")

	return goods, nil
}

func seedWarehouses(ctx context.Context, db *postgres.Client, count int) ([]int, error) {
	store := warehouse.New(db)

	warehouses := make([]int, 0, count)
	for i := 0; i < count; i++ {
		warehouseName := fmt.Sprintf("%s %s", gofakeit.City(), gofakeit.StreetName())
		createdWarehouse, err := store.Create(ctx, warehouse.CreateEntity{Name: warehouseName, IsAvailable: rand.Intn(2) == 1})
		if err != nil {
			logger.Log.Error("generating warehouse",
				"error", err.Error())
			return nil, err
		}
		warehouses = append(warehouses, createdWarehouse.ID)
	}
	logger.Log.Debug("warehouses created successfully")

	return warehouses, nil
}

func seedStocks(ctx context.Context, db *postgres.Client, count int, goodsID, warehousesID []int) error {
	store := stock.New(db)
	totalCreated := 0

	for _, warehouseID := range warehousesID {
		for _, goodID := range goodsID {
			if totalCreated >= count {
				return nil
			}

			_, err := store.Create(ctx, stock.CreateEntity{
				GoodID:      goodID,
				WarehouseID: warehouseID,
				Amount:      rand.Intn(10000),
			})
			if err != nil {
				logger.Log.Error("create stock",
					"error", err.Error())
				return err
			}
			totalCreated++
		}
	}

	logger.Log.Debug("stocks created successfully")
	return nil
}

func CheckIfDatabaseIsEmpty(ctx context.Context, db *postgres.Client) (bool, error) {
	var test int
	err := db.QueryRowContext(ctx, "SELECT 1 FROM stocks LIMIT 1").Scan(&test)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}

	return false, nil
}
