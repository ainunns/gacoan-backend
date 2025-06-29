package seed

import (
	"fp-kpl/infrastructure/database/migration/data"
	"fp-kpl/infrastructure/database/schema"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Menu(db *gorm.DB) error {
	hasTable := db.Migrator().HasTable(&schema.Menu{})
	if !hasTable {
		return db.AutoMigrate(&schema.Menu{})
	}

	menus := data.GetMenus(db)

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoNothing: true,
	}).CreateInBatches(menus, 100).Error
}
