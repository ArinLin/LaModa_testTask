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

// @Summary		Get Goods
// @Description	Get Goods List
// @Tags			Goods
// @Produce		json
// @Success		200	{object}	model.GetGoodsResponse
// @Failure		500	{object}	model.InternalResponse
// @Router			/goods [get]
func (r *Resolver) getGoods(c echo.Context) error {
	data, err := r.service.GetGoods(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, web.ErrorResponse(core.InternalErrorCode, err.Error(), nil))
	}

	pres := good.PresentList(data)

	return c.JSON(http.StatusOK, pres.Response())
}

// @Summary		Get Good By ID
// @Description	Get Goods By ID
// @Tags			Goods
// @Produce		json
// @Param			id	path		int	true	"Good ID"
// @Success		200	{object}	model.GetGoodByIDResponse
// @Failure		400	{object}	model.BadRequestInvalidIDResponse
// @Failure		404	{object}	model.GoodNotFoundResponse
// @Failure		500	{object}	model.InternalResponse
// @Router			/good/{id} [get]
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

// @Summary		Create Good
// @Description	Create Good
// @Tags			Goods
// @Produce		json
// @Param			good	body		model.CreateGoodRequest	true	"Params to create good"
// @Success		200		{object}	model.GetGoodByIDResponse
// @Failure		400		{object}	model.BadRequestInvalidBodyResponse
// @Failure		422		{object}	model.ValidationResponse
// @Failure		500		{object}	model.InternalResponse
// @Router			/good [post]
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

// @Summary		Update Good By ID
// @Description	Update Good By ID
// @Tags			Goods
// @Produce		json
// @Param			id		path		int					true	"Good ID"
// @Param			good	body		model.UpdateGoodRequest	true	"Params to update good"
// @Success		200		{object}	model.GetGoodByIDResponse
// @Failure		400		{object}	model.BadRequestInvalidIDResponse
// @Failure		404		{object}	model.GoodNotFoundResponse
// @Failure		422		{object}	model.ValidationResponse
// @Failure		500		{object}	model.InternalResponse
// @Router			/good/{id} [patch]
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

// @Summary		Delete Good
// @Description	Delete Good By ID
// @Tags			Goods
// @Produce		json
// @Param			id	path		int	true	"Good ID"
// @Success		200	{object}	model.DeleteGoodResponse
// @Failure		400	{object}	model.BadRequestInvalidIDResponse
// @Failure		404	{object}	model.GoodNotFoundResponse
// @Failure		500	{object}	model.InternalResponse
// @Router			/good/{id} [delete]
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
