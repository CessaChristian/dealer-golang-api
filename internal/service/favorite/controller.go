package favorite

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

func (c *Controller) AddFavorite(ctx echo.Context) error {
	userID := ctx.Get("user_id").(int)

	var req AddFavoriteRequest
	if err := ctx.Bind(&req); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, "invalid request")
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, err.Error())
	}

	if err := c.Svc.AddFavorite(userID, req.VehicleID); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	return utils.Success(ctx, "favorite added successfully", nil)
}

func (c *Controller) GetFavorites(ctx echo.Context) error {
	userID := ctx.Get("user_id").(int)

	data, err := c.Svc.GetFavorites(userID)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, "success", data)
}

func (c *Controller) GetAllFavoritesAdmin(ctx echo.Context) error {
	data, err := c.Svc.GetAllFavoritesAdmin()
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, "success", data)
}
