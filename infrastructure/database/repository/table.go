package repository

import (
	"context"
	"fp-kpl/domain/table"
	"fp-kpl/infrastructure/database/db_transaction"
	"fp-kpl/infrastructure/database/schema"
	"fp-kpl/infrastructure/database/validation"
)

type tableRepository struct {
	db *db_transaction.Repository
}

func NewTableRepository(db *db_transaction.Repository) table.Repository {
	return &tableRepository{db: db}
}

func (r *tableRepository) GetAllTables(ctx context.Context, tx interface{}) ([]table.Table, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return nil, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var tableSchemas []schema.Table

	query := db.WithContext(ctx).Model(&schema.Table{})
	if err = query.Find(&tableSchemas).Error; err != nil {
		return nil, err
	}

	tableEntities := make([]table.Table, len(tableSchemas))
	for i, tableSchema := range tableSchemas {
		tableEntities[i] = schema.TableSchemaToEntity(tableSchema)
	}

	return tableEntities, nil
}

func (r *tableRepository) GetTableByID(ctx context.Context, tx interface{}, id string) (table.Table, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return table.Table{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var tableSchema schema.Table

	if err = db.WithContext(ctx).Where("id = ?", id).Take(&tableSchema).Error; err != nil {
		return table.Table{}, err
	}

	tableEntity := schema.TableSchemaToEntity(tableSchema)
	return tableEntity, nil
}
