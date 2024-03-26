package stock

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"lamoda/internal/core"
	"lamoda/internal/tests"
	"lamoda/pkg/logger"
	"lamoda/pkg/postgres"
)

func TestCreate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresContainer, mappedPort := tests.CreateInventoryHubTestPostgresContainer(ctx)
	defer tests.TerminateInventoryHubTestContainer(ctx, postgresContainer)

	postgresClient, err := postgres.NewClient(postgres.Config{
		User:       tests.PostgresInventoryHubTestUser,
		Password:   tests.PostgresInventoryHubTestPassword,
		Host:       tests.PostgresInventoryHubTestHost + ":" + mappedPort,
		DBName:     tests.PostgresInventoryHubTestDB,
		DisableTLS: true,
	})
	if err != nil {
		t.Errorf("error with starting postgres client: %s", err.Error())
	}
	defer postgresClient.Close()

	store := New(postgresClient)

	if err := seedTestData(ctx, false, postgresClient); err != nil {
		t.Fatalf("error with seeding test data: %s", err)
	}

	for _, test := range []struct {
		Name    string
		Data    CreateEntity
		WantErr bool
		Err     string
	}{
		{
			Name: "Successful create stock",
			Data: CreateEntity{
				GoodID:      1,
				WarehouseID: 1,
				Amount:      100,
			},
		},
		{
			Name: "Creating a stock with non-existing good id",
			Data: CreateEntity{
				GoodID:      1000,
				WarehouseID: 1,
				Amount:      100,
			},
			WantErr: true,
			Err:     "ERROR: insert or update on table \"stocks\" violates foreign key constraint \"stocks_good_id_fkey\" (SQLSTATE 23503)",
		},
		{
			Name: "Creating a stock with non-existing warehouse id",
			Data: CreateEntity{
				GoodID:      1,
				WarehouseID: 1000,
				Amount:      100,
			},
			WantErr: true,
			Err:     "ERROR: insert or update on table \"stocks\" violates foreign key constraint \"stocks_warehouse_id_fkey\" (SQLSTATE 23503)",
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			entity, err := store.Create(ctx, test.Data)
			if err != nil {
				if !test.WantErr {
					t.Errorf("unexpected error: %s", err.Error())
				}
				if err.Error() != test.Err {
					t.Errorf("unexpected error. Expected %q but got %q", test.Err, err.Error())
				}
				return
			}
			if test.WantErr {
				t.Errorf("expected error but nothing got")
			}
			logger.Log.Info("Created stock", "entity", entity)

			if entity.GoodID != test.Data.GoodID {
				t.Errorf("wrong good id. Expected %d but got %d", test.Data.GoodID, entity)
			}
			if entity.WarehouseID != test.Data.WarehouseID {
				t.Errorf("wrong warehouse id. Expected %d but got %d", test.Data.WarehouseID, entity)
			}
			if entity.Amount != test.Data.Amount {
				t.Errorf("wrong amount. Expected %d but got %d", test.Data.Amount, entity)
			}
			if entity.Reserved != 0 {
				t.Errorf("reserved amount should be 0 but got %d", entity.Reserved)
			}
		})
	}
}

func TestGetByGoodID(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresContainer, mappedPort := tests.CreateInventoryHubTestPostgresContainer(ctx)
	defer tests.TerminateInventoryHubTestContainer(ctx, postgresContainer)

	postgresClient, err := postgres.NewClient(postgres.Config{
		User:       tests.PostgresInventoryHubTestUser,
		Password:   tests.PostgresInventoryHubTestPassword,
		Host:       tests.PostgresInventoryHubTestHost + ":" + mappedPort,
		DBName:     tests.PostgresInventoryHubTestDB,
		DisableTLS: true,
	})
	if err != nil {
		t.Errorf("error with starting postgres client: %s", err.Error())
	}
	defer postgresClient.Close()

	store := New(postgresClient)

	if err := seedTestData(ctx, true, postgresClient); err != nil {
		t.Fatalf("error with seeding test data: %s", err)
	}

	for _, test := range []struct {
		Name    string
		ID      int
		WantErr bool
		Err     string
	}{
		{
			Name: "Successful get a stocks by good id",
			ID:   1,
		},
		{
			Name: "Get non-existent stock",
			ID:   1000,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			entities, err := store.GetByGoodID(ctx, test.ID)
			if err != nil {
				if !test.WantErr {
					t.Errorf("unexpected error: %s", err.Error())
				}
				if err.Error() != test.Err {
					t.Errorf("unexpected error. Expected %q but got %q", test.Err, err.Error())
				}
				return
			}
			if test.WantErr {
				t.Errorf("expected error but nothing got")
			}
			logger.Log.Info("Get stocks by id", "entities", entities)

			for _, entity := range entities {
				if entity.GoodID != test.ID {
					t.Errorf("wrong good id. Expected %d but got %d", test.ID, entity.GoodID)
				}
			}
		})
	}
}

func TestGetGoodsAmountByWarehouseID(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresContainer, mappedPort := tests.CreateInventoryHubTestPostgresContainer(ctx)
	defer tests.TerminateInventoryHubTestContainer(ctx, postgresContainer)

	postgresClient, err := postgres.NewClient(postgres.Config{
		User:       tests.PostgresInventoryHubTestUser,
		Password:   tests.PostgresInventoryHubTestPassword,
		Host:       tests.PostgresInventoryHubTestHost + ":" + mappedPort,
		DBName:     tests.PostgresInventoryHubTestDB,
		DisableTLS: true,
	})
	if err != nil {
		t.Errorf("error with starting postgres client: %s", err.Error())
	}
	defer postgresClient.Close()

	store := New(postgresClient)

	if err := seedTestData(ctx, true, postgresClient); err != nil {
		t.Fatalf("error with seeding test data: %s", err)
	}

	for _, test := range []struct {
		Name    string
		ID      int
		Amount  int
		WantErr bool
		Err     string
	}{
		{
			Name:   "Successful get a goods amount by warehouse id",
			ID:     1,
			Amount: 100,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			entities, err := store.GetGoodsAmountByWarehouseID(ctx, test.ID)
			if err != nil {
				if !test.WantErr {
					t.Errorf("unexpected error: %s", err.Error())
				}
				if err.Error() != test.Err {
					t.Errorf("unexpected error. Expected %q but got %q", test.Err, err.Error())
				}
				return
			}
			if test.WantErr {
				t.Errorf("expected error but nothing got")
			}
			logger.Log.Info("Get goods amount by warehouse id", "entities", entities)

			for _, entity := range entities {
				if entity.Amount != test.Amount {
					t.Errorf("wrong amount. Expected %d but got %d", test.Amount, entity.Amount)
				}
				if entity.Reserved != 100 {
					t.Errorf("reserved amount should be 100 but got %d", entity.Reserved)
				}
			}
		})
	}
}

func TestReserveAmount(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresContainer, mappedPort := tests.CreateInventoryHubTestPostgresContainer(ctx)
	defer tests.TerminateInventoryHubTestContainer(ctx, postgresContainer)

	postgresClient, err := postgres.NewClient(postgres.Config{
		User:       tests.PostgresInventoryHubTestUser,
		Password:   tests.PostgresInventoryHubTestPassword,
		Host:       tests.PostgresInventoryHubTestHost + ":" + mappedPort,
		DBName:     tests.PostgresInventoryHubTestDB,
		DisableTLS: true,
	})
	if err != nil {
		t.Fatalf("error with starting postgres client: %s", err)
	}
	defer postgresClient.Close()

	store := New(postgresClient)

	if err := seedTestData(ctx, true, postgresClient); err != nil {
		t.Fatalf("error with seeding test data: %s", err)
	}

	for _, test := range []struct {
		Name    string
		Data    []ChangeAmountEntity
		WantErr bool
		Err     string
	}{
		{
			Name: "Successful reserve some goods",
			Data: []ChangeAmountEntity{
				{GoodID: 1, Amount: 10},
				{GoodID: 2, Amount: 15},
				{GoodID: 3, Amount: 30},
			},
		},
		{
			Name: "Successful reserve a lot of goods",
			Data: []ChangeAmountEntity{
				{GoodID: 3, Amount: 110},
				{GoodID: 4, Amount: 130},
				{GoodID: 5, Amount: 180},
			},
		},
		{
			Name: "Out of reserve amount goods",
			Data: []ChangeAmountEntity{
				{GoodID: 6, Amount: 10000},
			},
			WantErr: true,
			Err:     core.ErrNotEnoughAmount.Error(),
		},
		{
			Name: "Reserve non-existent goods",
			Data: []ChangeAmountEntity{
				{GoodID: 1000, Amount: 100},
			},
			WantErr: true,
			Err:     sql.ErrNoRows.Error(),
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			err := store.Reserve(ctx, test.Data)
			if err != nil {
				if !test.WantErr {
					t.Errorf("unexpected error: %s", err.Error())
				}
				if err.Error() != test.Err {
					t.Errorf("unexpected error. Expected %s but got %s", test.Err, err.Error())
				}
				return
			}
			if test.WantErr {
				t.Errorf("expected error but nothing got")
			}

			for _, data := range test.Data {
				entities, err := store.GetByGoodID(ctx, data.GoodID)
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				}
				logger.Log.Info("Stock for good id",
					"good id", data.GoodID,
					"entities", entities)

				for _, entity := range entities {
					if entity.Amount < 0 {
						t.Errorf("Amount should be positive")
					}
				}
			}
		})
	}
}

func TestReleaseAmount(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresContainer, mappedPort := tests.CreateInventoryHubTestPostgresContainer(ctx)
	defer tests.TerminateInventoryHubTestContainer(ctx, postgresContainer)

	postgresClient, err := postgres.NewClient(postgres.Config{
		User:       tests.PostgresInventoryHubTestUser,
		Password:   tests.PostgresInventoryHubTestPassword,
		Host:       tests.PostgresInventoryHubTestHost + ":" + mappedPort,
		DBName:     tests.PostgresInventoryHubTestDB,
		DisableTLS: true,
	})
	if err != nil {
		t.Fatalf("error with starting postgres client: %s", err)
	}
	defer postgresClient.Close()

	store := New(postgresClient)

	if err := seedTestData(ctx, true, postgresClient); err != nil {
		t.Fatalf("error with seeding test data: %s", err)
	}

	for _, test := range []struct {
		Name    string
		Data    []ChangeAmountEntity
		WantErr bool
		Err     string
	}{
		{
			Name: "Successful release some goods",
			Data: []ChangeAmountEntity{
				{GoodID: 1, Amount: 10},
				{GoodID: 2, Amount: 15},
				{GoodID: 3, Amount: 30},
			},
		},
		{
			Name: "Successful release a lot of goods",
			Data: []ChangeAmountEntity{
				{GoodID: 3, Amount: 110},
				{GoodID: 4, Amount: 130},
				{GoodID: 5, Amount: 180},
			},
		},
		{
			Name: "Out of release amount goods",
			Data: []ChangeAmountEntity{
				{GoodID: 6, Amount: 10000},
			},
			WantErr: true,
			Err:     core.ErrNotEnoughReserve.Error(),
		},
		{
			Name: "Release non-existent goods",
			Data: []ChangeAmountEntity{
				{GoodID: 1000, Amount: 100},
			},
			WantErr: true,
			Err:     sql.ErrNoRows.Error(),
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			err := store.Release(ctx, test.Data)
			if err != nil {
				if !test.WantErr {
					t.Errorf("unexpected error: %s", err.Error())
				}
				if err.Error() != test.Err {
					t.Errorf("unexpected error. Expected %s but got %s", test.Err, err.Error())
				}
				return
			}
			if test.WantErr {
				t.Errorf("expected error but nothing got")
			}

			for _, data := range test.Data {
				entities, err := store.GetByGoodID(ctx, data.GoodID)
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				}
				logger.Log.Info("Stock for good id",
					"good id", data.GoodID,
					"entities", entities)

				for _, entity := range entities {
					if entity.Reserved < 0 {
						t.Errorf("Reserved should be positive")
					}
				}
			}
		})
	}
}

func seedTestData(ctx context.Context, withStocks bool, client *postgres.Client) error {
	for i := 1; i <= 10; i++ {
		_, err := client.ExecContext(ctx,
			`INSERT INTO warehouses (name, is_available) VALUES ($1, $2)`,
			fmt.Sprintf("Warehouse %d", i), true)
		if err != nil {
			return err
		}
	}

	sizes := []string{"s", "m", "l"}
	for i := 1; i <= 10; i++ {
		_, err := client.ExecContext(ctx,
			`INSERT INTO goods (name, size) VALUES ($1, $2)`,
			fmt.Sprintf("Good %d", i), sizes[i%len(sizes)])
		if err != nil {
			return err
		}
	}

	if withStocks {
		for goodID := 1; goodID <= 10; goodID++ {
			for warehouseID := 1; warehouseID <= 10; warehouseID++ {
				_, err := client.ExecContext(ctx,
					`INSERT INTO stocks (good_id, warehouse_id, amount, reserved) VALUES ($1, $2, $3, $4)`,
					goodID, warehouseID, 100, 100)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
