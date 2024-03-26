package webutil

import (
	"errors"
	"net/http"
	"strconv"

	"lamoda/internal/core"
	"lamoda/pkg/web"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// var validate *validator.Validate

func BodyChecker(c echo.Context, entity interface{}) error {
	if err := c.Bind(entity); err != nil {
		c.JSON(http.StatusBadRequest, web.ErrorResponse(core.InvalidBodyCode, err.Error(), nil))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return Validate(c, entity)
}

func ParseID(c echo.Context) (int, error) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, web.ErrorResponse(core.IDRequiredCode, nil, nil))
		return 0, echo.NewHTTPError(http.StatusBadRequest, core.IDRequiredCode)
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, web.ErrorResponse(core.InvalidIDCode, err.Error(), nil))
		return 0, echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return id, nil
}

func Validate(c echo.Context, entity interface{}) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(entity); err != nil {
		var verrors validator.ValidationErrors
		ok := errors.As(err, &verrors)
		if !ok {
			c.JSON(http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, nil, nil))
			return echo.NewHTTPError(http.StatusInternalServerError, core.InternalErrorCode)
		}
		c.JSON(http.StatusUnprocessableEntity, web.ValidationErrorResponse(web.ValidationErrors(verrors), nil))
		return echo.NewHTTPError(http.StatusUnprocessableEntity, verrors)
	}

	return nil
}
