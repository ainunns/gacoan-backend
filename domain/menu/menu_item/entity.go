package menu

import (
	"fp-kpl/domain/identity"
	"fp-kpl/domain/shared"
	"time"
)

type Menu struct {
	ID          identity.ID
	CategoryID  identity.ID
	Name        string
	ImageURL    shared.URL
	Price       shared.Price
	IsAvailable bool
	CookingTime time.Duration
	Description string
	shared.Timestamp
}
