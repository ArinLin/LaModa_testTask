package swagger

import (
	"fmt"

	_ "lamoda/api/doc" // swagger doc

	"github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
	"github.com/urfave/cli/v2"
)

var Cmd = cli.Command{
	Name:  "swagger",
	Usage: "Provide swagger documentation",
	Flags: cmdFlags,
	OnUsageError: func(c *cli.Context, _ error, _ bool) error {
		return cli.ShowCommandHelp(c, "swagger")
	},
	Action: run,
}

func run(c *cli.Context) error {
	e := echo.New()

	e.GET("/doc/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", c.Int("swagger-port"))))

	return nil
}
