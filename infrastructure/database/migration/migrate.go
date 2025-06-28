package migration

import (
	"fp-kpl/infrastructure/database/schema"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&schema.User{},
		&schema.Table{},
		&schema.Category{},
		&schema.Menu{},
		&schema.Transaction{},
	); err != nil {
		return err
	}

	return nil
}
