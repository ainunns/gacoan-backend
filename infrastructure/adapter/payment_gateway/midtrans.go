package payment_gateway

import (
	"context"
	"fmt"
	"fp-kpl/domain/identity"
	"fp-kpl/domain/port"
	"fp-kpl/domain/transaction"
	"fp-kpl/infrastructure/database/schema"
	"fp-kpl/infrastructure/database/validation"
	"os"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

type midtransAdapter struct {
	db                       *gorm.DB
	transactionDomainService transaction.Service
}

func NewMidtransAdapter(db *gorm.DB, transactionDomainService transaction.Service) port.PaymentGatewayPort {
	return &midtransAdapter{
		db:                       db,
		transactionDomainService: transactionDomainService,
	}
}

func (m midtransAdapter) ProcessPayment(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (port.ProcessPaymentResponse, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return port.ProcessPaymentResponse{}, err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = m.db
	}

	transactionSchema := schema.TransactionEntityToSchema(transactionEntity)

	err = db.WithContext(ctx).
		Preload("User").
		Preload("Orders").
		Preload("Orders.Menu").
		First(&transactionSchema, "id = ?", transactionSchema.ID.String()).Error
	if err != nil {
		return port.ProcessPaymentResponse{}, err
	}

	var s = snap.Client{}
	s.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	var itemDetails []midtrans.ItemDetails
	for _, orderSchema := range transactionSchema.Orders {
		menuSchema := orderSchema.Menu
		itemDetails = append(itemDetails, midtrans.ItemDetails{
			ID:    menuSchema.ID.String(),
			Name:  menuSchema.Name,
			Price: menuSchema.Price.IntPart(),
			Qty:   int32(orderSchema.Quantity),
		})
	}

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transactionSchema.ID.String(),
			GrossAmt: transactionSchema.TotalPrice.IntPart(),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: transactionSchema.User.Name,
			Email: transactionSchema.User.Email,
			Phone: transactionSchema.User.PhoneNumber,
		},
		Items: &itemDetails,
	}

	snapResp, snapErr := s.CreateTransaction(req)
	if snapErr != nil {
		return port.ProcessPaymentResponse{}, snapErr.RawError
	}
	return port.ProcessPaymentResponse{
		Token:       snapResp.Token,
		PaymentLink: snapResp.RedirectURL,
	}, nil
}

func (m midtransAdapter) HookPayment(ctx context.Context, tx interface{}, transactionId uuid.UUID, datas map[string]interface{}) error {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return err
	}

	db := validatedTransaction.DB()
	if db == nil {
		db = m.db
	}

	identityTransactionId := identity.NewIDFromSchema(transactionId)

	var transactionData schema.Transaction
	err = db.WithContext(ctx).
		Where("id = ?", identityTransactionId.String()).
		First(&transactionData).Error
	if err != nil {
		return err
	}

	status, ok := datas["transaction_status"].(string)
	if !ok {
		return fmt.Errorf("transaction_status is required in datas")
	}

	if !isValidPaymentStatus(status) {
		return fmt.Errorf("invalid payment status: %s", status)
	}

	transactionData.PaymentStatus = status
	transactionData.PaymentCode = datas["transaction_id"].(string)

	queueCode, err := m.transactionDomainService.GenerateQueueCode(ctx, transactionData.ID.String())
	if err != nil {
		return fmt.Errorf("failed to generate queue code: %w", err)
	}

	err = db.WithContext(ctx).
		Model(&transactionData).
		Updates(map[string]interface{}{
			"payment_status": transactionData.PaymentStatus,
			"payment_code":   transactionData.PaymentCode,
			"queue_code":     queueCode,
		}).Error
	if err != nil {
		return err
	}

	return nil
}

func isValidPaymentStatus(status string) bool {
	for _, paymentStatus := range transaction.PaymentStatuses {
		if paymentStatus == status {
			return true
		}
	}
	return false
}
