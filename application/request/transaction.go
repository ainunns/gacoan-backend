package request

type (
	TransactionCreate struct {
		TableID string  `json:"table_id" form:"table_id" binding:"required"`
		Orders  []Order `json:"orders" form:"orders" binding:"required"`
	}

	Order struct {
		MenuID   string `json:"menu_id" form:"menu_id" binding:"required"`
		Quantity int    `json:"quantity" form:"quantity" binding:"required"`
	}
)
