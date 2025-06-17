package migration

import (
	"fp-kpl/infrastructure/database/refresh_token"
	"fp-kpl/infrastructure/database/user"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&user.User{},
		&refresh_token.RefreshToken{},
	); err != nil {
		return err
	}

	return nil
}
