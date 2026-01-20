package transaction

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateTransaction(ctx context.Context, tx pgx.Tx, t Transaction) (int64, error)
	CreateTransactionItems(ctx context.Context, tx pgx.Tx, items []TransactionItem) error
	GetByOrderID(ctx context.Context, orderID string) (Transaction, []TransactionItem, error)
	UpdateTransactionPaymentInfo(ctx context.Context, tx pgx.Tx, orderID string, bank, va *string) error
	InsertPaymentDetail(ctx context.Context, tx pgx.Tx, pd PaymentDetail) error
	UpdateTransactionStatus(ctx context.Context, tx pgx.Tx, orderID, status string) error
	DecreaseVehicleStock(ctx context.Context, tx pgx.Tx, vehicleID int64, qty int) error
	GetTransactionIDByOrderID(ctx context.Context, tx pgx.Tx, orderID string) (int64, error)
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repo{db: db}
}

func (r *repo) CreateTransaction(ctx context.Context, tx pgx.Tx, t Transaction) (int64, error) {
	var id int64
	err := tx.QueryRow(ctx,
		`INSERT INTO transactions (order_id, user_id, total_amount, payment_method, bank, va_number, status)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING transaction_id`,
		t.OrderID, t.UserID, t.TotalAmount, t.PaymentMethod, t.Bank, t.VaNumber, t.Status,
	).Scan(&id)
	return id, err
}

func (r *repo) CreateTransactionItems(ctx context.Context, tx pgx.Tx, items []TransactionItem) error {
	for _, it := range items {
		_, err := tx.Exec(ctx, `INSERT INTO transaction_items (transaction_id, vehicle_id, quantity, price) VALUES ($1,$2,$3,$4)`,
			it.TransactionID, it.VehicleID, it.Quantity, it.Price)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *repo) GetByOrderID(ctx context.Context, orderID string) (Transaction, []TransactionItem, error) {
	var t Transaction
	err := r.db.QueryRow(ctx, `
		SELECT transaction_id, order_id, user_id, total_amount, payment_method, bank, va_number, status, created_at, updated_at
		FROM transactions WHERE order_id = $1`, orderID,
	).Scan(&t.ID, &t.OrderID, &t.UserID, &t.TotalAmount, &t.PaymentMethod, &t.Bank, &t.VaNumber, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return Transaction{}, nil, err
	}

	rows, err := r.db.Query(ctx, `SELECT detail_id, transaction_id, vehicle_id, quantity, price FROM transaction_items WHERE transaction_id = $1`, t.ID)
	if err != nil {
		return t, nil, err
	}
	defer rows.Close()

	var items []TransactionItem
	for rows.Next() {
		var it TransactionItem
		if err := rows.Scan(&it.DetailID, &it.TransactionID, &it.VehicleID, &it.Quantity, &it.Price); err != nil {
			return t, nil, err
		}
		items = append(items, it)
	}
	return t, items, nil
}

func (r *repo) UpdateTransactionPaymentInfo(ctx context.Context, tx pgx.Tx, orderID string, bank, va *string) error {
	_, err := tx.Exec(ctx, `UPDATE transactions SET bank = $1, va_number = $2, updated_at = CURRENT_TIMESTAMP WHERE order_id = $3`, bank, va, orderID)
	return err
}

func (r *repo) InsertPaymentDetail(ctx context.Context, tx pgx.Tx, pd PaymentDetail) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO payment_detail (
			transaction_id, midtrans_transaction_id, payment_type,
			transaction_status, fraud_status, gross_amount, paid_at, note
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, pd.TransactionID, pd.MidtransTransactionID, pd.PaymentType, pd.TransactionStatus, pd.FraudStatus, pd.GrossAmount, pd.PaidAt, pd.Note)

	return err
}

func (r *repo) UpdateTransactionStatus(ctx context.Context, tx pgx.Tx, orderID, status string) error {
	_, err := tx.Exec(ctx, `UPDATE transactions SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE order_id = $2`, status, orderID)
	return err
}

func (r *repo) DecreaseVehicleStock(ctx context.Context, tx pgx.Tx, vehicleID int64, qty int) error {
	res, err := tx.Exec(ctx, `UPDATE vehicles SET stock = stock - $1 WHERE vehicle_id = $2 AND stock >= $1`, qty, vehicleID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("insufficient stock for vehicle %d", vehicleID)
	}
	return nil
}

func (r *repo) GetTransactionIDByOrderID(ctx context.Context, tx pgx.Tx, orderID string) (int64, error) {
	var id int64
	err := tx.QueryRow(ctx, `SELECT transaction_id FROM transactions WHERE order_id=$1`, orderID).Scan(&id)
	return id, err
}
