package api

import (
	"lamoda/api/hub/rest/handlers"

	"github.com/labstack/echo/v4"
)

type Container struct {
	Server   *echo.Echo
	Resolver *handlers.Resolver
}

// NewContainer creates a new application struct.
func NewContainer(
	server *echo.Echo,
	resolver *handlers.Resolver,
) Container {
	return Container{
		Resolver: resolver,
		Server:   server,
	}
}
