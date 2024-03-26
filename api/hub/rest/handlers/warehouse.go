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

// @Summary		Get Warehouses
// @Description	Get Warehouses List
// @Tags			Warehouses
// @Produce		json
// @Success		200	{object}	model.GetWarehousesResponse
// @Failure		500	{object}	model.InternalResponse
// @Router			/warehouses [get]
func (r *Resolver) getWarehouses(c echo.Context) error {
	data, err := r.service.GetWarehouses(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, err.Error(), nil))
	}

	pres := warehouse.PresentList(data)

	return c.JSON(http.StatusOK, pres.Response())
}

// @Summary		Get Warehouse By ID
// @Description	Get Warehouse By ID
// @Tags			Warehouses
// @Produce		json
// @Param			id	path		int	true	"Warehouse ID"
// @Success		200	{object}	model.GetWarehouseByIDResponse
// @Failure		400	{object}	model.BadRequestInvalidIDResponse
// @Failure		404	{object}	model.WarehouseNotFoundResponse
// @Failure		500	{object}	model.InternalResponse
// @Router			/warehouse/{id} [get]
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

// @Summary		Create Warehouse
// @Description	Create Warehouse
// @Tags			Warehouses
// @Produce		json
// @Param			warehouse	body		model.CreateWarehouseRequest	true	"Params to create warehouse"
// @Success		200	{object}	model.CreateWarehouseResponse
// @Failure		400	{object}	model.BadRequestInvalidBodyResponse
// @Failure		422		{object}	model.ValidationResponse
// @Failure		500	{object}	model.InternalResponse
// @Router			/warehouse [post]
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

// @Summary		Update Warehouse
// @Description	Update Warehouse By ID
// @Tags			Warehouses
// @Produce		json
// @Param			id	path		int	true	"Warehouse ID"
// @Param			warehouse	body		model.UpdateWarehouseRequest	true	"Params to update warehouse"
// @Success		200	{object}	model.UpdateWarehouseResponse
// @Failure		400	{object}	model.BadRequestInvalidIDResponse
// @Failure		404	{object}	model.WarehouseNotFoundResponse
// @Failure		422		{object}	model.ValidationResponse
// @Failure		500	{object}	model.InternalResponse
// @Router			/warehouse/{id} [patch]
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

// @Summary		Delete Warehouse
// @Description	Delete Warehouse By Id
// @Tags			Warehouses
// @Produce		json
// @Param			id	path		int	true	"Warehouse ID"
// @Success		200	{object}	model.DeleteWarehouseResponse
// @Failure		400	{object}	model.BadRequestInvalidIDResponse
// @Failure		404	{object}	model.WarehouseNotFoundResponse
// @Failure		500	{object}	model.InternalResponse
// @Router			/warehouse/{id} [delete]
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
