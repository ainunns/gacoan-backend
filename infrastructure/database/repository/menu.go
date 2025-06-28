package repository

import (
	"context"
	menu "fp-kpl/domain/menu/menu_item"
	"fp-kpl/infrastructure/database/db_transaction"
	"fp-kpl/infrastructure/database/schema"
	"fp-kpl/infrastructure/database/validation"
)

type menuRepository struct {
	db *db_transaction.Repository
}

func NewMenuRepository(db *db_transaction.Repository) menu.Repository {
	return &menuRepository{db: db}
}

func (r *menuRepository) GetAllMenus(ctx context.Context, tx interface{}) ([]menu.Menu, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return nil, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var menuSchemas []schema.Menu

	query := db.WithContext(ctx).Model(&schema.Menu{})
	if err = query.Find(&menuSchemas).Error; err != nil {
		return nil, err
	}

	menuEntities := make([]menu.Menu, len(menuSchemas))
	for i, menuSchema := range menuSchemas {
		menuEntities[i] = schema.MenuSchemaToEntity(menuSchema)
	}

	return menuEntities, nil
}

func (r *menuRepository) GetMenuByID(ctx context.Context, tx interface{}, id string) (menu.Menu, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return menu.Menu{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var menuSchema schema.Menu

	if err = db.WithContext(ctx).Where("id = ?", id).Take(&menuSchema).Error; err != nil {
		return menu.Menu{}, err
	}

	menuEntity := schema.MenuSchemaToEntity(menuSchema)
	return menuEntity, nil
}

func (r *menuRepository) GetMenusByCategoryID(ctx context.Context, tx interface{}, categoryID string) ([]menu.Menu, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return nil, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var menuSchemas []schema.Menu

	if err = db.WithContext(ctx).Where("category_id = ?", categoryID).Find(&menuSchemas).Error; err != nil {
		return nil, err
	}

	menuEntities := make([]menu.Menu, len(menuSchemas))
	for i, menuSchema := range menuSchemas {
		menuEntities[i] = schema.MenuSchemaToEntity(menuSchema)
	}

	return menuEntities, nil
}

func (r *menuRepository) UpdateMenuAvailability(ctx context.Context, tx interface{}, id string, isAvailable bool) (menu.Menu, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return menu.Menu{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var menuSchema schema.Menu

	if err = db.WithContext(ctx).Where("id = ?", id).Take(&menuSchema).Error; err != nil {
		return menu.Menu{}, err
	}

	menuSchema.IsAvailable = isAvailable

	if err = db.WithContext(ctx).Save(&menuSchema).Error; err != nil {
		return menu.Menu{}, err
	}

	menuEntity := schema.MenuSchemaToEntity(menuSchema)
	return menuEntity, nil
}
