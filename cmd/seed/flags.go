package seed

import "github.com/urfave/cli/v2"

var cmdFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "postgres-user",
		Usage:   "postgres user",
		EnvVars: []string{"POSTGRES_USER"},
		Value:   "postgres",
	},
	&cli.StringFlag{
		Name:    "postgres-password",
		Usage:   "postgres password",
		EnvVars: []string{"POSTGRES_PASSWORD"},
		Value:   "postgres",
	},
	&cli.StringFlag{
		Name:    "postgres-host",
		Usage:   "postgres host",
		EnvVars: []string{"POSTGRES_HOST"},
		Value:   "localhost:5432",
	},
	&cli.StringFlag{
		Name:    "postgres-db-name",
		Usage:   "postgres database name",
		EnvVars: []string{"POSTGRES_DB_NAME"},
		Value:   "inventory_hub",
	},
	&cli.BoolFlag{
		Name:    "postgres-disable-tls",
		Usage:   "postgres disable tls",
		EnvVars: []string{"POSTGRES_DISABLE_TLS"},
		Value:   true,
	},
	&cli.IntFlag{
		Name:    "goods",
		Usage:   "goods count",
		EnvVars: []string{"GOODS"},
		Value:   15,
	},
	&cli.IntFlag{
		Name:    "stocks",
		Usage:   "stocks count",
		EnvVars: []string{"STOCKS"},
		Value:   100,
	},
	&cli.IntFlag{
		Name:    "warehouses",
		Usage:   "warehouses count",
		EnvVars: []string{"WAREHOUSES"},
		Value:   5,
	},
}
