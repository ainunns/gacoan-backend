package order

import "context"

type (
	Repository interface {
		CreateOrder(ctx context.Context, tx interface{}, orderEntity Order) (Order, error)
		GetOrderByTransactionID(ctx context.Context, tx interface{}, transactionID string) ([]Order, error)
	}
)
