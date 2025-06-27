package table

import (
	"context"
	"fp-kpl/domain/table"
	"fp-kpl/infrastructure/database/db_transaction"
	"fp-kpl/infrastructure/database/validation"
)

type repository struct {
	db *db_transaction.Repository
}

func NewRepository(db *db_transaction.Repository) table.Repository {
	return &repository{db: db}
}

func (r *repository) GetAllTables(ctx context.Context, tx interface{}) ([]table.Table, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return nil, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var tableSchemas []Table

	query := db.WithContext(ctx).Model(&Table{})
	if err = query.Find(&tableSchemas).Error; err != nil {
		return nil, err
	}

	tableEntities := make([]table.Table, len(tableSchemas))
	for i, tableSchema := range tableSchemas {
		tableEntities[i] = SchemaToEntity(tableSchema)
	}

	return tableEntities, nil
}

func (r *repository) GetTableByID(ctx context.Context, tx interface{}, id string) (table.Table, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return table.Table{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var tableSchema Table

	if err = db.WithContext(ctx).Where("id = ?", id).Take(&tableSchema).Error; err != nil {
		return table.Table{}, err
	}

	tableEntity := SchemaToEntity(tableSchema)
	return tableEntity, nil
}
