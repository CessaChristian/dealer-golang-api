package auth

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

func (c *Controller) Register(ctx echo.Context) error {
	var req RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, "invalid request")
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, err.Error())
	}

	if err := c.Svc.Register(req); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	return utils.Success(ctx, "register success", nil)
}

func (c *Controller) Login(ctx echo.Context) error {
	var req LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, "invalid request")
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, err.Error())
	}

	token, user, err := c.Svc.Login(req)
	if err != nil {
		return utils.Error(ctx, 401, err.Error())
	}

	return utils.Success(ctx, "login success", LoginResponse{
		Token: token,
		User:  user,
	})
}
