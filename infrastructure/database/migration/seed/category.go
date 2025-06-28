package seed

import (
	"fp-kpl/infrastructure/database/migration/data"
	"fp-kpl/infrastructure/database/schema"

	"gorm.io/gorm"
)

func Category(db *gorm.DB) error {
	hasTable := db.Migrator().HasTable(&schema.Category{})
	if !hasTable {
		return db.AutoMigrate(&schema.Category{})
	}

	return db.CreateInBatches(data.Categories, 100).Error
}
