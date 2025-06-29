package repository

import (
	"context"
	"errors"
	"fmt"
	"fp-kpl/domain/transaction"
	"fp-kpl/infrastructure/database/db_transaction"
	"fp-kpl/infrastructure/database/schema"
	"fp-kpl/infrastructure/database/validation"
	"fp-kpl/platform/pagination"
	"time"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *db_transaction.Repository
}

func NewTransactionRepository(db *db_transaction.Repository) transaction.Repository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (transaction.Transaction, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return transaction.Transaction{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	transactionSchema := schema.TransactionEntityToSchema(transactionEntity)
	if err = db.WithContext(ctx).Create(&transactionSchema).Error; err != nil {
		return transaction.Transaction{}, err
	}

	transactionEntity = schema.TransactionSchemaToEntity(transactionSchema)
	return transactionEntity, nil
}

func (r *transactionRepository) GetAllTransactionsWithPagination(ctx context.Context, tx interface{}, userID string, req pagination.Request) (pagination.ResponseWithData, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return pagination.ResponseWithData{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var transactionSchemas []schema.Transaction
	var count int64

	req.Default()

	query := db.WithContext(ctx).Model(&transactionSchemas).Where("user_id = ?", userID)

	if req.Search != "" {
		query = query.Where("queue_code LIKE ?", "%"+req.Search+"%")
	}

	if err = query.Count(&count).Error; err != nil {
		return pagination.ResponseWithData{}, err
	}

	if err = query.Scopes(pagination.Paginate(req)).
		Preload("Table").
		Preload("Orders").
		Preload("Orders.Menu").
		Find(&transactionSchemas).Error; err != nil {
		return pagination.ResponseWithData{}, err
	}

	totalPage := pagination.TotalPage(count, int64(req.PerPage))

	// Return schema data directly to avoid N+1 queries
	data := make([]any, len(transactionSchemas))
	for i, transactionSchema := range transactionSchemas {
		data[i] = transactionSchema
	}

	return pagination.ResponseWithData{
		Data: data,
		Response: pagination.Response{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   count,
			MaxPage: totalPage,
		},
	}, nil
}

func (r *transactionRepository) GetTransactionByID(ctx context.Context, tx interface{}, userID string, id string) (transaction.Transaction, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return transaction.Transaction{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var transactionSchema schema.Transaction

	query := db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID)

	if err = query.Preload("Table").
		Preload("Orders").
		Preload("Orders.Menu").
		Take(&transactionSchema).Error; err != nil {
		return transaction.Transaction{}, err
	}

	transactionEntity := schema.TransactionSchemaToEntity(transactionSchema)
	return transactionEntity, nil
}

func (r *transactionRepository) GetLatestQueueCode(ctx context.Context, tx interface{}, id string) (string, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return "", err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var transactionSchema schema.Transaction
	today := time.Now().Format("2006-01-02")
	if err = db.WithContext(ctx).
		Where("DATE(created_at) = ?", today).
		Where("queue_code IS NOT NULL").
		Where("id != ?", id).
		Order("queue_code DESC").
		First(&transactionSchema).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "Q0001", nil
		}
		return "", err
	}

	latestQueueCode := transactionSchema.QueueCode
	num := 1
	if latestQueueCode != nil && len(*latestQueueCode) == 5 {
		fmt.Sscanf(*latestQueueCode, "Q%04d", &num)
		num++
	}

	return fmt.Sprintf("Q%04d", num), nil
}

func (r *transactionRepository) UpdateTransaction(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (transaction.Transaction, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return transaction.Transaction{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	transactionSchema := schema.TransactionEntityToSchema(transactionEntity)
	if err = db.WithContext(ctx).Where("id = ?", transactionEntity.ID).Updates(&transactionSchema).Error; err != nil {
		return transaction.Transaction{}, err
	}

	transactionEntity = schema.TransactionSchemaToEntity(transactionSchema)
	return transactionEntity, nil
}
