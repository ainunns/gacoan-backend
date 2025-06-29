package transaction

import (
	"context"
	"fp-kpl/platform/pagination"
)

type Repository interface {
	CreateTransaction(ctx context.Context, tx interface{}, transactionEntity Transaction) (Transaction, error)
	GetAllTransactionsWithPagination(ctx context.Context, tx interface{}, userID string, req pagination.Request) (pagination.ResponseWithData, error)
	GetTransactionByID(ctx context.Context, tx interface{}, userID string, id string) (interface{}, error)
	GetLatestQueueCode(ctx context.Context, tx interface{}, id string) (string, error)
	GetNextOrder(ctx context.Context, tx interface{}, userID string) (interface{}, error)
	UpdateTransactionCookingStatusStart(ctx context.Context, tx interface{}, transactionID string) (Transaction, error)
	GetTransactionByQueueCode(ctx context.Context, tx interface{}, queueCode string) (interface{}, error)
}
