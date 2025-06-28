package schema

import (
	"fp-kpl/domain/identity"
	"fp-kpl/domain/shared"
	"fp-kpl/domain/transaction"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transaction struct {
	ID            uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();column:id"`
	UserID        uuid.UUID       `gorm:"type:uuid;not null;column:user_id"`
	TableID       uuid.UUID       `gorm:"type:uuid;not null;column:table_id"`
	PaymentCode   string          `gorm:"type:varchar(255);not null;column:payment_code"`
	PaymentStatus string          `gorm:"type:varchar(255);not null;column:payment_status"`
	OrderStatus   string          `gorm:"type:varchar(255);not null;column:order_status"`
	ServedAt      *time.Time      `gorm:"type:timestamp with time zone;column:served_at"`
	QueueCode     *string         `gorm:"type:varchar(255);column:queue_code"`
	TotalPrice    decimal.Decimal `gorm:"type:decimal(12,2);not null;default:0;column:total_price"`
	CreatedAt     time.Time       `gorm:"type:timestamp with time zone;column:created_at"`
	UpdatedAt     time.Time       `gorm:"type:timestamp with time zone;column:updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"type:timestamp with time zone;column:deleted_at"`

	User   *User   `gorm:"foreignKey:UserID"`
	Table  *Table  `gorm:"foreignKey:TableID"`
	Orders []Order `gorm:"foreignKey:TransactionID"`
}

func TransactionEntityToSchema(entity transaction.Transaction) Transaction {
	var deletedAtTime time.Time
	if entity.DeletedAt != nil {
		deletedAtTime = *entity.DeletedAt
	} else {
		deletedAtTime = time.Time{}
	}
	return Transaction{
		ID:            entity.ID.ID,
		UserID:        entity.UserID.ID,
		TableID:       entity.TableID.ID,
		PaymentCode:   entity.Payment.Code,
		PaymentStatus: entity.Payment.Status,
		OrderStatus:   entity.OrderStatus.Status,
		ServedAt:      entity.ServedAt,
		QueueCode:     &entity.QueueCode.Code,
		TotalPrice:    entity.TotalPrice.Price,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		DeletedAt: gorm.DeletedAt{
			Time:  deletedAtTime,
			Valid: entity.DeletedAt != nil,
		},
	}
}

func TransactionSchemaToEntity(schema Transaction) transaction.Transaction {
	var queueCode transaction.QueueCode
	if schema.QueueCode != nil {
		queueCode = transaction.NewQueueCodeFromSchema(*schema.QueueCode, true)
	} else {
		queueCode = transaction.QueueCode{
			Code:  "Q0000",
			Valid: false,
		}
	}
	return transaction.Transaction{
		ID:          identity.NewIDFromSchema(schema.ID),
		UserID:      identity.NewIDFromSchema(schema.UserID),
		TableID:     identity.NewIDFromSchema(schema.TableID),
		Payment:     transaction.NewPaymentFromSchema(schema.PaymentCode, schema.PaymentStatus),
		OrderStatus: transaction.NewOrderStatusFromSchema(schema.OrderStatus),
		ServedAt:    schema.ServedAt,
		QueueCode:   queueCode,
		TotalPrice:  shared.NewPriceFromSchema(schema.TotalPrice),
		Timestamp: shared.Timestamp{
			CreatedAt: schema.CreatedAt,
			UpdatedAt: schema.UpdatedAt,
			DeletedAt: &schema.DeletedAt.Time,
		},
	}
}
