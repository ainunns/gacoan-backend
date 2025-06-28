package repository

import (
	"context"
	"fp-kpl/domain/menu/category"
	"fp-kpl/infrastructure/database/db_transaction"
	"fp-kpl/infrastructure/database/schema"
	"fp-kpl/infrastructure/database/validation"
)

type categoryRepository struct {
	db *db_transaction.Repository
}

func NewCategoryRepository(db *db_transaction.Repository) category.Repository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetAllCategories(ctx context.Context, tx interface{}) ([]category.Category, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return nil, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var categorySchemas []schema.Category

	query := db.WithContext(ctx).Model(&schema.Category{})
	if err = query.Find(&categorySchemas).Error; err != nil {
		return nil, err
	}

	categoryEntities := make([]category.Category, len(categorySchemas))
	for i, categorySchema := range categorySchemas {
		categoryEntities[i] = schema.CategorySchemaToEntity(categorySchema)
	}

	return categoryEntities, nil
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, tx interface{}, id string) (category.Category, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return category.Category{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var categorySchema schema.Category

	if err = db.WithContext(ctx).Where("id = ?", id).Take(&categorySchema).Error; err != nil {
		return category.Category{}, err
	}

	categoryEntity := schema.CategorySchemaToEntity(categorySchema)
	return categoryEntity, nil
}
