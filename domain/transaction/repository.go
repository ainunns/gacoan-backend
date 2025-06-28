package transaction

import "context"

type Repository interface {
	CreateTransaction(ctx context.Context, tx interface{}, transactionEntity Transaction) (Transaction, error)
	GetAllTransactions(ctx context.Context, tx interface{}) ([]Transaction, error)
	GetTransactionByID(ctx context.Context, tx interface{}, id string) (Transaction, error)
	GetLatestQueueCode(ctx context.Context, tx interface{}) (QueueCode, error)
}
