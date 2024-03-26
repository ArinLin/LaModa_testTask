package inject

import (
	"lamoda/internal/store/good"
	"lamoda/internal/store/stock"
	"lamoda/internal/store/warehouse"
	"lamoda/pkg/postgres"

	"github.com/google/wire"
	"github.com/urfave/cli/v2"
)

var storeSet = wire.NewSet( // nolint
	providePostgresConfig,
	postgres.NewClient,
	good.New,
	stock.New,
	warehouse.New,
)

func providePostgresConfig(c *cli.Context) postgres.Config {
	return postgres.Config{
		User:       c.String("postgres-user"),
		Password:   c.String("postgres-password"),
		Host:       c.String("postgres-host"),
		DBName:     c.String("postgres-db-name"),
		DisableTLS: c.Bool("postgres-disable-tls"),
	}
}
