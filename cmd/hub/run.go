package hub

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"lamoda/api/hub/inject"
	"lamoda/pkg/logger"

	"github.com/urfave/cli/v2"
)

var Cmd = cli.Command{
	Name:  "hub",
	Usage: "Run service",
	Flags: cmdFlags,
	OnUsageError: func(c *cli.Context, _ error, _ bool) error {
		return cli.ShowCommandHelp(c, "hub")
	},
	Action: run,
}

func run(c *cli.Context) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
			return
		case s := <-sig:
			logger.Log.Info("signal received",
				"signal", s.String())
			cancel()
		}
	}()

	app, err := inject.InitializeApplication(c, ctx)
	if err != nil {
		logger.Log.Error("main: cannot initialize server",
			"error", err.Error())
		os.Exit(1)
	}

	go func() {
		app.Server.Start(c.String("server-host"))
	}()

	<-ctx.Done()

	_ = app.Server.Shutdown(ctx)
	logger.Log.Debug("ctx end received")

	return nil
}
