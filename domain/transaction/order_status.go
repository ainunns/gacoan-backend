package transaction

import "fmt"

const (
	OrderStatusPending      = "pending"
	OrderStatusPreparing    = "preparing"
	OrderStatusReadyToServe = "ready_to_serve"
	OrderStatusDelivering   = "delivering"
	OrderStatusServed       = "served"
)

var (
	OrderStatuses = []string{
		OrderStatusPending,
		OrderStatusPreparing,
		OrderStatusReadyToServe,
		OrderStatusDelivering,
		OrderStatusServed,
	}
)

type OrderStatus struct {
	Status string
}

func NewOrderStatus(status string) (OrderStatus, error) {
	if !isValidOrderStatus(status) {
		return OrderStatus{}, fmt.Errorf("invalid order status: %s", status)
	}
	return OrderStatus{
		Status: status,
	}, nil
}

func isValidOrderStatus(status string) bool {
	for _, orderStatus := range OrderStatuses {
		if orderStatus == status {
			return true
		}
	}
	return false
}
