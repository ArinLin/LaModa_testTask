package hub

import "github.com/urfave/cli/v2"

var cmdFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "server-host",
		Usage:   "server host",
		EnvVars: []string{"SERVER_HOST"},
		Value:   "localhost:8080",
	},
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
}
