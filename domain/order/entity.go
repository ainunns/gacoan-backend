package order

import (
	"fp-kpl/domain/identity"
	"fp-kpl/domain/shared"
)

type Order struct {
	ID            identity.ID
	TransactionID identity.ID
	MenuID        identity.ID
	Quantity      int
	shared.Timestamp
}
