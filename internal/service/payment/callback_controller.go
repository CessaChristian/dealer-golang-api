package payment

import (
	"net/http"
	"strconv"
	"time"

	"dealer_golang_api/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type CallbackController struct {
	DB     *pgxpool.Pool
	Repo   Repository
	PaySvc *MidtransService
}

func NewCallbackController(db *pgxpool.Pool, repo Repository, svc *MidtransService) *CallbackController {
	return &CallbackController{
		DB:     db,
		Repo:   repo,
		PaySvc: svc,
	}
}

// HandleMidtransCallback godoc
// @Summary Midtrans payment callback
// @Description Endpoint for Midtrans server to notify payment status (DO NOT CALL MANUALLY)
// @Tags Payment
// @Accept json
// @Produce json
// @Param payload body MidtransCallbackRequest true "Midtrans callback payload"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/payments/callback [post]
func (c *CallbackController) HandleCallback(ctx echo.Context) error {
	var req MidtransCallbackRequest
	if err := ctx.Bind(&req); err != nil {
		return utils.Error(ctx, http.StatusBadRequest, "invalid payload")
	}

	// 1️⃣ Validate signature (WAJIB)
	if !c.PaySvc.ValidateSignature(
		req.OrderID,
		req.StatusCode,     
		req.GrossAmount,
		req.SignatureKey,
	) {
		return utils.Error(ctx, http.StatusBadRequest, "invalid signature")
	}

	// 2️⃣ Begin DB transaction
	tx, err := c.DB.Begin(ctx.Request().Context())
	if err != nil {
		return utils.Error(ctx, http.StatusInternalServerError, "db error")
	}
	defer tx.Rollback(ctx.Request().Context())

	// 3️⃣ Ambil status transaksi SAAT INI
	currentStatus, err := c.Repo.GetTransactionStatus(
		ctx.Request().Context(),
		tx,
		req.OrderID,
	)
	if err != nil {
		return utils.Error(ctx, http.StatusNotFound, "transaction not found")
	}

	// 4️⃣ IDEMPOTENCY CHECK (INI KUNCI UTAMA)
	if isFinalStatus(currentStatus) {
		// Callback duplicate → ACK saja
		return utils.Success(ctx, "callback already processed", nil)
	}

	// 5️⃣ Ambil transaction_id (internal)
	trxID, err := c.Repo.GetTransactionByOrderID(
		ctx.Request().Context(),
		tx,
		req.OrderID,
	)
	if err != nil {
		return utils.Error(ctx, http.StatusNotFound, "transaction not found")
	}

	// 6️⃣ Insert payment detail (HANYA JIKA STATUS BERUBAH)
	if currentStatus != req.TransactionStatus {
		gross, _ := strconv.ParseFloat(req.GrossAmount, 64)
		now := time.Now()

		pd := PaymentDetail{
			TransactionID:         trxID,
			MidtransTransactionID: req.TransactionID,
			PaymentType:           req.PaymentType,
			TransactionStatus:     req.TransactionStatus,
			FraudStatus:           req.FraudStatus,
			GrossAmount:           gross,
			PaidAt:                &now,
		}

		if err := c.Repo.InsertPaymentDetail(
			ctx.Request().Context(),
			tx,
			pd,
		); err != nil {
			return utils.Error(ctx, 500, err.Error())
		}
	}

	// 7️⃣ Update status transaksi
	if err := c.Repo.UpdateTransactionStatus(
		ctx.Request().Context(),
		tx,
		req.OrderID,
		req.TransactionStatus,
	); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	// 8️⃣ Kurangi stok HANYA SEKALI (settlement pertama)
	if req.TransactionStatus == "settlement" && currentStatus != "settlement" {
		items, err := c.Repo.GetTransactionItems(
			ctx.Request().Context(),
			tx,
			trxID,
		)
		if err != nil {
			return utils.Error(ctx, 500, "failed to load transaction items")
		}

		for _, it := range items {
			if err := c.Repo.DecreaseVehicleStock(
				ctx.Request().Context(),
				tx,
				it.VehicleID,
				it.Quantity,
			); err != nil {
				return utils.Error(ctx, 500, "failed to update stock")
			}
		}
	}

	// 9️⃣ Commit DB
	if err := tx.Commit(ctx.Request().Context()); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	// 🔟 ACK ke Midtrans (WAJIB 200)
	return utils.Success(ctx, "callback processed", nil)
}


func isFinalStatus(status string) bool {
	switch status {
	case "settlement", "expire", "cancel", "deny":
		return true
	default:
		return false
	}
}

