package transaction

import (
	"net/http"

	"dealer_golang_api/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	Svc *Service
	DB  *pgxpool.Pool
}

func NewController(svc *Service, db *pgxpool.Pool) *Controller {
	return &Controller{Svc: svc, DB: db}
}

func (c *Controller) CreateTransaction(ctx echo.Context) error {
	userIDIface := ctx.Get("user_id")
	if userIDIface == nil {
		return utils.Error(ctx, http.StatusUnauthorized, "unauthenticated")
	}
	userID := int64(userIDIface.(int))
	var req CreateTransactionRequest
	if err := ctx.Bind(&req); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, "invalid request")
	}

	resp, err := c.Svc.CreateTransaction(ctx.Request().Context(), userID, req)
	if err != nil {
		return utils.Error(ctx, http.StatusBadRequest, err.Error())
	}
	return utils.Success(ctx, "transaction created", resp)
}

func (c *Controller) GetTransaction(ctx echo.Context) error {
	orderID := ctx.Param("order_id")
	t, items, err := c.Svc.GetTransaction(ctx.Request().Context(), orderID)
	if err != nil {
		return utils.Error(ctx, 404, "transaction not found")
	}

	type itemResp struct {
		VehicleID int64   `json:"vehicle_id"`
		Quantity  int     `json:"quantity"`
		Price     float64 `json:"price"`
	}
	var itemsResp []itemResp
	for _, it := range items {
		itemsResp = append(itemsResp, itemResp{
			VehicleID: it.VehicleID,
			Quantity:  it.Quantity,
			Price:     it.Price,
		})
	}

	data := map[string]interface{}{
		"order_id":       t.OrderID,
		"status":         t.Status,
		"payment_method": t.PaymentMethod,
		"amount":         t.TotalAmount,
		"items":          itemsResp,
		"va_number":      t.VaNumber,
		"bank":           t.Bank,
	}

	return utils.Success(ctx, "success", data)
}
