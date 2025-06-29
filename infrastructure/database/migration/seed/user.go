package seed

import (
	"fp-kpl/infrastructure/database/migration/data"
	"fp-kpl/infrastructure/database/schema"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func User(db *gorm.DB) error {
	hasTable := db.Migrator().HasTable(&schema.User{})
	if !hasTable {
		return db.AutoMigrate(&schema.User{})
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoNothing: true,
	}).CreateInBatches(data.Users, 100).Error
}
