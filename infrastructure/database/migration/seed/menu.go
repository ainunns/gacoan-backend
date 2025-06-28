package seed

import (
	"fp-kpl/infrastructure/database/migration/data"
	"fp-kpl/infrastructure/database/schema"

	"gorm.io/gorm"
)

func Menu(db *gorm.DB) error {
	hasTable := db.Migrator().HasTable(&schema.Menu{})
	if !hasTable {
		return db.AutoMigrate(&schema.Menu{})
	}

	menus := data.GetMenus(db)

	return db.CreateInBatches(menus, 100).Error
}
