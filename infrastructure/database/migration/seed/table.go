package seed

import (
	"fp-kpl/infrastructure/database/migration/data"
	"fp-kpl/infrastructure/database/table"

	"gorm.io/gorm"
)

func Table(db *gorm.DB) error {
	hasTable := db.Migrator().HasTable(&table.Table{})
	if !hasTable {
		return db.AutoMigrate(&table.Table{})
	}

	return db.CreateInBatches(data.Tables, 100).Error
}
