package tests

import (
	"context"
	"fmt"
	"lamoda/pkg/logger"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	PostgresInventoryHubTestHost     = "localhost"
	postgresInventoryHubTestPort     = "5432"
	PostgresInventoryHubTestUser     = "postgres"
	PostgresInventoryHubTestPassword = "postgres"
	PostgresInventoryHubTestDB       = "inventory_hub"
)

func CreateInventoryHubTestPostgresContainer(ctx context.Context) (testcontainers.Container, string) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15.3-bullseye",
		ExposedPorts: []string{postgresInventoryHubTestPort + "/tcp"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		Env: map[string]string{
			"POSTGRES_USER":     PostgresInventoryHubTestUser,
			"POSTGRES_PASSWORD": PostgresInventoryHubTestPassword,
			"POSTGRES_DB":       PostgresInventoryHubTestDB,
		},
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic("starting postgres container: " + err.Error())
	}

	mappedPort, err := postgresContainer.MappedPort(ctx, postgresInventoryHubTestPort)
	if err != nil {
		panic("getting mapped port: " + err.Error())
	}

	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		PostgresInventoryHubTestUser,
		PostgresInventoryHubTestPassword,
		PostgresInventoryHubTestHost+":"+mappedPort.Port(),
		PostgresInventoryHubTestDB,
	)

	// Wait for the container to be ready
	time.Sleep(500 * time.Millisecond)
	if err = runMigrations(ctx, databaseURL); err != nil {
		panic("run migrations: " + err.Error())
	}

	if err := logger.SetupLogger("local"); err != nil {
		panic("setting up logger: " + err.Error())
	}

	return postgresContainer, mappedPort.Port()
}

func TerminateInventoryHubTestContainer(ctx context.Context, container testcontainers.Container) {
	if err := container.Terminate(ctx); err != nil {
		panic(err)
	}
}

func runMigrations(ctx context.Context, connString string) error {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	migrator, err := migrate.NewMigrator(context.Background(), conn, "schema_version")
	if err != nil {
		return err
	}

	err = migrator.LoadMigrations(os.DirFS("../../../migrations/sql"))
	if err != nil {
		return err
	}

	return migrator.Migrate(context.Background())
}
