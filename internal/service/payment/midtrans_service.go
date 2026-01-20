package payment

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"dealer_golang_api/internal/service/transaction"
)

type MidtransService struct {
	ServerKey  string
	HttpClient *http.Client
}

func NewMidtransService() *MidtransService {
	return &MidtransService{
		ServerKey:  os.Getenv("MIDTRANS_SERVER_KEY"),
		HttpClient: &http.Client{},
	}
}

func basicAuth(serverKey string) string {
	token := serverKey + ":"
	return base64.StdEncoding.EncodeToString([]byte(token))
}

// Extracts the first va_number string from the raw response (if exists)
func extractVaNumber(raw interface{}) *string {
	if raw == nil {
		return nil
	}
	arr, ok := raw.([]interface{})
	if !ok || len(arr) == 0 {
		return nil
	}
	first := arr[0]
	m, ok := first.(map[string]interface{})
	if !ok {
		return nil
	}
	if v, ok := m["va_number"].(string); ok && v != "" {
		return &v
	}
	return nil
}

func (s *MidtransService) Charge(ctx context.Context, orderID string, amount float64, method string, opts map[string]interface{}) (transaction.PaymentResponseFromGateway, error) {
	url := "https://api.sandbox.midtrans.com/v2/charge"

	payload := map[string]interface{}{
		"transaction_details": map[string]interface{}{
			"order_id":     orderID,
			"gross_amount": amount,
		},
		"payment_type": method,
	}

	if method == "bank_transfer" {
		payload["bank_transfer"] = map[string]interface{}{
			"bank": opts["bank"],
		}
	}
	if method == "qris" {
		payload["qris"] = map[string]interface{}{}
	}
	if method == "gopay" {
		payload["gopay"] = map[string]interface{}{
			"enable_callback": true,
			"callback_url":    "https://example.com/callback",
		}
	}

	jsonBody, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+basicAuth(s.ServerKey))

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return transaction.PaymentResponseFromGateway{}, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var raw map[string]interface{}
	_ = json.Unmarshal(bodyBytes, &raw)

	var qr *string
	if v, ok := raw["qr_string"].(string); ok {
		qr = &v
	}

	txID := ""
	if v, ok := raw["transaction_id"].(string); ok {
		txID = v
	}

	// extract FIRST va_number (if present)
	vaSingle := extractVaNumber(raw["va_numbers"])

	return transaction.PaymentResponseFromGateway{
		OrderID:               orderID,
		MidtransTransactionID: txID,
		PaymentMethod:         method,
		VaNumber:              vaSingle,
		QRString:              qr,
		Amount:                amount,
	}, nil
}

func (s *MidtransService) ValidateSignature(orderID, statusCode, grossAmount, signature string) bool {
	raw := orderID + statusCode + grossAmount + s.ServerKey
	hash := fmt.Sprintf("%x", sha512.Sum512([]byte(raw)))
	return hash == signature
}
