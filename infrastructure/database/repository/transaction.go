package repository

import (
	"context"
	"errors"
	"fmt"
	"fp-kpl/application/response"
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

	transactionQueries := make([]transaction.Query, len(transactionSchemas))
	for i, transactionSchema := range transactionSchemas {
		transactionQueries[i].Transaction = schema.TransactionSchemaToEntity(transactionSchema)
		for j, orderSchema := range transactionSchema.Orders {
			transactionQueries[i].Orders = append(transactionQueries[i].Orders, transaction.OrderQuery{
				Order: schema.OrderSchemaToEntity(orderSchema),
			})
			transactionQueries[i].Orders[j].Menu = schema.MenuSchemaToEntity(*orderSchema.Menu)
		}
		transactionQueries[i].Table = schema.TableSchemaToEntity(*transactionSchema.Table)
	}

	totalPage := pagination.TotalPage(count, int64(req.PerPage))

	data := make([]any, len(transactionQueries))
	for i, transactionQuery := range transactionQueries {
		data[i] = transactionQuery
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

func (r *transactionRepository) GetAllReadyToServeTransactionList(ctx context.Context, tx interface{}, req pagination.Request) (pagination.ResponseWithData, error) {
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

	query := db.WithContext(ctx).Model(&transactionSchemas).
		Where("payment_status IN ?", []string{transaction.PaymentStatusSettlement, transaction.PaymentStatusCapture}).
		Where("served_at IS NULL").
		Where("order_status = ?", transaction.OrderStatusReadyToServe)

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
		Order("created_at DESC").
		Find(&transactionSchemas).Error; err != nil {
		return pagination.ResponseWithData{}, err
	}

	transactionQueries := make([]transaction.Query, len(transactionSchemas))
	for i, transactionSchema := range transactionSchemas {
		transactionQueries[i].Transaction = schema.TransactionSchemaToEntity(transactionSchema)
		for j, orderSchema := range transactionSchema.Orders {
			transactionQueries[i].Orders = append(transactionQueries[i].Orders, transaction.OrderQuery{
				Order: schema.OrderSchemaToEntity(orderSchema),
			})
			transactionQueries[i].Orders[j].Menu = schema.MenuSchemaToEntity(*orderSchema.Menu)
		}
		transactionQueries[i].Table = schema.TableSchemaToEntity(*transactionSchema.Table)
	}

	totalPage := pagination.TotalPage(count, int64(req.PerPage))

	data := make([]any, len(transactionQueries))
	for i, transactionQuery := range transactionQueries {
		data[i] = transactionQuery
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

func (r *transactionRepository) GetDetailedTransactionByID(ctx context.Context, tx interface{}, id string) (transaction.Query, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return transaction.Query{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var transactionSchema schema.Transaction

	query := db.WithContext(ctx).Where("id = ?", id)

	if err = query.Preload("Table").
		Preload("Orders").
		Preload("Orders.Menu").
		Take(&transactionSchema).Error; err != nil {
		return transaction.Query{}, err
	}

	var transactionQuery transaction.Query
	transactionQuery.Transaction = schema.TransactionSchemaToEntity(transactionSchema)
	for i, orderSchema := range transactionSchema.Orders {
		transactionQuery.Orders = append(transactionQuery.Orders, transaction.OrderQuery{
			Order: schema.OrderSchemaToEntity(orderSchema),
		})
		transactionQuery.Orders[i].Menu = schema.MenuSchemaToEntity(*orderSchema.Menu)
	}
	transactionQuery.Table = schema.TableSchemaToEntity(*transactionSchema.Table)

	return transactionQuery, nil
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

func (r *transactionRepository) GetNextOrder(ctx context.Context, tx interface{}) (response.NextOrder, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return response.NextOrder{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var transactionSchema schema.Transaction
	today := time.Now().Format("2006-01-02")

	query := db.WithContext(ctx).Where("payment_status IN ?", []string{transaction.PaymentStatusSettlement, transaction.PaymentStatusCapture})

	if err = query.Where("order_status = ?", "pending").
		Where("DATE(created_at) = ?", today).
		Preload("Table").
		Preload("Orders").
		Preload("Orders.Menu").
		Order("created_at ASC").First(&transactionSchema).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NextOrder{}, nil
		}
		return response.NextOrder{}, err
	}

	var orderResponses []response.OrderForTransaction
	for _, orderSchema := range transactionSchema.Orders {
		orderResponses = append(orderResponses, response.OrderForTransaction{
			Menu: response.MenuForTransaction{
				ID:    orderSchema.Menu.ID.String(),
				Name:  orderSchema.Menu.Name,
				Price: orderSchema.Menu.Price.String(),
			},
			Quantity: orderSchema.Quantity,
		})
	}

	return response.NextOrder{
		QueueCode: *transactionSchema.QueueCode,
		Orders:    orderResponses,
	}, nil
}

func (r *transactionRepository) GetTransactionByQueueCode(ctx context.Context, tx interface{}, queueCode string) (transaction.Query, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return transaction.Query{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var transactionSchema schema.Transaction
	today := time.Now().Format("2006-01-02")

	if err = db.WithContext(ctx).Where("queue_code = ?", queueCode).
		Where("DATE(created_at) = ?", today).
		Preload("Table").
		Preload("Orders").
		Preload("Orders.Menu").
		First(&transactionSchema).Error; err != nil {
		return transaction.Query{}, err
	}

	var transactionQuery transaction.Query
	transactionQuery.Transaction = schema.TransactionSchemaToEntity(transactionSchema)
	for i, orderSchema := range transactionSchema.Orders {
		transactionQuery.Orders = append(transactionQuery.Orders, transaction.OrderQuery{
			Order: schema.OrderSchemaToEntity(orderSchema),
		})
		transactionQuery.Orders[i].Menu = schema.MenuSchemaToEntity(*orderSchema.Menu)
	}
	transactionQuery.Table = schema.TableSchemaToEntity(*transactionSchema.Table)

	return transactionQuery, nil
}

func (r *transactionRepository) UpdateCookedAt(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return transaction.Transaction{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var transactionSchema schema.Transaction

	if err = db.WithContext(ctx).Model(&transactionSchema).Where("id = ?", transactionID).Update("cooked_at", time.Now()).Error; err != nil {
		return transaction.Transaction{}, err
	}

	transactionEntity := schema.TransactionSchemaToEntity(transactionSchema)
	return transactionEntity, nil
}

func (r *transactionRepository) UpdateTransactionCookingStatusStart(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return transaction.Transaction{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var transactionSchema schema.Transaction

	if err = db.WithContext(ctx).Model(&transactionSchema).Where("id = ?", transactionID).Update("order_status", transaction.OrderStatusPreparing).Error; err != nil {
		return transaction.Transaction{}, err
	}

	transactionEntity := schema.TransactionSchemaToEntity(transactionSchema)
	return transactionEntity, nil
}

func (r *transactionRepository) UpdateTransactionCookingStatusFinish(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return transaction.Transaction{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var transactionSchema schema.Transaction

	if err = db.WithContext(ctx).Model(&transactionSchema).Where("id = ?", transactionID).Update("order_status", transaction.OrderStatusReadyToServe).Error; err != nil {
		return transaction.Transaction{}, err
	}

	transactionEntity := schema.TransactionSchemaToEntity(transactionSchema)
	return transactionEntity, nil
}

func (r *transactionRepository) UpdateTransactionDeliveringStatusStart(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return transaction.Transaction{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var transactionSchema schema.Transaction

	if err = db.WithContext(ctx).Model(&transactionSchema).Where("id = ?", transactionID).Update("order_status", transaction.OrderStatusDelivering).Error; err != nil {
		return transaction.Transaction{}, err
	}

	transactionEntity := schema.TransactionSchemaToEntity(transactionSchema)
	return transactionEntity, nil
}

func (r *transactionRepository) UpdateTransactionDeliveringStatusFinish(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return transaction.Transaction{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var transactionSchema schema.Transaction

	if err = db.WithContext(ctx).Model(&transactionSchema).Where("id = ?", transactionID).Update("order_status", transaction.OrderStatusServed).Error; err != nil {
		return transaction.Transaction{}, err
	}

	transactionEntity := schema.TransactionSchemaToEntity(transactionSchema)
	return transactionEntity, nil
}

func (r *transactionRepository) UpdateServedAt(ctx context.Context, tx interface{}, transactionID string) (transaction.Transaction, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return transaction.Transaction{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var transactionSchema schema.Transaction

	if err = db.WithContext(ctx).Model(&transactionSchema).Where("id = ?", transactionID).Update("served_at", time.Now()).Error; err != nil {
		return transaction.Transaction{}, err
	}

	transactionEntity := schema.TransactionSchemaToEntity(transactionSchema)
	return transactionEntity, nil
}
