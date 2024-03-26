package main

import (
	"os"

	"lamoda/cmd/hub"
	"lamoda/cmd/seed"
	"lamoda/cmd/swagger"
	"lamoda/pkg/logger"

	"github.com/urfave/cli/v2"
)

var globalFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "log-level",
		Usage:   "log level",
		EnvVars: []string{"LOG-LEVEL"},
		Value:   "local",
	},
}

// @title		Lamoda Intentory Hub API
// @version		1.0
// @description	Lamoda Intentory Hub API Service Documentation
// @host		server:8080
// @BasePath	/hub/api/v1
func main() {
	app := &cli.App{
		Usage: "Lamoda Intentory Hub",
		Commands: []*cli.Command{
			&hub.Cmd,
			&seed.Cmd,
			&swagger.Cmd,
		},
		Flags: globalFlags,
		OnUsageError: func(c *cli.Context, _ error, _ bool) error {
			return cli.ShowAppHelp(c)
		},
		Before: func(ctx *cli.Context) error {
			return logger.SetupLogger(ctx.String("log-level"))
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Log.Error("run app",
			"error", err)
	}
}
