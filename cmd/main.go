package main

import (
	"github.com/urfave/cli/v2"
)

var globalFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "log-level",
		Usage:   "log level",
		EnvVars: []string{"LOG_LEVEL"},
		Value:   "local",
	},
}

func main() {}
