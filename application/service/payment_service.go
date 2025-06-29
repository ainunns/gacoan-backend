package service

import (
	"context"
	"fmt"
	"fp-kpl/application/request"
	"fp-kpl/application/response"
	"fp-kpl/domain/transaction"
	"fp-kpl/infrastructure/database/validation"
	"os"
)

type (
	PaymentService interface {
		PaymentNotification(ctx context.Context, req request.MidtransNotification) (response.TransactionDetail, error)
	}

	paymentService struct {
		transactionRepository transaction.Repository
		transaction           interface{}
	}
)

func NewPaymentService(
	transactionRepository transaction.Repository,
	transaction interface{},
) PaymentService {
	return &paymentService{
		transactionRepository: transactionRepository,
		transaction:           transaction,
	}
}

func (s *paymentService) PaymentNotification(ctx context.Context, req request.MidtransNotification) (response.TransactionDetail, error) {
	validatedTransaction, err := validation.ValidateTransaction(s.transaction)
	if err != nil {
		return response.TransactionDetail{}, err
	}

	tx, err := validatedTransaction.Begin(ctx)
	if err != nil {
		return response.TransactionDetail{}, err
	}

	signature := fmt.Sprintf("%s%s%s%s", req.OrderID, req.StatusCode, req.GrossAmount, os.Getenv("MIDTRANS_SERVER_KEY"))
	if signature != req.SignatureKey {
		return response.TransactionDetail{}, fmt.Errorf("invalid signature key")
	}

	transactionEntity, err := s.transactionRepository.GetTransactionByID(ctx, tx, req.OrderID)
	if err != nil {
		return response.TransactionDetail{}, err
	}

}
