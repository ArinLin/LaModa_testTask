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
		case core.ErrStockNotFound:
			return c.JSON(http.StatusNotFound, web.ErrorResponse(core.StockNotFoundCode, nil, nil))
		default:
			return c.JSON(http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, err.Error(), nil))
		}
	}

	pres := stock.PresentList(data)

	return c.JSON(http.StatusOK, pres.Response())
}

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
