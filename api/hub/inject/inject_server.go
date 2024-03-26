package inject

import (
	"lamoda/api/hub/rest/handlers"
	"lamoda/internal/service/hub"

	"github.com/labstack/echo/v4"

	"github.com/google/wire"
)

var serverSet = wire.NewSet( // nolint
	hub.New,
	echo.New,
	handlers.NewResolver,
)
