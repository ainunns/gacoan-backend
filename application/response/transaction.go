package response

type (
	TransactionCreate struct {
		TransactionID string                      `json:"transaction_id" form:"transaction_id" binding:"required"`
		TotalPrice    string                      `json:"total_price" form:"total_price" binding:"required"`
		PaymentLink   string                      `json:"payment_link" form:"payment_link" binding:"required"`
		Orders        []OrderForTransactionCreate `json:"orders" form:"orders" binding:"required"`
	}

	OrderForTransactionCreate struct {
		Menu     MenuForTransactionCreate `json:"menu" form:"menu" binding:"required"`
		Quantity int                      `json:"quantity" form:"quantity" binding:"required"`
	}

	MenuForTransactionCreate struct {
		ID    string `json:"id" form:"id" binding:"required"`
		Name  string `json:"name" form:"name" binding:"required"`
		Price string `json:"price" form:"price" binding:"required"`
	}
)
