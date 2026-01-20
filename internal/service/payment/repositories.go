package payment

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Repository interface {
	InsertPaymentDetail(ctx context.Context, tx pgx.Tx, d PaymentDetail) error
	UpdateTransactionStatus(ctx context.Context, tx pgx.Tx, orderID, status string) error
	DecreaseVehicleStock(ctx context.Context, tx pgx.Tx, vehicleID int64, qty int) error
	GetTransactionByOrderID(ctx context.Context, tx pgx.Tx, orderID string) (int64, error)
	GetTransactionItems(ctx context.Context, tx pgx.Tx, transactionID int64) ([]struct {
		VehicleID int64
		Quantity  int
	}, error)
	GetTransactionStatus(ctx context.Context, tx pgx.Tx, orderID string) (string, error)
}

type repo struct{}

func NewRepository() Repository { return &repo{} }

func (r *repo) InsertPaymentDetail(ctx context.Context, tx pgx.Tx, d PaymentDetail) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO payment_detail (
			transaction_id, midtrans_transaction_id, payment_type,
			transaction_status, fraud_status, gross_amount, paid_at, note
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, d.TransactionID, d.MidtransTransactionID, d.PaymentType, d.TransactionStatus, d.FraudStatus, d.GrossAmount, d.PaidAt, d.Note)
	return err
}

func (r *repo) UpdateTransactionStatus(ctx context.Context, tx pgx.Tx, orderID, status string) error {
	_, err := tx.Exec(ctx, `UPDATE transactions SET status=$1, updated_at=NOW() WHERE order_id=$2`, status, orderID)
	return err
}

func (r *repo) GetTransactionByOrderID(ctx context.Context, tx pgx.Tx, orderID string) (int64, error) {
	var id int64
	err := tx.QueryRow(ctx, `SELECT transaction_id FROM transactions WHERE order_id=$1`, orderID).Scan(&id)
	return id, err
}

func (r *repo) GetTransactionItems(ctx context.Context, tx pgx.Tx, transactionID int64) ([]struct {
	VehicleID int64
	Quantity  int
}, error) {
	rows, err := tx.Query(ctx, `SELECT vehicle_id, quantity FROM transaction_items WHERE transaction_id=$1`, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []struct {
		VehicleID int64
		Quantity  int
	}
	for rows.Next() {
		var rec struct {
			VehicleID int64
			Quantity  int
		}
		rows.Scan(&rec.VehicleID, &rec.Quantity)
		list = append(list, rec)
	}
	return list, nil
}

func (r *repo) DecreaseVehicleStock(ctx context.Context, tx pgx.Tx, vehicleID int64, qty int) error {
	_, err := tx.Exec(ctx, `UPDATE vehicles SET stock = stock - $1 WHERE vehicle_id=$2 AND stock >= $1`, qty, vehicleID)
	return err
}

func (r *repo) GetTransactionStatus(
	ctx context.Context,
	tx pgx.Tx,
	orderID string,
) (string, error) {
	var status string
	err := tx.QueryRow(ctx,
		`SELECT status FROM transactions WHERE order_id = $1`,
		orderID,
	).Scan(&status)
	return status, err
}
