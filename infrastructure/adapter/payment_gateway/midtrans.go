package payment_gateway

import (
	"context"
	"fp-kpl/domain/port"
	"fp-kpl/domain/transaction"
	"fp-kpl/infrastructure/database/schema"
	"fp-kpl/infrastructure/database/validation"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
	"os"
)

type midtransAdapter struct {
	db *gorm.DB
}

func NewMidtransAdapter(db *gorm.DB) port.PaymentGatewayPort {
	return &midtransAdapter{
		db: db,
	}
}

func (m midtransAdapter) ProcessPayment(ctx context.Context, tx interface{}, transactionEntity transaction.Transaction) (string, error) {
	validatedTransaction, err := validation.ValidateTransaction(tx)
	if err != nil {
		return "", err
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
		return "", err
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
		return "", snapErr.RawError
	}
	return snapResp.RedirectURL, nil
}
