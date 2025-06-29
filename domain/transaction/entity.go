package transaction

import (
	"fp-kpl/domain/identity"
	"fp-kpl/domain/shared"
	"time"
)

type Transaction struct {
	ID          identity.ID
	UserID      identity.ID
	TableID     identity.ID
	Payment     Payment
	OrderStatus OrderStatus
	CookedAt    *time.Time
	ServedAt    *time.Time
	QueueCode   QueueCode
	TotalPrice  shared.Price
	shared.Timestamp
}
