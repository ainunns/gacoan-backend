package response

type TransactionDetail struct {
	ID              string `json:"id"`
	QueueCode       string `json:"queue_code"`
	PaymentCode     string `json:"payment_code"`
	PaymentStatus   string `json:"payment_status"`
	OrderStatus     string `json:"order_status"`
	ServedAt        string `json:"served_at"`
	TotalPrice      string `json:"total_price"`
	TableNumber     string `json:"table_number"`
	UserName        string `json:"user_name"`
	UserEmail       string `json:"user_email"`
	UserPhoneNumber string `json:"user_phone_number"`
}
