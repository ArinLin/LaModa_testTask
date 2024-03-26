package handlers

import (
	"database/sql"
	"net/http"

	"lamoda/api/hub/rest/presenters/warehouse"
	"lamoda/internal/core"
	"lamoda/internal/service/hub"
	"lamoda/pkg/web"
	"lamoda/pkg/webutil"

	"github.com/labstack/echo/v4"
)

func (r *Resolver) getWarehouses(c echo.Context) error {
	data, err := r.service.GetWarehouses(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, err.Error(), nil))
	}

	pres := warehouse.PresentList(data)

	return c.JSON(http.StatusOK, pres.Response())
}

func (r *Resolver) getWarehouseByID(c echo.Context) error {
	id, err := webutil.ParseID(c)
	if err != nil {
		return err
	}

	data, err := r.service.GetWarehouseByID(c.Request().Context(), id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return c.JSON(http.StatusNotFound, web.ErrorResponse(core.WarehouseNotFoundCode, nil, nil))
		default:
			return c.JSON(http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, err.Error(), nil))
		}
	}

	pres := warehouse.PresentWarehouse(data)

	return c.JSON(http.StatusOK, pres.Response(core.WarehouseReceivedCode))
}

func (r *Resolver) createWarehouse(c echo.Context) error {
	var model hub.CreateWarehouseModel
	if err := webutil.BodyChecker(c, &model); err != nil {
		return err
	}

	data, err := r.service.CreateWarehouse(c.Request().Context(), model)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, err.Error(), nil))
	}

	pres := warehouse.PresentWarehouse(data)

	return c.JSON(http.StatusCreated, pres.Response(core.WarehouseCreatedCode))
}

func (r *Resolver) updateWarehouse(c echo.Context) error {
	id, err := webutil.ParseID(c)
	if err != nil {
		return err
	}

	var model hub.UpdateWarehouseModel
	if err := webutil.BodyChecker(c, &model); err != nil {
		return err
	}

	data, err := r.service.UpdateWarehouse(c.Request().Context(), id, model)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return c.JSON(http.StatusNotFound, web.ErrorResponse(core.WarehouseNotFoundCode, nil, nil))
		default:
			return c.JSON(http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, err.Error(), nil))
		}
	}

	pres := warehouse.PresentWarehouse(data)

	return c.JSON(http.StatusOK, pres.Response(core.WarehouseUpdatedCode))
}

func (r *Resolver) deleteWarehouse(c echo.Context) error {
	id, err := webutil.ParseID(c)
	if err != nil {
		return err
	}

	if err = r.service.DeleteWarehouse(c.Request().Context(), id); err != nil {
		switch err {
		case sql.ErrNoRows:
			return c.JSON(http.StatusNotFound, web.ErrorResponse(core.WarehouseNotFoundCode, nil, nil))
		default:
			return c.JSON(http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, err.Error(), nil))
		}
	}

	return c.JSON(http.StatusOK, web.OKResponse(core.WarehouseDeletedCode, nil, nil))
}
