package migration

import (
	"fp-kpl/infrastructure/database/user"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&user.User{},
	); err != nil {
		return err
	}

	return nil
}
