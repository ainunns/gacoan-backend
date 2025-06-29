package response

import "github.com/shopspring/decimal"

type (
	TransactionCreate struct {
		TransactionID string                      `json:"transaction_id"`
		TotalPrice    string                      `json:"total_price"`
		Token         string                      `json:"token"`
		PaymentLink   string                      `json:"payment_link"`
		Orders        []OrderForTransactionCreate `json:"orders"`
	}

	OrderForTransactionCreate struct {
		Menu     MenuForTransaction `json:"menu"`
		Quantity int                `json:"quantity"`
	}

	MenuForTransaction struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Price string `json:"price"`
	}

	Transaction struct {
		ID           string                `json:"id"`
		QueueCode    string                `json:"queue_code"`
		EstimateTime string                `json:"estimate_time"`
		Orders       []OrderForTransaction `json:"orders"`
		TotalPrice   decimal.Decimal       `json:"total_price"`
		Table        Table                 `json:"table"`
		OrderStatus  string                `json:"order_status"`
		IsDelayed    bool                  `json:"is_delayed"`
	}

	OrderForTransaction struct {
		Menu     MenuForTransaction `json:"menu"`
		Quantity int                `json:"quantity"`
	}

	TransactionForWaiter struct {
		QueueCode string           `json:"queue_code"`
		Orders    []OrderForWaiter `json:"orders"`
		Table     TableForWaiter   `json:"table"`
	}
	OrderForWaiter struct {
		Menu     MenuForWaiter `json:"menu"`
		Quantity int           `json:"quantity"`
	}

	MenuForWaiter struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	TableForWaiter struct {
		ID          string `json:"id"`
		TableNumber string `json:"table_number"`
	}

	NextOrder struct {
		QueueCode string                `json:"queue_code"`
		Orders    []OrderForTransaction `json:"orders"`
	}

	StartCooking struct {
		QueueCode string                `json:"queue_code"`
		Orders    []OrderForTransaction `json:"orders"`
	}

	FinishCooking struct {
		QueueCode string                `json:"queue_code"`
		Orders    []OrderForTransaction `json:"orders"`
	}

	StartDelivering struct {
		QueueCode string                `json:"queue_code"`
		Orders    []OrderForTransaction `json:"orders"`
	}

	FinishDelivering struct{}
)
