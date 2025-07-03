package order

import (
	"context"
	"fp-kpl/domain/shared"
	"github.com/shopspring/decimal"
)

type (
	Service interface {
		CalculatePrice(ctx context.Context, price shared.Price, quantity int64) (shared.Price, error)
	}

	service struct{}
)

func NewService() Service {
	return &service{}
}

func (s service) CalculatePrice(ctx context.Context, price shared.Price, quantity int64) (shared.Price, error) {
	if quantity <= 0 {
		return shared.Price{}, ErrorInvalidQuantity
	}

	orderPrice := price.Price.Mul(decimal.NewFromInt(quantity))
	return shared.NewPrice(orderPrice)
}
