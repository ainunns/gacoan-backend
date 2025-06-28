package response

import "github.com/shopspring/decimal"

type (
	Menu struct {
		ID          string          `json:"id"`
		Name        string          `json:"name"`
		Description string          `json:"description"`
		ImageUrl    string          `json:"image_url"`
		IsAvailable bool            `json:"is_available"`
		Price       decimal.Decimal `json:"price"`
		Category    Category        `json:"category"`
	}
)
