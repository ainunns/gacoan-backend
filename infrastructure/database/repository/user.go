package repository

import (
	"context"
	"fp-kpl/domain/user"
	"fp-kpl/infrastructure/database/db_transaction"
	"fp-kpl/infrastructure/database/schema"
	"fp-kpl/infrastructure/database/validation"
)

type userRepository struct {
	db *db_transaction.Repository
}

func NewUserRepository(db *db_transaction.Repository) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) Register(ctx context.Context, tx interface{}, userEntity user.User) (user.User, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return user.User{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	userSchema := schema.UserEntityToSchema(userEntity)
	if err = db.WithContext(ctx).Create(&userSchema).Error; err != nil {
		return user.User{}, err
	}

	userEntity = schema.UserSchemaToEntity(userSchema)
	return userEntity, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, tx interface{}, id string) (user.User, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return user.User{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var userSchema schema.User
	if err = db.WithContext(ctx).Where("id = ?", id).Take(&userSchema).Error; err != nil {
		return user.User{}, err
	}

	userEntity := schema.UserSchemaToEntity(userSchema)
	return userEntity, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, tx interface{}, email string) (user.User, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return user.User{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var userSchema schema.User
	if err = db.WithContext(ctx).Where("email = ?", email).Take(&userSchema).Error; err != nil {
		return user.User{}, err
	}

	userEntity := schema.UserSchemaToEntity(userSchema)
	return userEntity, nil
}

func (r *userRepository) CheckEmail(ctx context.Context, tx interface{}, email string) (user.User, bool, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return user.User{}, false, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = r.db.DB()
	}

	var userSchema schema.User
	if err = db.WithContext(ctx).Where("email = ?", email).Take(&userSchema).Error; err != nil {
		return user.User{}, false, err
	}

	userEntity := schema.UserSchemaToEntity(userSchema)
	return userEntity, true, nil
}
