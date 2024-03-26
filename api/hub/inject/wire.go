//go:build wireinject
// +build wireinject

package inject

import (
	"context"
	"lamoda/api/hub"

	"github.com/google/wire"
	"github.com/urfave/cli/v2"
)

func InitializeApplication(c *cli.Context, appCtx context.Context) (api.Container, error) {
	wire.Build(
		serverSet,
		storeSet,
		api.NewContainer,
	)
	return api.Container{}, nil
}
