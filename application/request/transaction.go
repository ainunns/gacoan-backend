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

	StartCooking struct {
		QueueCode string `json:"queue_code" form:"queue_code" binding:"required"`
	}

	FinishCooking struct {
		QueueCode string `json:"queue_code" form:"queue_code" binding:"required"`
	}

	StartDelivering struct {
		QueueCode string `json:"queue_code" form:"queue_code" binding:"required"`
	}

	FinishDelivering struct {
		QueueCode string `json:"queue_code" form:"queue_code" binding:"required"`
	}
)
