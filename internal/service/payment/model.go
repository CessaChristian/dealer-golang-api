package payment

import "time"

type PaymentDetail struct {
	PaymentID             int64      `db:"payment_id"`
	TransactionID         int64      `db:"transaction_id"`
	MidtransTransactionID string     `db:"midtrans_transaction_id"`
	PaymentType           string     `db:"payment_type"`
	TransactionStatus     string     `db:"transaction_status"`
	FraudStatus           string     `db:"fraud_status"`
	GrossAmount           float64    `db:"gross_amount"`
	PaidAt                *time.Time `db:"paid_at"`
	Note                  *string    `db:"note"`
}
