package stock

import (
	"context"
	"database/sql"

	"lamoda/internal/core"
	"lamoda/pkg/logger"
	"lamoda/pkg/postgres"

	"github.com/jmoiron/sqlx"
)

type (
	Store interface {
		Create(context.Context, CreateEntity) (Entity, error)
		GetByGoodID(context.Context, int) ([]Entity, error)
		GetGoodsAmountByWarehouseID(context.Context, int) ([]AmountEntity, error)
		Reserve(context.Context, []ChangeAmountEntity) error
		Release(context.Context, []ChangeAmountEntity) error
	}

	storeImpl struct {
		client *postgres.Client
	}

	ChangeAmountEntity struct {
		GoodID int
		Amount int
	}

	CreateEntity struct {
		GoodID      int
		WarehouseID int
		Amount      int
	}

	Entity struct {
		GoodID      int `db:"good_id"`
		WarehouseID int `db:"warehouse_id"`
		Amount      int `db:"amount"`
		Reserved    int `db:"reserved"`
	}
	// какой товар в каком количестве лежит в определенном складе
	AmountEntity struct {
		GoodID   int    `db:"good_id"`
		GoodName string `db:"good_name"`
		Amount   int    `db:"amount"`
		Reserved int    `db:"reserved"`
	}
	// перечисление складов, с доступными товарами и зарезервированными для конкретного GoodID
	WarehouseStock struct {
		WarehouseID int `db:"warehouse_id"`
		Amount      int `db:"amount"`
		Reserved    int `db:"reserved"`
	}
)

func New(client *postgres.Client) Store {
	return &storeImpl{
		client: client,
	}
}

func (s *storeImpl) Create(ctx context.Context, entity CreateEntity) (Entity, error) {
	var newStock Entity

	err := s.client.QueryRowxContext(
		ctx,
		`
            INSERT INTO stocks (good_id, warehouse_id, amount)
            VALUES ($1, $2, $3)
			RETURNING *
		`,
		entity.GoodID,
		entity.WarehouseID,
		entity.Amount,
	).StructScan(&newStock)
	if err != nil {
		logger.Log.Error("create stock",
			"error", err)
	}

	return newStock, err
}

func (s *storeImpl) GetByGoodID(ctx context.Context, id int) ([]Entity, error) {
	var stock []Entity

	err := s.client.SelectContext(
		ctx,
		&stock,
		`
			SELECT *
			FROM stocks
			WHERE good_id = $1
		`,
		id)
	if err != nil {
		logger.Log.Error("get stock by id",
			"error", err)
	}

	return stock, err
}

func (s *storeImpl) GetGoodsAmountByWarehouseID(ctx context.Context, warehouseID int) ([]AmountEntity, error) {
	var goodsAmount []AmountEntity

	err := s.client.SelectContext(
		ctx,
		&goodsAmount,
		`
			SELECT g.id as good_id, g.name as good_name, amount, reserved
			FROM stocks s
			JOIN goods g ON g.id = s.good_id
			WHERE s.warehouse_id = $1
			ORDER BY g.id
		`,
		warehouseID)
	if err != nil {
		logger.Log.Error("get all goods for warehouse",
			"warehouse id", warehouseID,
			"error", err)
	}

	if len(goodsAmount) == 0 {
		return nil, core.ErrStockNotFound
	}

	return goodsAmount, err
}

func (s *storeImpl) Reserve(ctx context.Context, goods []ChangeAmountEntity) error {
	tx, err := s.client.BeginTxx(ctx, nil)
	if err != nil {
		logger.Log.Error("start transaction",
			"error", err)
		return err
	}

	if err = reserveAmountInTx(ctx, goods, tx); err != nil {
		logger.Log.Error("reserve amount in tx",
			"error", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		logger.Log.Error("commit transaction",
			"error", err)
	}

	return err
}

func (s *storeImpl) Release(ctx context.Context, goods []ChangeAmountEntity) error {
	tx, err := s.client.BeginTxx(ctx, nil)
	if err != nil {
		logger.Log.Error("start transaction",
			"error", err)
		return err
	}

	if err := releaseAmountInTx(ctx, goods, tx); err != nil {
		logger.Log.Error("release amount in tx",
			"error", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		logger.Log.Error("commit transaction",
			"error", err)
	}

	return err
}

// в каких складах в каком количестве есть этот товар
func getWarehouseStockByGoodIDInTx(ctx context.Context, goodID int, orderBy string, tx *sqlx.Tx) ([]WarehouseStock, error) {
	var stocks []WarehouseStock

	err := tx.SelectContext(
		ctx,
		&stocks,
		`
			SELECT warehouse_id, amount, reserved
			FROM stocks
			JOIN warehouses w ON w.id = warehouse_id
			WHERE good_id = $1
			AND w.is_available = true
			ORDER BY $2 DESC
		`,
		goodID,
		orderBy,
	)

	return stocks, err
}

func reserveAmountInTx(ctx context.Context, goods []ChangeAmountEntity, tx *sqlx.Tx) error {
	for _, good := range goods {
		var totalReserved int
		stocks, err := getWarehouseStockByGoodIDInTx(ctx, good.GoodID, "amount", tx)
		if err != nil {
			logger.Log.Error("get goods amount for warehouses",
				"error", err)
			return err
		}
		if len(stocks) == 0 {
			return sql.ErrNoRows
		}

		for _, stock := range stocks {
			if totalReserved >= good.Amount {
				break
			}

			reserveAmount := min(stock.Amount, good.Amount-totalReserved)
			if reserveAmount <= 0 {
				break
			}

			_, err := tx.ExecContext(ctx,
				`
				UPDATE stocks
				SET amount = amount - $1, reserved = reserved + $1
				WHERE good_id = $2
				AND warehouse_id = $3
			`,
				reserveAmount,
				good.GoodID,
				stock.WarehouseID)
			if err != nil {
				logger.Log.Error("reserving stock",
					"error", err)
				return err
			}

			totalReserved += reserveAmount
		}

		if totalReserved < good.Amount {
			return core.ErrNotEnoughAmount
		}
	}

	return nil
}

func releaseAmountInTx(ctx context.Context, goods []ChangeAmountEntity, tx *sqlx.Tx) error {
	for _, good := range goods {
		var totalReleased int
		stocks, err := getWarehouseStockByGoodIDInTx(ctx, good.GoodID, "reserved", tx)
		if err != nil {
			logger.Log.Error("get goods amount for warehouses",
				"error", err)
			return err
		}
		if len(stocks) == 0 {
			return sql.ErrNoRows
		}

		for _, stock := range stocks {
			if totalReleased >= good.Amount {
				break
			}

			releaseAmount := min(stock.Reserved, good.Amount-totalReleased)
			if releaseAmount <= 0 {
				break
			}

			_, err := tx.ExecContext(ctx,
				`
				UPDATE stocks
				SET amount = amount + $1, reserved = reserved - $1
				WHERE good_id = $2
				AND warehouse_id = $3
			`,
				releaseAmount,
				good.GoodID,
				stock.WarehouseID)
			if err != nil {
				logger.Log.Error("releasing stock",
					"error", err)
				return err
			}

			totalReleased += releaseAmount
		}

		if totalReleased < good.Amount {
			return core.ErrNotEnoughReserve
		}
	}

	return nil
}
