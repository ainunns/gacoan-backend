package order

import "context"

type (
	Repository interface {
		CreateOrder(ctx context.Context, tx interface{}, orderEntity Order) (Order, error)
		GetAllOrder(ctx context.Context, tx interface{}) ([]Order, error)
		GetOrderByID(ctx context.Context, tx interface{}, id string) (Order, error)
		GetOrderByTransactionID(ctx context.Context, tx interface{}, transactionID string) ([]Order, error)
		UpdateOrder(ctx context.Context, tx interface{}, orderEntity Order) (Order, error)
	}
)
