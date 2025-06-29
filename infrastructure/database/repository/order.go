package repository

import (
	"context"
	"fp-kpl/domain/order"
	"fp-kpl/infrastructure/database/db_transaction"
	"fp-kpl/infrastructure/database/schema"
	"fp-kpl/infrastructure/database/validation"
)

type orderRepository struct {
	db *db_transaction.Repository
}

func NewOrderRepository(db *db_transaction.Repository) order.Repository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(ctx context.Context, tx interface{}, orderEntity order.Order) (order.Order, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return order.Order{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	orderSchema := schema.OrderEntityToSchema(orderEntity)
	if err = db.WithContext(ctx).Create(&orderSchema).Error; err != nil {
		return order.Order{}, err
	}

	orderEntity = schema.OrderSchemaToEntity(orderSchema)
	return orderEntity, nil
}

func (r *orderRepository) GetOrdersByTransactionID(ctx context.Context, tx interface{}, transactionID string) ([]order.Order, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return nil, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var orderSchemas []schema.Order
	if err = db.WithContext(ctx).Where("transaction_id = ?", transactionID).Find(&orderSchemas).Error; err != nil {
		return nil, err
	}

	orderEntities := make([]order.Order, len(orderSchemas))
	for i, orderSchema := range orderSchemas {
		orderEntities[i] = schema.OrderSchemaToEntity(orderSchema)
	}

	return orderEntities, nil
}
