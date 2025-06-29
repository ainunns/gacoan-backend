package response

import "github.com/shopspring/decimal"

type (
	TransactionCreate struct {
		TransactionID string                      `json:"transaction_id" form:"transaction_id" binding:"required"`
		TotalPrice    string                      `json:"total_price" form:"total_price" binding:"required"`
		PaymentLink   string                      `json:"payment_link" form:"payment_link" binding:"required"`
		Orders        []OrderForTransactionCreate `json:"orders" form:"orders" binding:"required"`
	}

	OrderForTransactionCreate struct {
		Menu     MenuForTransaction `json:"menu" form:"menu" binding:"required"`
		Quantity int                `json:"quantity" form:"quantity" binding:"required"`
	}

	MenuForTransaction struct {
		ID    string `json:"id" form:"id" binding:"required"`
		Name  string `json:"name" form:"name" binding:"required"`
		Price string `json:"price" form:"price" binding:"required"`
	}

	Transaction struct {
		ID           string                `json:"id" form:"id" binding:"required"`
		QueueCode    string                `json:"queue_code" form:"queue_code" binding:"required"`
		EstimateTime string                `json:"estimate_time" form:"estimate_time" binding:"required"`
		Orders       []OrderForTransaction `json:"orders" form:"orders" binding:"required"`
		TotalPrice   decimal.Decimal       `json:"total_price" form:"total_price" binding:"required"`
		Table        Table                 `json:"table" form:"table" binding:"required"`
		OrderStatus  string                `json:"order_status" form:"order_status" binding:"required"`
		IsDelayed    bool                  `json:"is_delayed" form:"is_delayed" binding:"required"`
	}

	OrderForTransaction struct {
		Menu     MenuForTransaction `json:"menu" form:"menu" binding:"required"`
		Quantity int                `json:"quantity" form:"quantity" binding:"required"`
	}
)
