package handlers

import (
	"database/sql"
	"net/http"

	"lamoda/api/hub/rest/presenters/good"
	"lamoda/internal/core"
	"lamoda/internal/service/hub"
	"lamoda/pkg/web"
	"lamoda/pkg/webutil"

	"github.com/labstack/echo/v4"
)

func (r *Resolver) getGoods(c echo.Context) error {
	data, err := r.service.GetGoods(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, err.Error(), nil))
	}

	pres := good.PresentList(data)

	return c.JSON(http.StatusOK, pres.Response())
}

func (r *Resolver) getGoodByID(c echo.Context) error {
	id, err := webutil.ParseID(c)
	if err != nil {
		return err
	}

	data, err := r.service.GetGoodByID(c.Request().Context(), id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return c.JSON(http.StatusNotFound, web.ErrorResponse(core.GoodNotFoundCode, nil, nil))
		default:
			return c.JSON(http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, err.Error(), nil))
		}
	}

	pres := good.PresentGood(data)

	return c.JSON(http.StatusOK, pres.Response(core.GoodReceivedCode))
}

func (r *Resolver) createGood(c echo.Context) error {
	var model hub.CreateGoodModel
	if err := webutil.BodyChecker(c, &model); err != nil {
		return err
	}

	data, err := r.service.CreateGood(c.Request().Context(), model)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, err.Error(), nil))
	}

	pres := good.PresentGood(data)

	return c.JSON(http.StatusCreated, pres.Response(core.GoodCreatedCode))
}

func (r *Resolver) updateGood(c echo.Context) error {
	id, err := webutil.ParseID(c)
	if err != nil {
		return err
	}

	var model hub.UpdateGoodModel
	if err := webutil.BodyChecker(c, &model); err != nil {
		return err
	}

	data, err := r.service.UpdateGood(c.Request().Context(), id, model)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return c.JSON(http.StatusNotFound, web.ErrorResponse(core.GoodNotFoundCode, nil, nil))
		default:
			return c.JSON(http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, err.Error(), nil))
		}
	}

	pres := good.PresentGood(data)

	return c.JSON(http.StatusOK, pres.Response(core.GoodUpdatedCode))
}

func (r *Resolver) deleteGood(c echo.Context) error {
	id, err := webutil.ParseID(c)
	if err != nil {
		return err
	}

	if err = r.service.DeleteGood(c.Request().Context(), id); err != nil {
		switch err {
		case sql.ErrNoRows:
			return c.JSON(http.StatusNotFound, web.ErrorResponse(core.GoodNotFoundCode, nil, nil))
		default:
			return c.JSON(http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, err.Error(), nil))
		}
	}

	return c.JSON(http.StatusOK, web.OKResponse(core.GoodDeletedCode, nil, nil))
}
