package handlers

import (
	"database/sql"
	"net/http"

	"lamoda/api/hub/rest/presenters/stock"
	"lamoda/internal/core"
	"lamoda/internal/service/hub"
	"lamoda/pkg/web"
	"lamoda/pkg/webutil"

	"github.com/labstack/echo/v4"
)

// @Summary		Get Stocks By Warehouse ID
// @Description	Get Stocks By Warehouse ID
// @Tags			Stocks
// @Produce		json
// @Param			id	path		int	true	"Warehouse ID"
// @Success		200	{object}	model.GetStocksResponse
// @Failure		404	{object}	model.WarehouseNotFoundResponse
// @Failure		500	{object}	model.InternalResponse
// @Router			/stocks/warehouse/{id} [get]
func (r *Resolver) getStockByWarehouseID(c echo.Context) error {
	id, err := webutil.ParseID(c)
	if err != nil {
		return err
	}

	data, err := r.service.GetStockByWarehouseID(c.Request().Context(), id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return c.JSON(http.StatusNotFound, web.ErrorResponse(core.WarehouseNotFoundCode, nil, nil))
		default:
			return c.JSON(http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, err.Error(), nil))
		}
	}

	pres := stock.PresentList(data)

	return c.JSON(http.StatusOK, pres.Response())
}

// @Summary		Reserve Stock
// @Description	Reserve Stock
// @Tags			Stocks
// @Produce		json
// @Param			warehouse	body		model.ChangeAmountRequest	true	"List of goods with amount to reserve"
// @Success		200	{object}	model.ReserveStockResponse
// @Failure		400	{object}	model.BadRequestInvalidBodyResponse
// @Failure		404	{object}	model.GoodNotFoundResponse
// @Failure		422	{object}	model.ValidationResponse
// @Failure		500	{object}	model.InternalResponse
// @Router		/stocks/reserve [put]
func (r *Resolver) reserveStock(c echo.Context) error {
	var model hub.ChangeStocksAmountModel
	if err := webutil.BodyChecker(c, &model); err != nil {
		return err
	}

	if err := r.service.ReserveStocks(c.Request().Context(), model); err != nil {
		switch err {
		case sql.ErrNoRows:
			return c.JSON(http.StatusNotFound, web.ErrorResponse(core.GoodNotFoundCode, nil, nil))
		case core.ErrNotEnoughAmount:
			return c.JSON(http.StatusBadRequest, web.ErrorResponse(core.NotEnoughAmountCode, nil, nil))
		default:
			return c.JSON(http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, err.Error(), nil))
		}
	}

	return c.JSON(http.StatusCreated, web.OKResponse(core.StocksReservedCode, nil, nil))
}

// @Summary		Release Stock
// @Description	Release Stock
// @Tags			Stocks
// @Produce		json
// @Param			stocks	body		model.ChangeAmountRequest	true	"List of goods with amount to release"
// @Success		200	{object}	model.ReleaseStockResponse
// @Failure		400	{object}	model.BadRequestInvalidBodyResponse
// @Failure		404	{object}	model.GoodNotFoundResponse
// @Failure		422	{object}	model.ValidationResponse
// @Failure		500	{object}	model.InternalResponse
// @Router		/stocks/release [put]
func (r *Resolver) releaseStock(c echo.Context) error {
	var model hub.ChangeStocksAmountModel
	if err := webutil.BodyChecker(c, &model); err != nil {
		return err
	}

	if err := r.service.ReleaseStocks(c.Request().Context(), model); err != nil {
		switch err {
		case sql.ErrNoRows:
			return c.JSON(http.StatusNotFound, web.ErrorResponse(core.GoodNotFoundCode, nil, nil))
		case core.ErrNotEnoughReserve:
			return c.JSON(http.StatusBadRequest, web.ErrorResponse(core.NotEnoughReserveCode, nil, nil))
		default:
			return c.JSON(http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, err.Error(), nil))
		}
	}

	return c.JSON(http.StatusCreated, web.OKResponse(core.StocksReleasedCode, nil, nil))
}
