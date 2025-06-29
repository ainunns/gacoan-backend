package repository

import (
	"context"
	"errors"
	"fmt"
	"fp-kpl/domain/transaction"
	"fp-kpl/infrastructure/database/db_transaction"
	"fp-kpl/infrastructure/database/schema"
	"fp-kpl/infrastructure/database/validation"
	"gorm.io/gorm"
	"time"
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

func (r *transactionRepository) GetAllTransactions(ctx context.Context, tx interface{}) ([]transaction.Transaction, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return nil, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var transactionSchemas []schema.Transaction
	query := db.WithContext(ctx).Model(&transactionSchemas)
	if err = query.Find(&transactionSchemas).Error; err != nil {
		return nil, err
	}

	transactionEntities := make([]transaction.Transaction, len(transactionSchemas))
	for i, transactionSchema := range transactionSchemas {
		transactionEntities[i] = schema.TransactionSchemaToEntity(transactionSchema)
	}

	return transactionEntities, nil
}

func (r *transactionRepository) GetTransactionByID(ctx context.Context, tx interface{}, id string) (transaction.Transaction, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return transaction.Transaction{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var transactionSchema schema.Transaction

	if err = db.WithContext(ctx).Where("id = ?", id).Take(&transactionSchema).Error; err != nil {
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
