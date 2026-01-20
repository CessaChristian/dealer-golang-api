package transaction

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"dealer_golang_api/internal/service/vehicle"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PaymentService interface: Charge returns typed PaymentResponseFromGateway
type PaymentService interface {
	Charge(ctx context.Context, orderID string, amount float64, paymentMethod string, opts map[string]interface{}) (PaymentResponseFromGateway, error)
}

type Service struct {
	repo        Repository
	db          *pgxpool.Pool
	vehicleRepo vehicle.VehicleRepository
	paymentSvc  PaymentService
}

func NewService(db *pgxpool.Pool, repo Repository, vRepo vehicle.VehicleRepository, pSvc PaymentService) *Service {
	return &Service{
		repo:        repo,
		db:          db,
		vehicleRepo: vRepo,
		paymentSvc:  pSvc,
	}
}

func CreateOrderID() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	r := hex.EncodeToString(b)
	now := time.Now().UTC().Format("20060102T150405")
	return fmt.Sprintf("TX-%s-%s", now, strings.ToUpper(r))
}

func (s *Service) CreateTransaction(ctx context.Context, userID int64, req CreateTransactionRequest) (CreateTransactionResponse, error) {
	if len(req.Items) == 0 {
		return CreateTransactionResponse{}, errors.New("no items provided")
	}

	var total float64
	var items []TransactionItem

	for _, it := range req.Items {
		veh, err := s.vehicleRepo.GetByID(ctx, int(it.VehicleID))
		if err != nil {
			return CreateTransactionResponse{}, fmt.Errorf("vehicle %d not found", it.VehicleID)
		}
		if veh.Stock < it.Quantity {
			return CreateTransactionResponse{}, fmt.Errorf("vehicle %d has insufficient stock", it.VehicleID)
		}
		price := veh.Price
		total += price * float64(it.Quantity)
		items = append(items, TransactionItem{
			VehicleID: it.VehicleID,
			Quantity:  it.Quantity,
			Price:     price,
		})
	}

	orderID := CreateOrderID()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return CreateTransactionResponse{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	t := Transaction{
		OrderID:       orderID,
		UserID:        userID,
		TotalAmount:   total,
		PaymentMethod: req.PaymentMethod,
		Bank:          nil,
		VaNumber:      nil,
		Status:        "pending",
	}

	tid, err := s.repo.CreateTransaction(ctx, tx, t)
	if err != nil {
		return CreateTransactionResponse{}, err
	}

	for i := range items {
		items[i].TransactionID = tid
	}

	if err := s.repo.CreateTransactionItems(ctx, tx, items); err != nil {
		return CreateTransactionResponse{}, err
	}

	opts := map[string]interface{}{}
	if req.Bank != "" {
		opts["bank"] = req.Bank
	}

	payResp, err := s.paymentSvc.Charge(ctx, orderID, total, req.PaymentMethod, opts)
	if err != nil {
		return CreateTransactionResponse{}, err
	}

	var bankPtr *string
	var vaPtr *string
	if req.Bank != "" {
		bankPtr = &req.Bank
	}
	if payResp.VaNumber != nil {
		vaPtr = payResp.VaNumber
	}

	// update transaction with bank and single VA (string) if available
	if err := s.repo.UpdateTransactionPaymentInfo(ctx, tx, orderID, bankPtr, vaPtr); err != nil {
		return CreateTransactionResponse{}, err
	}

	// store payment detail (no JSON, only primitive types)
	pd := PaymentDetail{
		TransactionID:         tid,
		MidtransTransactionID: &payResp.MidtransTransactionID,
		PaymentType:           &payResp.PaymentMethod,
		TransactionStatus:     ptrString("pending"),
		GrossAmount:           &payResp.Amount,
	}
	if err := s.repo.InsertPaymentDetail(ctx, tx, pd); err != nil {
		return CreateTransactionResponse{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return CreateTransactionResponse{}, err
	}

	resp := CreateTransactionResponse{
		OrderID:       orderID,
		PaymentMethod: req.PaymentMethod,
		Amount:        total,
		Status:        "pending",
		VaNumber:      payResp.VaNumber,
		QRString:      payResp.QRString,
	}

	return resp, nil
}

func ptrString(s string) *string { return &s }

func (s *Service) GetTransaction(ctx context.Context, orderID string) (Transaction, []TransactionItem, error) {
	return s.repo.GetByOrderID(ctx, orderID)
}
