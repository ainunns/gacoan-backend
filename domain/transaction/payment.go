package transaction

import "fmt"

const (
	PaymentStatusCapture    = "capture"
	PaymentStatusSettlement = "settlement"
	PaymentStatusCancel     = "cancel"
	PaymentStatusDeny       = "deny"
	PaymentStatusExpire     = "expire"
	PaymentStatusPending    = "pending"
)

var (
	PaymentStatuses = []string{
		PaymentStatusCapture,
		PaymentStatusSettlement,
		PaymentStatusCancel,
		PaymentStatusDeny,
		PaymentStatusExpire,
		PaymentStatusPending,
	}
)

type Payment struct {
	Code   string
	Status string
}

func NewPayment(code, status string) (Payment, error) {
	if !isValidPaymentStatus(status) {
		return Payment{}, fmt.Errorf("invalid payment status: %s", status)
	}
	return Payment{
		Code:   code,
		Status: status,
	}, nil
}

func NewPaymentFromSchema(code, status string) Payment {
	return Payment{
		Code:   code,
		Status: status,
	}
}

func isValidPaymentStatus(status string) bool {
	for _, paymentStatus := range PaymentStatuses {
		if paymentStatus == status {
			return true
		}
	}
	return false
}
