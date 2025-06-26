package shared

import "github.com/shopspring/decimal"

type Price struct {
	Price decimal.Decimal
}

func NewPrice(price float64) *Price {
	return &Price{Price: decimal.NewFromFloat(price)}
}

func (p *Price) GetPrice() decimal.Decimal {
	return p.Price
}

func (p *Price) GetPriceString() string {
	return p.Price.String()
}
