package transaction

// Request DTOs
type CreateTransactionItemRequest struct {
	VehicleID int64 `json:"vehicle_id" validate:"required,gt=0"`
	Quantity  int   `json:"quantity" validate:"required,gt=0"`
}

type CreateTransactionRequest struct {
	PaymentMethod string                         `json:"payment_method" validate:"required,oneof=bank_transfer qris gopay"`
	Bank          string                         `json:"bank,omitempty" validate:"omitempty,oneof=bca bni bri permata"`
	Items         []CreateTransactionItemRequest `json:"items" validate:"required,min=1,dive"`
}

// Response DTOs
type PaymentResponseFromGateway struct {
	OrderID               string  `json:"order_id"`
	MidtransTransactionID string  `json:"midtrans_transaction_id"`
	PaymentMethod         string  `json:"payment_method"`
	VaNumber              *string `json:"va_number,omitempty"`
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
