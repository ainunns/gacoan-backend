package table

import (
	"fp-kpl/domain/identity"
	"fp-kpl/domain/shared"
)

type Table struct {
	ID          identity.ID
	TableNumber string
	shared.Timestamp
}
