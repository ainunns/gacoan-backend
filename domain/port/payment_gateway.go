package port

import (
	"context"
	"fp-kpl/domain/transaction"
)

type (
	PaymentGatewayPort interface {
		ProcessPayment(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (string, error)
	}
)
