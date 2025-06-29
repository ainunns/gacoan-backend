package seed

import (
	"fp-kpl/infrastructure/database/migration/data"
	"fp-kpl/infrastructure/database/schema"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Table(db *gorm.DB) error {
	hasTable := db.Migrator().HasTable(&schema.Table{})
	if !hasTable {
		return db.AutoMigrate(&schema.Table{})
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "table_number"}},
		DoNothing: true,
	}).CreateInBatches(data.Tables, 100).Error
}
