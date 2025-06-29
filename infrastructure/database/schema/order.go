package schema

import (
	"fp-kpl/domain/identity"
	"fp-kpl/domain/order"
	"fp-kpl/domain/shared"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();column:id"`
	TransactionID uuid.UUID      `gorm:"type:uuid;not null;column:transaction_id"`
	MenuID        uuid.UUID      `gorm:"type:uuid;not null;column:menu_id"`
	Quantity      int            `gorm:"type:int;not null;column:quantity"`
	CreatedAt     time.Time      `gorm:"type:timestamp with time zone;column:created_at"`
	UpdatedAt     time.Time      `gorm:"type:timestamp with time zone;column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"type:timestamp with time zone;column:deleted_at"`

	Transaction *Transaction `gorm:""`
	Menu        *Menu        `gorm:""`
}

func OrderEntityToSchema(entity order.Order) Order {
	var deletedAtTime time.Time
	if entity.DeletedAt != nil {
		deletedAtTime = *entity.DeletedAt
	} else {
		deletedAtTime = time.Time{}
	}

	return Order{
		ID:            entity.ID.ID,
		TransactionID: entity.TransactionID.ID,
		MenuID:        entity.MenuID.ID,
		Quantity:      entity.Quantity,
		CreatedAt:     entity.Timestamp.CreatedAt,
		UpdatedAt:     entity.Timestamp.UpdatedAt,
		DeletedAt: gorm.DeletedAt{
			Time:  deletedAtTime,
			Valid: entity.DeletedAt != nil,
		},
	}
}

func OrderSchemaToEntity(schema Order) order.Order {
	return order.Order{
		ID:            identity.NewIDFromSchema(schema.ID),
		TransactionID: identity.NewIDFromSchema(schema.TransactionID),
		MenuID:        identity.NewIDFromSchema(schema.MenuID),
		Quantity:      schema.Quantity,
		Timestamp: shared.Timestamp{
			CreatedAt: schema.CreatedAt,
			UpdatedAt: schema.UpdatedAt,
			DeletedAt: &schema.DeletedAt.Time,
		},
	}
}
