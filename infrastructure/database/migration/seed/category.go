package seed

import (
	"fp-kpl/infrastructure/database/migration/data"
	"fp-kpl/infrastructure/database/schema"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Category(db *gorm.DB) error {
	hasTable := db.Migrator().HasTable(&schema.Category{})
	if !hasTable {
		return db.AutoMigrate(&schema.Category{})
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoNothing: true,
	}).CreateInBatches(data.Categories, 100).Error
}
