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
	ImgUrl      string
	Price       shared.Price
	IsAvailable bool
	CookingTime time.Time
	Description string
	shared.Timestamp
}
