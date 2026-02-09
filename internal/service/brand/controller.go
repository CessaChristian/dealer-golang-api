package brand

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
	var req CreateBrandRequest

	if err := ctx.Bind(&req); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, "invalid request")
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, err.Error())
	}

	if err := c.Svc.Create(req.Name); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	return utils.Success(ctx, "brand created", nil)
}

func (c *Controller) GetAll(ctx echo.Context) error {
	brands, err := c.Svc.GetAll()
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, "success", brands)
}
