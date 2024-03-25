package warehouse

import (
	"context"
	"database/sql"
	"time"

	"lamoda/pkg/logger"
	"lamoda/pkg/postgres"

	"github.com/jmoiron/sqlx"
)

type (
	Store interface {
		Create(context.Context, CreateEntity) (Entity, error)
		GetByID(context.Context, int) (Entity, error)
		GetAll(context.Context) ([]Entity, error)
		Update(context.Context, int, UpdateEntity) (Entity, error)
		Delete(context.Context, int) error
	}

	storeImpl struct {
		client *postgres.Client
	}

	CreateEntity struct {
		Name        string
		IsAvailable bool
	}

	UpdateEntity struct {
		Name        *string
		IsAvailable *bool
	}

	Entity struct {
		ID          int        `db:"id"`
		Name        string     `db:"name"`
		IsAvailable bool       `db:"is_available"`
		CreatedAt   time.Time  `db:"created_at"`
		UpdatedAt   time.Time  `db:"updated_at"`
		DeletedAt   *time.Time `db:"deleted_at"`
	}
)

func New(client *postgres.Client) Store {
	return &storeImpl{
		client: client,
	}
}

func (s *storeImpl) Create(ctx context.Context, entity CreateEntity) (Entity, error) {
	var newWarehouse Entity

	err := s.client.QueryRowxContext(
		ctx,
		`
			INSERT INTO warehouses (name, is_available)
			VALUES ($1, $2)
			RETURNING *
		`,
		entity.Name,
		entity.IsAvailable,
	).StructScan(&newWarehouse)
	if err != nil {
		logger.Log.Error("create new warehouse",
			"error", err)
	}

	return newWarehouse, err
}

func (s *storeImpl) GetByID(ctx context.Context, id int) (Entity, error) {
	var warehouse Entity

	err := s.client.GetContext(
		ctx,
		&warehouse,
		`
			SELECT *
			FROM warehouses
			WHERE id = $1
		`,
		id)
	if err != nil {
		logger.Log.Error("get warehouse by id",
			"error", err)
	}

	return warehouse, err
}

func (s *storeImpl) GetAll(ctx context.Context) ([]Entity, error) {
	var warehouses []Entity

	err := s.client.SelectContext(
		ctx,
		&warehouses,
		`
			SELECT *
			FROM warehouses
			WHERE deleted_at IS NULL
			ORDER BY id
		`)
	if err != nil {
		logger.Log.Error("get all warehouses",
			"error", err)
	}

	return warehouses, err
}

func (s *storeImpl) Update(ctx context.Context, id int, entity UpdateEntity) (Entity, error) {
	tx, err := s.client.BeginTxx(ctx, nil)
	if err != nil {
		logger.Log.Error("start transaction",
			"error", err)
		return Entity{}, err
	}

	good, err := getWarehouseByIDInTx(ctx, id, tx)
	if err != nil {
		logger.Log.Error("get warehouse by id",
			"error", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return Entity{}, rollbackErr
		}
		return Entity{}, err
	}
	if entity.Name != nil {
		good.Name = *entity.Name
	}
	if entity.IsAvailable != nil {
		good.IsAvailable = *entity.IsAvailable
	}

	var updatedWarehouse Entity
	err = tx.QueryRowxContext(
		ctx,
		`
			UPDATE warehouses
			SET name = $1, is_available = $2, updated_at = NOW()
			WHERE id = $3
			RETURNING *
		`,
		good.Name,
		good.IsAvailable,
		id,
	).StructScan(&updatedWarehouse)
	if err != nil {
		logger.Log.Error("update warehouse in tx",
			"error", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return Entity{}, rollbackErr
		}
		return Entity{}, err
	}

	if err := tx.Commit(); err != nil {
		logger.Log.Error("commit transaction",
			"error", err)
	}

	return updatedWarehouse, err
}

func (s *storeImpl) Delete(ctx context.Context, id int) error {
	res, err := s.client.ExecContext(
		ctx,
		`
			UPDATE warehouses
			SET updated_at = NOW(), deleted_at = NOW()
			WHERE id = $1
		`,
		id,
	)
	if err != nil {
		logger.Log.Error("delete good",
			"error", err)
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func getWarehouseByIDInTx(ctx context.Context, id int, tx *sqlx.Tx) (Entity, error) {
	var warehouse Entity
	err := tx.GetContext(
		ctx,
		&warehouse,
		`
			SELECT *
			FROM warehouses
			WHERE id = $1
		`,
		id)

	return warehouse, err
}
