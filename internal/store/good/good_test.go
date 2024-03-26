package good

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
			Name: "Successful creating a good",
			Data: CreateEntity{
				Name: "T-Shirt",
				Size: "s",
			},
		},
		{
			Name: "Error creating a good with unsupported size",
			Data: CreateEntity{
				Name: "Pants",
				Size: "square",
			},
			WantErr: true,
			Err:     "ERROR: invalid input value for enum size: \"square\" (SQLSTATE 22P02)",
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
			logger.Log.Info("Created good", "entity", entity)

			if entity.ID == 0 {
				t.Errorf("entity ID is 0")
			}
			if entity.Name != test.Data.Name {
				t.Errorf("wrong name. Expected %q but got %q", test.Data.Name, entity.Name)
			}
			if entity.Size != test.Data.Size {
				t.Errorf("wrong size. Expected %q but got %q", test.Data.Size, entity.Size)
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

	good, err := store.Create(ctx, CreateEntity{
		Name: "Hat",
		Size: "m",
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
			Name: "Successful get a good",
			ID:   good.ID,
		},
		{
			Name:    "Get non-existent good",
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
			logger.Log.Info("Get good by id", "entity", entity)

			if entity.Name != good.Name {
				t.Errorf("wrong name. Expected %q but got %q", good.Name, entity.Name)
			}
			if entity.Size != good.Size {
				t.Errorf("wrong size. Expected %q but got %q", good.Size, entity.Size)
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

	var goods []Entity
	for _, data := range []CreateEntity{
		{
			Name: "Hat",
			Size: "m",
		},
		{
			Name: "Pants",
			Size: "l",
		},
		{
			Name: "T-Shirt",
			Size: "l",
		},
		{
			Name: "Hoodie",
			Size: "s",
		},
		{
			Name: "Socks",
			Size: "m",
		},
	} {
		entity, err := store.Create(ctx, data)
		if err != nil {
			t.Errorf("error with creating test data: %s", err.Error())
		}

		goods = append(goods, entity)
	}

	for _, test := range []struct {
		Name    string
		WantErr bool
		Err     string
	}{
		{
			Name: "Successful get all goods",
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
			logger.Log.Info("Get goods list", "entity", entity)

			for i, good := range goods {
				if entity[i].Name != good.Name {
					t.Errorf("wrong name. Expected %q but got %q", good.Name, entity[i].Name)
				}
				if entity[i].Size != good.Size {
					t.Errorf("wrong size. Expected %q but got %q", good.Size, entity[i].Size)
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

	var goods []Entity
	for _, data := range []CreateEntity{
		{
			Name: "T-Shirt",
			Size: "s",
		},
		{
			Name: "Jacket",
			Size: "m",
		},
		{
			Name: "Shorts",
			Size: "s",
		},
	} {
		entity, err := store.Create(ctx, data)
		if err != nil {
			t.Errorf("error with creating goods: %s", err.Error())
		}

		goods = append(goods, entity)
	}

	newName := "Bomber"
	newSize := "l"

	var id int
	for _, test := range []struct {
		Name    string
		Data    UpdateEntity
		WantErr bool
		Err     string
	}{
		{
			Name: "Successful full update a good",
			Data: UpdateEntity{
				Name: &newName,
				Size: &newSize,
			},
		},
		{
			Name: "Successful update name only",
			Data: UpdateEntity{
				Name: &newName,
			},
		},
		{
			Name: "Successful update size only",
			Data: UpdateEntity{
				Size: &newSize,
			},
		},

		{
			Name:    "Update non-existent good",
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
			logger.Log.Info("Update good", "entity", entity)

			if test.Data.Name != nil {
				if entity.Name != *test.Data.Name {
					t.Errorf("wrong name. Expected %q but got %q", *test.Data.Name, entity.Name)
				}
			} else {
				if entity.Name != goods[id-1].Name {
					t.Errorf("wrong name. Expected %q but got %q", goods[id-1].Name, entity.Name)
				}
			}
			if test.Data.Size != nil {
				if entity.Size != *test.Data.Size {
					t.Errorf("wrong size. Expected %q but got %q", *test.Data.Size, entity.Size)
				}
			} else {
				if entity.Size != goods[id-1].Size {
					t.Errorf("wrong size. Expected %q but got %q", goods[id-1].Size, entity.Size)
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

	good, err := store.Create(ctx, CreateEntity{
		Name: "Gloves",
		Size: "l",
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
			Name: "Successful delete a good",
			ID:   good.ID,
		},
		{
			Name: "Delete already deleted good",
			ID:   good.ID,
		},
		{
			Name:    "Delete non-existent good",
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
		})
	}
}
