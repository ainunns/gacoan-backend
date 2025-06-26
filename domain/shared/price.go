package shared

import (
	"errors"

	"github.com/shopspring/decimal"
)

type Price struct {
	Price decimal.Decimal
}

func NewPrice(price decimal.Decimal) (Price, error) {
	if !isValidPrice(price) {
		return Price{}, errors.New("price must be greater than zero")
	}
	return Price{Price: price}, nil
}

func isValidPrice(price decimal.Decimal) bool {
	return price.GreaterThanOrEqual(decimal.Zero)
}
