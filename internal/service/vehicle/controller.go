package vehicle

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

func (c *Controller) Create(ctx echo.Context) error {
	var req CreateVehicleRequest
	if err := ctx.Bind(&req); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, "invalid request")
	}

	if err := c.Svc.Create(req); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	return utils.Success(ctx, "vehicle created", nil)
}

func (c *Controller) Update(ctx echo.Context) error {
	id := utils.ToInt(ctx.Param("id"))

	var req UpdateVehicleRequest
	ctx.Bind(&req)

	if err := c.Svc.Update(id, req); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	return utils.Success(ctx, "vehicle updated", nil)
}

func (c *Controller) Import(ctx echo.Context) error {
	var req ImportVehicleRequest
	ctx.Bind(&req)

	if err := c.Svc.ImportCar(req); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	return utils.Success(ctx, "vehicle imported", nil)
}

func (c *Controller) GetAll(ctx echo.Context) error {
	list, err := c.Svc.GetAll()
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, "success", list)
}

func (c *Controller) GetByID(ctx echo.Context) error {
	id := utils.ToInt(ctx.Param("id"))
	v, err := c.Svc.GetByID(id)
	if err != nil {
		return utils.Error(ctx, 404, "vehicle not found")
	}

	return utils.Success(ctx, "success", v)
}

func (c *Controller) LowStock(ctx echo.Context) error {
	list, err := c.Svc.LowStock()
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, "success", list)
}
