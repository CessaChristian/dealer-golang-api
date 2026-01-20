package transaction

// Request DTOs
type CreateTransactionItemRequest struct {
	VehicleID int64 `json:"vehicle_id"`
	Quantity  int   `json:"quantity"`
}

type CreateTransactionRequest struct {
	PaymentMethod string                         `json:"payment_method"`
	Bank          string                         `json:"bank,omitempty"`
	Items         []CreateTransactionItemRequest `json:"items"`
}

// Response DTOs (typed; no actions, no va_numbers array)
type PaymentResponseFromGateway struct {
	OrderID               string  `json:"order_id"`
	MidtransTransactionID string  `json:"midtrans_transaction_id"`
	PaymentMethod         string  `json:"payment_method"`
	VaNumber              *string `json:"va_number,omitempty"` // single VA for DB
	QRString              *string `json:"qr_string,omitempty"`
	Amount                float64 `json:"amount"`
}

type CreateTransactionResponse struct {
	OrderID       string  `json:"order_id"`
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
	VaNumber      *string `json:"va_number,omitempty"`
	QRString      *string `json:"qr_string,omitempty"`
}
