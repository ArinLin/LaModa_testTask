package seed

import (
	"lamoda/internal/service/seed"
	"lamoda/pkg/postgres"

	"github.com/urfave/cli/v2"
)

var Cmd = cli.Command{
	Name:  "seed",
	Usage: "Fill database with data",
	Flags: cmdFlags,
	OnUsageError: func(c *cli.Context, _ error, _ bool) error {
		return cli.ShowCommandHelp(c, "seed")
	},
	Action: run,
}

func run(c *cli.Context) error {
	postgresCfg := postgres.Config{
		User:       c.String("postgres-user"),
		Password:   c.String("postgres-password"),
		Host:       c.String("postgres-host"),
		DBName:     c.String("postgres-db-name"),
		DisableTLS: c.Bool("postgres-disable-tls"),
	}
	postgresClient, err := postgres.NewClient(postgresCfg)
	if err != nil {
		return err
	}

	empty, err := seed.CheckIfDatabaseIsEmpty(c.Context, postgresClient)
	if err != nil {
		return err
	}
	if !empty {
		return nil
	}

	counter := seed.Counter{
		Goods:      c.Int("goods"),
		Stocks:     c.Int("stocks"),
		Warehouses: c.Int("warehouses"),
	}
	err = seed.Seed(c.Context, postgresClient, counter)
	if err != nil {
		return err
	}

	postgresClient.Close()
	return nil
}
