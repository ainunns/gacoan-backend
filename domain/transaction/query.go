package transaction

import (
	menu "fp-kpl/domain/menu/menu_item"
	"fp-kpl/domain/order"
	"fp-kpl/domain/table"
)

type (
	Query struct {
		Transaction Transaction  `json:"transaction"`
		Orders      []OrderQuery `json:"orders"`
		Table       table.Table  `json:"table"`
	}

	OrderQuery struct {
		Order order.Order `json:"order"`
		Menu  menu.Menu   `json:"menu"`
	}
)
