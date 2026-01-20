package report

import (
	"dealer_golang_api/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	Svc *Service
}

func NewController(s *Service) *Controller {
	return &Controller{Svc: s}
}

// ========== LOW STOCK ==========

func (c *Controller) LowStockJSON(ctx echo.Context) error {
	data, err := c.Svc.LowStockJSON(ctx.Request().Context())
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, "success", data)
}

func (c *Controller) LowStockCSV(ctx echo.Context) error {
	csvText, err := c.Svc.LowStockCSV(ctx.Request().Context())
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	ctx.Response().Header().Set("Content-Type", "text/csv")
	ctx.Response().Header().Set("Content-Disposition", "attachment; filename=low_stock_report.csv")

	return ctx.String(http.StatusOK, csvText)
}

// ========== FAVORITES ==========

func (c *Controller) FavoriteJSON(ctx echo.Context) error {
	data, err := c.Svc.FavoriteJSON(ctx.Request().Context())
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, "success", data)
}

func (c *Controller) FavoriteCSV(ctx echo.Context) error {
	csvText, err := c.Svc.FavoriteCSV(ctx.Request().Context())
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	ctx.Response().Header().Set("Content-Type", "text/csv")
	ctx.Response().Header().Set("Content-Disposition", "attachment; filename=favorite_report.csv")

	return ctx.String(http.StatusOK, csvText)
}
