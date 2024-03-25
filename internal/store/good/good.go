package good

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
		Name string
		Size string
	}

	UpdateEntity struct {
		Name *string
		Size *string
	}

	Entity struct {
		ID        int        `db:"id"`
		Name      string     `db:"name"`
		Size      string     `db:"size"`
		CreatedAt time.Time  `db:"created_at"`
		UpdatedAt time.Time  `db:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at"`
	}
)

func New(client *postgres.Client) Store {
	return &storeImpl{
		client: client,
	}
}

func (s *storeImpl) Create(ctx context.Context, entity CreateEntity) (Entity, error) {
	var newGood Entity

	err := s.client.QueryRowxContext(
		ctx,
		`
			INSERT INTO goods (name, size)
			VALUES ($1, $2)
			RETURNING *
		`,
		entity.Name,
		entity.Size,
	).StructScan(&newGood)
	if err != nil {
		logger.Log.Error("create new good",
			"error", err)
	}

	return newGood, err
}

func (s *storeImpl) GetByID(ctx context.Context, id int) (Entity, error) {
	var good Entity

	err := s.client.GetContext(
		ctx,
		&good,
		`
			SELECT *
			FROM goods
			WHERE id = $1
		`,
		id)
	if err != nil {
		logger.Log.Error("get good by id",
			"error", err)
	}

	return good, err
}

func (s *storeImpl) GetAll(ctx context.Context) ([]Entity, error) {
	var goods []Entity

	err := s.client.SelectContext(
		ctx,
		&goods,
		`
			SELECT *
			FROM goods
			WHERE deleted_at IS NULL
			ORDER BY id
		`)
	if err != nil {
		logger.Log.Error("get all goods",
			"error", err)
	}

	return goods, err
}

func (s *storeImpl) Update(ctx context.Context, id int, entity UpdateEntity) (Entity, error) {
	tx, err := s.client.BeginTxx(ctx, nil)
	if err != nil {
		logger.Log.Error("start transaction",
			"error", err)
		return Entity{}, err
	}

	good, err := getGoodByIDInTx(ctx, id, tx)
	if err != nil {
		logger.Log.Error("get good by id",
			"error", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return Entity{}, rollbackErr
		}
		return Entity{}, err
	}
	if entity.Name != nil {
		good.Name = *entity.Name
	}
	if entity.Size != nil {
		good.Size = *entity.Size
	}

	var updatedGood Entity
	err = tx.QueryRowxContext(
		ctx,
		`
			UPDATE goods
			SET name = $1, size = $2, updated_at = NOW()
			WHERE id = $3
			RETURNING *
		`,
		good.Name,
		good.Size,
		id,
	).StructScan(&updatedGood)
	if err != nil {
		logger.Log.Error("update good in tx",
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

	return updatedGood, err
}

func (s *storeImpl) Delete(ctx context.Context, id int) error {
	res, err := s.client.ExecContext(
		ctx,
		`
			UPDATE goods
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

func getGoodByIDInTx(ctx context.Context, id int, tx *sqlx.Tx) (Entity, error) {
	var good Entity
	err := tx.GetContext(
		ctx,
		&good,
		`
			SELECT *
			FROM goods
			WHERE id = $1
		`,
		id)

	return good, err
}
