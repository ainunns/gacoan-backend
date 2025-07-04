package order

import "errors"

var (
	ErrorInvalidQuantity          = errors.New("invalid quantity, must be greater than zero")
	ErrorGetOrdersByTransactionID = errors.New("failed to get orders by transaction id")
)
