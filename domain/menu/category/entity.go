package category

import (
	"fp-kpl/domain/identity"
	"fp-kpl/domain/shared"
)

type Category struct {
	ID   identity.ID
	Name string
	shared.Timestamp
}
