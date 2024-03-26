package swagger

import "github.com/urfave/cli/v2"

var cmdFlags = []cli.Flag{
	&cli.IntFlag{
		Name:    "swagger-port",
		Usage:   "swagger port",
		EnvVars: []string{"SWAGGER_PORT"},
		Value:   9000,
	},
}
