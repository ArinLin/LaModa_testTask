package warehouse

import (
	"context"
	"database/sql"
	"testing"

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

	for _, test := range []struct {
		Name    string
		Data    CreateEntity
		WantErr bool
		Err     string
	}{
		{
			Name: "Successful creating a warehouse",
			Data: CreateEntity{
				Name:        "Logopark Sofino",
				IsAvailable: true,
			},
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
			logger.Log.Info("Created warehouse", "entity", entity)

			if entity.ID == 0 {
				t.Errorf("entity ID is 0")
			}
			if entity.Name != test.Data.Name {
				t.Errorf("wrong name. Expected %q but got %q", test.Data.Name, entity.Name)
			}
			if entity.IsAvailable != test.Data.IsAvailable {
				t.Errorf("wrong availability. Expected %t but got %t", test.Data.IsAvailable, entity.IsAvailable)
			}
			if entity.DeletedAt != nil {
				t.Errorf("deleted at is not nil")
			}
		})
	}
}

func TestGetByID(t *testing.T) {
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

	warehouse, err := store.Create(ctx, CreateEntity{
		Name:        "Logopark Bykovo",
		IsAvailable: true,
	})
	if err != nil {
		t.Errorf("error with creating test data: %s", err.Error())
	}

	for _, test := range []struct {
		Name    string
		ID      int
		WantErr bool
		Err     string
	}{
		{
			Name: "Successful get a warehouse",
			ID:   warehouse.ID,
		},
		{
			Name:    "Get non-existent warehouse",
			WantErr: true,
			Err:     sql.ErrNoRows.Error(),
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			entity, err := store.GetByID(ctx, test.ID)
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
			logger.Log.Info("Get warehouse by id", "entity", entity)

			if entity.Name != warehouse.Name {
				t.Errorf("wrong name. Expected %q but got %q", warehouse.Name, entity.Name)
			}
			if entity.IsAvailable != warehouse.IsAvailable {
				t.Errorf("wrong availability. Expected %t but got %t", warehouse.IsAvailable, entity.IsAvailable)
			}
			if entity.DeletedAt != nil {
				t.Error("wrong deleted at. Expected nil")
			}
		})
	}
}

func TestGetAll(t *testing.T) {
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

	var warehouses []Entity
	for _, data := range []CreateEntity{
		{
			Name:        "Baklovif",
			IsAvailable: true,
		},
		{
			Name:        "Gedmicuc",
			IsAvailable: true,
		},
		{
			Name:        "Kotunzas",
			IsAvailable: false,
		},
		{
			Name:        "Guickuw",
			IsAvailable: true,
		},
		{
			Name:        "Wabuiko",
			IsAvailable: false,
		},
	} {
		entity, err := store.Create(ctx, data)
		if err != nil {
			t.Errorf("error with creating test data: %s", err.Error())
		}

		warehouses = append(warehouses, entity)
	}

	for _, test := range []struct {
		Name string

		WantErr bool
		Err     string
	}{
		{
			Name: "Successful get all warehouses",
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			entity, err := store.GetAll(ctx)
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
			logger.Log.Info("Get warehouses list", "entity", entity)

			if len(entity) != len(warehouses) {
				t.Errorf("wrong number of warehouses. Expected %d but got %d", len(warehouses), len(entity))
			}
			for i, warehouse := range warehouses {
				if entity[i].Name != warehouse.Name {
					t.Errorf("wrong name. Expected %q but got %q", warehouse.Name, entity[i].Name)
				}
				if entity[i].IsAvailable != warehouse.IsAvailable {
					t.Errorf("wrong size. Expected %t but got %t", warehouse.IsAvailable, entity[i].IsAvailable)
				}
				if entity[i].DeletedAt != nil {
					t.Error("wrong deleted at. Expected nil")
				}
			}
			for i := 0; i < len(entity)-1; i++ {
				if entity[i].ID > entity[i+1].ID {
					t.Errorf("wrong order. Expected %d to be after %d", i+1, i)
				}
			}
		})
	}
}

func TestUpdate(t *testing.T) {
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

	var warehouses []Entity
	for _, data := range []CreateEntity{
		{
			Name:        "Pemdulbup",
			IsAvailable: true,
		},
		{
			Name:        "Vujonfu",
			IsAvailable: true,
		},
		{
			Name:        "Pamnijuv",
			IsAvailable: true,
		},
	} {
		entity, err := store.Create(ctx, data)
		if err != nil {
			t.Errorf("error with creating warehouse: %s", err.Error())
		}

		warehouses = append(warehouses, entity)
	}

	newName := "Urhiffu"
	anotherNewName := "Magagpi"
	newAvailable := false

	var id int
	for _, test := range []struct {
		Name    string
		Data    UpdateEntity
		WantErr bool
		Err     string
	}{
		{
			Name: "Successful full update a warehouse",
			Data: UpdateEntity{
				Name:        &newName,
				IsAvailable: &newAvailable,
			},
		},
		{
			Name: "Successful update name only",
			Data: UpdateEntity{
				Name: &anotherNewName,
			},
		},
		{
			Name: "Successful update availability only",
			Data: UpdateEntity{
				IsAvailable: &newAvailable,
			},
		},

		{
			Name:    "Update non-existent warehouse",
			WantErr: true,
			Err:     sql.ErrNoRows.Error(),
		},
	} {
		id++

		t.Run(test.Name, func(t *testing.T) {
			entity, err := store.Update(ctx, id, test.Data)
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
			logger.Log.Info("Update warehouse", "entity", entity)

			if test.Data.Name != nil {
				if entity.Name != *test.Data.Name {
					t.Errorf("wrong name. Expected %q but got %q", *test.Data.Name, entity.Name)
				}
			} else {
				if entity.Name != warehouses[id-1].Name {
					t.Errorf("wrong name. Expected %q but got %q", warehouses[id-1].Name, entity.Name)
				}
			}
			if test.Data.IsAvailable != nil {
				if entity.IsAvailable != *test.Data.IsAvailable {
					t.Errorf("wrong availability. Expected %t but got %t", *test.Data.IsAvailable, entity.IsAvailable)
				}
			} else {
				if entity.IsAvailable != warehouses[id-1].IsAvailable {
					t.Errorf("wrong availability. Expected %t but got %t", warehouses[id-1].IsAvailable, entity.IsAvailable)
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
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

	warehouse, err := store.Create(ctx, CreateEntity{
		Name:        "Cuketkes",
		IsAvailable: true,
	})
	if err != nil {
		t.Errorf("error with creating test data: %s", err.Error())
	}

	for _, test := range []struct {
		Name    string
		ID      int
		WantErr bool
		Err     string
	}{
		{
			Name: "Successful delete a warehouse",
			ID:   warehouse.ID,
		},
		{
			Name: "Delete already deleted warehouse",
			ID:   warehouse.ID,
		},
		{
			Name:    "Delete non-existent warehouse",
			ID:      100,
			WantErr: true,
			Err:     sql.ErrNoRows.Error(),
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			err := store.Delete(ctx, test.ID)
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

			entity, err := store.GetByID(ctx, test.ID)
			if err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}
			if entity.DeletedAt == nil {
				t.Errorf("deleted at should not be nil")
			}
			if entity.IsAvailable {
				t.Error("availability should be false")
			}
		})
	}
}
