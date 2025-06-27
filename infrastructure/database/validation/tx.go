package validation

import (
	"fp-kpl/infrastructure/database/db_transaction"
	"gorm.io/gorm"
)

func ValidateTransaction(tx interface{}) (*db_transaction.Repository, error) {
	db := &db_transaction.Repository{}
	if tx == nil {
		return db, nil
	}

	db, ok := tx.(*db_transaction.Repository)
	if !ok {
		return nil, gorm.ErrInvalidTransaction
	}

	return db, nil
}
