package transaction

import (
	"context"
	"fp-kpl/application/response"
	"fp-kpl/platform/pagination"
)

type Repository interface {
	CreateTransaction(ctx context.Context, tx interface{}, transactionEntity Transaction) (Transaction, error)
	GetAllTransactionsWithPagination(ctx context.Context, tx interface{}, userID string, req pagination.Request) (pagination.ResponseWithData, error)
	GetAllReadyToServeTransactionList(ctx context.Context, tx interface{}, req pagination.Request) (pagination.ResponseWithData, error)
	GetDetailedTransactionByID(ctx context.Context, tx interface{}, id string) (Query, error)
	GetLatestQueueCode(ctx context.Context, tx interface{}, id string) (string, error)
	GetNextOrder(ctx context.Context, tx interface{}) (response.NextOrder, error)
	UpdateCookedAt(ctx context.Context, tx interface{}, transactionID string) (Transaction, error)
	UpdateTransactionCookingStatusStart(ctx context.Context, tx interface{}, transactionID string) (Transaction, error)
	UpdateTransactionCookingStatusFinish(ctx context.Context, tx interface{}, transactionID string) (Transaction, error)
	UpdateTransactionDeliveringStatusStart(ctx context.Context, tx interface{}, transactionID string) (Transaction, error)
	UpdateTransactionDeliveringStatusFinish(ctx context.Context, tx interface{}, transactionID string) (Transaction, error)
	UpdateServedAt(ctx context.Context, tx interface{}, transactionID string) (Transaction, error)
	GetTransactionByQueueCode(ctx context.Context, tx interface{}, queueCode string) (Query, error)
}
