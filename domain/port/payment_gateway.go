package port

import (
	"context"
	"fp-kpl/domain/transaction"
	"github.com/google/uuid"
)

type (
	PaymentGatewayPort interface {
		ProcessPayment(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (ProcessPaymentResponse, error)
		HookPayment(ctx context.Context, tx interface{}, transactionId uuid.UUID, datas map[string]interface{}) error
	}

	ProcessPaymentResponse struct {
		Token       string
		PaymentLink string
	}
)
