package user

import (
	"dealer_golang_api/utils"


	"github.com/labstack/echo/v4"
)

type Controller struct {
	Svc *Service
}

func NewController(s *Service) *Controller {
	return &Controller{Svc: s}
}

func (c *Controller) GetAll(ctx echo.Context) error {
	users, err := c.Svc.GetAllUsers()
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, "success", users)
}

func (c *Controller) GetByID(ctx echo.Context) error {
	id := utils.ToInt(ctx.Param("id"))

	user, err := c.Svc.GetByID(id)
	if err != nil {
		return utils.Error(ctx, 404, "user not found")
	}

	return utils.Success(ctx, "success", user)
}
