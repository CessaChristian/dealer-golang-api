package transaction

import "time"

type Transaction struct {
	ID            int64     `db:"transaction_id"`
	OrderID       string    `db:"order_id"`
	UserID        int64     `db:"user_id"`
	TotalAmount   float64   `db:"total_amount"`
	PaymentMethod string    `db:"payment_method"`
	Bank          *string   `db:"bank"`
	VaNumber      *string   `db:"va_number"`
	Status        string    `db:"status"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type TransactionItem struct {
	DetailID      int64   `db:"detail_id"`
	TransactionID int64   `db:"transaction_id"`
	VehicleID     int64   `db:"vehicle_id"`
	Quantity      int     `db:"quantity"`
	Price         float64 `db:"price"`
}

type PaymentDetail struct {
	PaymentID             int64       `db:"payment_id"`
	TransactionID         int64       `db:"transaction_id"`
	MidtransTransactionID *string     `db:"midtrans_transaction_id"`
	PaymentType           *string     `db:"payment_type"`
	TransactionStatus     *string     `db:"transaction_status"`
	FraudStatus           *string     `db:"fraud_status"`
	GrossAmount           *float64    `db:"gross_amount"`
	PaidAt                *time.Time  `db:"paid_at"`
	Note                  *string     `db:"note"`
}
