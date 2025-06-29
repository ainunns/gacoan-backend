package transaction

import (
	"context"
	"fp-kpl/platform/pagination"
)

type Repository interface {
	CreateTransaction(ctx context.Context, tx interface{}, transactionEntity Transaction) (Transaction, error)
	GetAllTransactionsWithPagination(ctx context.Context, tx interface{}, userID string, req pagination.Request) (pagination.ResponseWithData, error)
	GetTransactionByID(ctx context.Context, tx interface{}, id string) (Transaction, error)
	GetLatestQueueCode(ctx context.Context, tx interface{}, id string) (string, error)
	UpdateTransaction(ctx context.Context, tx interface{}, transactionEntity Transaction) (Transaction, error)
}
