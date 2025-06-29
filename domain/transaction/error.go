package transaction

import (
	"errors"
)

var (
	ErrorInvalidTransaction = errors.New("invalid transaction")
	ErrorGetAllTransactions = errors.New("failed to get all transactions")
)
