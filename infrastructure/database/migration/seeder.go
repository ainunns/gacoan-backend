package migration

import (
	"fp-kpl/infrastructure/database/migration/seed"

	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	if err := seed.Table(db); err != nil {
		return err
	}

	return nil
}
